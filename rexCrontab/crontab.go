package rexCrontab

import (
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type Crontab struct {
	*cron.Cron
	TaskPool             map[string]*Task
	Lock                 sync.Mutex
	Register, UnRegister chan *Task
	Close                chan int
	TaskCount            int
}

type Task struct {
	TaskId string // note: 最好用uuid
	Name   string
	Spec   string
	JobId  cron.EntryID
	Job    cron.Job
}

func New() *Crontab {
	c := cron.New(cron.WithSeconds())
	return &Crontab{
		Cron:       c,
		TaskPool:   make(map[string]*Task),
		Lock:       sync.Mutex{},
		Register:   make(chan *Task),
		UnRegister: make(chan *Task),
		Close:      make(chan int),
		TaskCount:  0,
	}
}

func (c *Crontab) Run() {
	c.Start()
	defer c.Stop()

	for {
		select {
		case num := <-c.Close:
			logx.Infof("pool close signal = %d", num)
			break
		case task := <-c.Register:
			//注册定时任务
			c.Lock.Lock()
			jobId, err := c.AddJob(task.Spec, task.Job)
			if err != nil {
				logx.Errorf("task register failed, and task id = %s, task name = %s, and task ID = %d", task.TaskId, task.Name, task.JobId)
				return
			}
			task.JobId = jobId
			c.TaskPool[task.TaskId] = task
			c.TaskCount += 1
			c.Lock.Unlock()
			logx.Infof("task register success, and task id = %s, task name = %s, and task ID = %d", task.TaskId, task.Name, task.JobId)
		case task := <-c.UnRegister:
			//注销客户端
			c.Lock.Lock()
			if _, ok := c.TaskPool[task.TaskId]; ok {
				c.Remove(task.JobId)
				//删除分组中的任务
				delete(c.TaskPool, task.TaskId)
				//任务数量减1
				c.TaskCount -= 1
				logx.Infof("task unregister success, and task id = %s, task name = %s, and task ID = %d", task.TaskId, task.Name, task.JobId)
			}
			c.Lock.Unlock()
		}
	}
}
