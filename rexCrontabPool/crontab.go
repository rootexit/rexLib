package rexCrontabPool

import (
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type CrontabPool struct {
	*cron.Cron
	TaskPool   map[string]*Task
	Lock       sync.Mutex
	Register   chan *Task
	Update     chan *Task
	UnRegister chan string
	Close      chan int
	TaskCount  int
}

type Task struct {
	TaskUuid string // note: 最好用uuid
	Name     string
	Spec     string
	JobId    cron.EntryID
	Job      cron.Job
}

func NewCrontabPool() *CrontabPool {
	c := cron.New(cron.WithSeconds())
	return &CrontabPool{
		Cron:       c,
		TaskPool:   make(map[string]*Task),
		Lock:       sync.Mutex{},
		Register:   make(chan *Task),
		Update:     make(chan *Task),
		UnRegister: make(chan string),
		Close:      make(chan int),
		TaskCount:  0,
	}
}

func (c *CrontabPool) Run() {
	c.Start()
	defer c.Stop()

	for {
		select {
		case num := <-c.Close:
			logx.Infof("pool close signal = %d", num)
			break
		case task := <-c.Register:
			c.Lock.Lock()
			// note: 如果任务已经存在，则不再注册，但是万一需要更新，这个需要改进
			if _, ok := c.TaskPool[task.TaskUuid]; ok {
				logx.Infof("task already exists, and task uuid = %s, task name = %s, and task ID = %d", task.TaskUuid, task.Name, task.JobId)
				return
			}
			jobId, err := c.AddJob(task.Spec, task.Job)
			if err != nil {
				logx.Errorf("task register failed, and task uuid = %s, task name = %s, and task ID = %d", task.TaskUuid, task.Name, task.JobId)
				return
			}
			task.JobId = jobId
			c.TaskPool[task.TaskUuid] = task
			c.TaskCount += 1
			c.Lock.Unlock()
			logx.Infof("task register success, and task uuid = %s, task name = %s, and task ID = %d", task.TaskUuid, task.Name, task.JobId)
		case task := <-c.Update:
			c.Lock.Lock()
			// note: 如果是更新通道，说明有可能任务会存在
			if _, ok := c.TaskPool[task.TaskUuid]; ok {
				c.Remove(c.TaskPool[task.TaskUuid].JobId)
				//删除分组中的任务
				delete(c.TaskPool, task.TaskUuid)
				//任务数量减1
				c.TaskCount -= 1
				logx.Infof("task delete old success, and task uuid = %s", task.TaskUuid)
			}
			jobId, err := c.AddJob(task.Spec, task.Job)
			if err != nil {
				logx.Errorf("task register failed, and task uuid = %s, task name = %s, and task ID = %d", task.TaskUuid, task.Name, task.JobId)
				return
			}
			task.JobId = jobId
			c.TaskPool[task.TaskUuid] = task
			c.TaskCount += 1
			c.Lock.Unlock()
			logx.Infof("task update success, and task uuid = %s, task name = %s, and task ID = %d", task.TaskUuid, task.Name, task.JobId)
		case taskUuid := <-c.UnRegister:
			//注销客户端
			c.Lock.Lock()
			if _, ok := c.TaskPool[taskUuid]; ok {
				c.Remove(c.TaskPool[taskUuid].JobId)
				//删除分组中的任务
				delete(c.TaskPool, taskUuid)
				//任务数量减1
				c.TaskCount -= 1
				logx.Infof("task unregister success, and task uuid = %s", taskUuid)
			}
			c.Lock.Unlock()
		}
	}
}
