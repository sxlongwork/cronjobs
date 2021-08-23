package master

import (
	"context"
	"cronjobs/src/crontab/common"
	"cronjobs/src/crontab/master/config"
	"encoding/json"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var GOL_JOBMGR *JobMgr

func InitJobMgr() (err error) {
	var (
		cfg    clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
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
	GOL_JOBMGR = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var (
		key    string
		putRes *clientv3.PutResponse
		data   []byte
	)
	key = common.JOB_SAVE_DIR + job.JobName
	if data, err = json.Marshal(job); err != nil {
		return
	}

	if putRes, err = jobMgr.kv.Put(context.TODO(), key, string(data), clientv3.WithPrevKV()); err != nil {
		return
	}

	if putRes.PrevKv != nil {
		oldJob = &common.Job{}
		if err = json.Unmarshal(putRes.PrevKv.Value, oldJob); err != nil {
			return
		}
	}
	return
}

func (jobMgr *JobMgr) DelJob(name string) (oldJob *common.Job, err error) {

	var (
		key    string
		delRes *clientv3.DeleteResponse
	)
	key = common.JOB_SAVE_DIR + name

	if delRes, err = GOL_JOBMGR.kv.Delete(context.TODO(), key, clientv3.WithPrevKV()); err != nil {
		return
	}

	if len(delRes.PrevKvs) != 0 {
		if err = json.Unmarshal(delRes.PrevKvs[0].Value, oldJob); err != nil {
			err = nil
			return
		}
	}
	return
}

func (JobMgr *JobMgr) ListJob() (jobs []*common.Job, err error) {
	//
	var (
		getRes *clientv3.GetResponse

		job *common.Job
	)
	jobs = make([]*common.Job, 0)
	if getRes, err = GOL_JOBMGR.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	if len(getRes.Kvs) != 0 {
		for _, v := range getRes.Kvs {
			// fmt.Println(string(v.Value))
			job = &common.Job{}
			if err = json.Unmarshal(v.Value, job); err != nil {
				return
			}

			jobs = append(jobs, job)
		}
		// for _, v := range jobs {
		// 	fmt.Println(*v)
		// }
	}
	return
}

func (JobMgr *JobMgr) KillJob(name string) (err error) {
	// 将需要杀掉的任务名称存储到/cron/kill/目录下即可，worker中有协程会监听这个目录
	var (
		key           string
		leaseGraneRes *clientv3.LeaseGrantResponse
		leaseID       clientv3.LeaseID
		// putRes        *clientv3.PutResponse
	)
	key = common.JOB_KILLER_DIR + name

	// 申请租约，设置1s后自动删除该key
	if leaseGraneRes, err = GOL_JOBMGR.lease.Grant(context.TODO(), 1); err != nil {
		return
	}
	leaseID = leaseGraneRes.ID

	// put key
	if _, err = GOL_JOBMGR.kv.Put(context.TODO(), key, "", clientv3.WithLease(leaseID)); err != nil {
		return
	}
	return
}
