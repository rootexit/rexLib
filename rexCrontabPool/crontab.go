package rexCrontabPool

import (
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type CrontabPool struct {
	Cron       *cron.Cron
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
	c.Cron.Start()
	defer c.Cron.Stop()

	for {
		select {
		case num := <-c.Close:
			logx.Infof("pool close signal = %d", num)
			return
		case task := <-c.Register:
			if task == nil {
				logx.Errorf("register task is nil")
				continue
			}
			c.Lock.Lock()
			// note: 如果任务已经存在，则不再注册，但是万一需要更新，这个需要改进
			if _, ok := c.TaskPool[task.TaskUuid]; ok {
				logx.Infof("task already exists, and task uuid = %s, task name = %s, and task ID = %d", task.TaskUuid, task.Name, task.JobId)
				c.Lock.Unlock()
				continue
			}
			c.Lock.Unlock()

			jobId, err := c.Cron.AddJob(task.Spec, task.Job)
			if err != nil {
				logx.Errorf("task register failed, and task uuid = %s, task name = %s, and task ID = %d", task.TaskUuid, task.Name, task.JobId)
				c.Lock.Unlock()
				continue
			}

			task.JobId = jobId
			c.TaskPool[task.TaskUuid] = task
			c.TaskCount = len(c.TaskPool)
			logx.Infof("task register success, and task uuid = %s, task name = %s, and task ID = %d", task.TaskUuid, task.Name, task.JobId)
		case task := <-c.Update:
			if task == nil {
				logx.Errorf("update task is nil")
				continue
			}

			jobId, err := c.Cron.AddJob(task.Spec, task.Job)
			if err != nil {
				logx.Errorf("task register failed, and task uuid = %s, task name = %s, and task ID = %d", task.TaskUuid, task.Name, task.JobId)
				c.Lock.Unlock()
				continue
			}

			c.Lock.Lock()
			// note: 如果是更新通道，说明有可能任务会存在
			if oldTask, ok := c.TaskPool[task.TaskUuid]; ok {
				c.Cron.Remove(oldTask.JobId)
			}
			c.Lock.Unlock()

			task.JobId = jobId
			c.TaskPool[task.TaskUuid] = task
			logx.Infof("task update success, and task uuid = %s, task name = %s, and task ID = %d", task.TaskUuid, task.Name, task.JobId)
		case taskUuid := <-c.UnRegister:
			//注销客户端
			c.Lock.Lock()
			if _, ok := c.TaskPool[taskUuid]; ok {
				c.Cron.Remove(c.TaskPool[taskUuid].JobId)
				//删除分组中的任务
				delete(c.TaskPool, taskUuid)
				//任务数量减1
				c.TaskCount = len(c.TaskPool)
				logx.Infof("task unregister success, and task uuid = %s", taskUuid)
			}
			c.Lock.Unlock()
		}
	}
}
