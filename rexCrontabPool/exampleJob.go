package rexCrontabPool

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/rootexit/rexLib/rexDao"
	"github.com/rootexit/rexLib/rexRequest"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type PeriodicJob struct {
	dao           rexDao.Dao
	store         rexDao.RedisDao
	rs            *redsync.Redsync
	RequestClient rexRequest.RequestClient
	Data          PeriodicJobData
	fn            func(taskUuid, taskName string)
}

type PeriodicJobData struct {
	taskUuid string
	taskName string
}

func (j *PeriodicJob) Run() {
	// 分布式锁 key
	mutex := j.rs.NewMutex(fmt.Sprintf("lock-%s", j.Data.taskUuid), redsync.WithExpiry(5*time.Minute))
	if err := mutex.Lock(); err == nil {
		// 获得锁，执行任务
		defer mutex.Unlock()
		logx.Infof("任务uuid: %s, 任务名称: %s, 我来执行任务", j.Data.taskUuid, j.Data.taskName)
		j.fn(j.Data.taskUuid, j.Data.taskName)
	} else {
		logx.Infof("任务uuid: %s, 任务名称: %s, 其他节点已在执行任务，跳过", j.Data.taskUuid, j.Data.taskName)
	}
}

func AddTask(uuidStr string, taskName, spec string, fn func(taskUuid, taskName string)) (taskUuid string, err error) {
	// note: 添加一个周期任务
	//job := &Task{
	//	TaskUuid: uuidStr,
	//	Name:     taskName,
	//	Spec:     spec,
	//	Job: &PeriodicJob{
	//		//dao:           s.Dao,
	//		//store:         s.StoreDao,
	//		//rs:            s.RdSync,
	//		//RequestClient: s.RequestClient,
	//		fn: fn,
	//		Data: PeriodicJobData{
	//			taskUuid: uuidStr,
	//			taskName: taskName,
	//		},
	//	},
	//}
	//CrontabPool.Register <- job
	return uuidStr, nil
}
