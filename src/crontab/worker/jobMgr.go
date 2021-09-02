package worker

import (
	"context"
	"cronjobs/src/crontab/common"
	"cronjobs/src/crontab/worker/config"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var GOL_JOBMGR *JobMgr

/*
任务管理器初始化
*/
func InitJobMgr() (err error) {
	var (
		cfg     clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)
	cfg = clientv3.Config{
		Endpoints:   config.GOL_CONFIG.Endpoints,
		DialTimeout: time.Duration(config.GOL_CONFIG.DialTimeout) * time.Millisecond,
	}
	if client, err = clientv3.New(cfg); err != nil {
		return
	}
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)
	GOL_JOBMGR = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}
	// 启动监听新增/修改/删除任务协程
	GOL_JOBMGR.WatchJobs()
	// 启动监听杀死任务协程
	GOL_JOBMGR.watchKillJob()

	return
}

func (jobMgr *JobMgr) CreateLock(jobname string) (lock *JobLock) {
	lock = InitLock(jobname, jobMgr.kv, jobMgr.lease)
	return
}

/*
监听/cron/job/下的任务变化，新增/修改/删除
*/
func (jobMgr *JobMgr) WatchJobs() (err error) {
	var (
		getRes    *clientv3.GetResponse
		job       *common.Job
		creation  int64
		watchChan clientv3.WatchChan
		watchRes  clientv3.WatchResponse
		event     *clientv3.Event
		jobEvent  *common.JobEvent
		jobName   string
	)

	if getRes, err = GOL_JOBMGR.client.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	for _, value := range getRes.Kvs {
		// 反序列化job
		if job, err = common.Unmarshal(value.Value); err != nil {
			return
		}
		// 将job推送给调度协程
		jobEvent = common.BuildJobEvent(common.JOB_PUT_EVENT, job)
		GOL_SCHDULE.PutJobEvent(jobEvent)
		// fmt.Println("jobevent:", jobEvent.Job.JobName)

	}

	// 获取当前的creation，并从当前版本开始一直监听该目录
	go func() {
		creation = getRes.Header.Revision + 1
		watchChan = GOL_JOBMGR.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix(), clientv3.WithRev(creation))
		for watchRes = range watchChan {
			for _, event = range watchRes.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					// put操作
					if job, err = common.Unmarshal(event.Kv.Value); err != nil {
						continue
					}
					// 构建event，将job推送给调度协程
					jobEvent = common.BuildJobEvent(common.JOB_PUT_EVENT, job)
				case clientv3.EventTypeDelete:
					// DELETE操作
					jobName = common.GetJobSaveName(string(event.Kv.Key))
					job = &common.Job{JobName: jobName}
					// 构建event，将job推送给调度协程
					jobEvent = common.BuildJobEvent(common.JOB_DEL_EVENT, job)
				}
				// 推送jobEvent到调度协程
				GOL_SCHDULE.PutJobEvent(jobEvent)
				// fmt.Println("jobevent:", jobEvent.EventType, jobEvent.Job.JobName)
			}
		}
	}()

	return

}

/*
监听/cron/kill/目录，put操作杀死任务，delete操作忽略
*/
func (jobMgr *JobMgr) watchKillJob() (err error) {
	var (
		watchChan clientv3.WatchChan
		watchRes  clientv3.WatchResponse
		event     *clientv3.Event
		job       *common.Job
		jobEvent  *common.JobEvent
	)
	go func() {
		// 监听强杀任务事件 /cron/kill/
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_KILL_DIR, clientv3.WithPrefix())

		for watchRes = range watchChan {
			for _, event = range watchRes.Events {
				switch event.Type {
				case clientv3.EventTypePut:

					job = &common.Job{JobName: common.GetJobKillName(string(event.Kv.Key))}
					// 构建event，将job推送给调度协程
					jobEvent = common.BuildJobEvent(common.JOB_KILL_EVENT, job)
					// 推送jobEvent到调度协程
					GOL_SCHDULE.PutJobEvent(jobEvent)
				case clientv3.EventTypeDelete:
				}

			}
		}
	}()
	return
}
