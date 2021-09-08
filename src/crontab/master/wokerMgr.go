package master

import (
	"context"
	"cronjobs/src/crontab/common"
	"cronjobs/src/crontab/master/config"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

var GOL_WORKERMGR *WorkerMgr

type WorkerMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

func InitWorkerMgr() (err error) {
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
	GOL_WORKERMGR = &WorkerMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

/*
获取worker节点列表
*/
func (worker *WorkerMgr) GetWorkerList() (workerList []string, err error) {
	var (
		getRes *clientv3.GetResponse
		kv     *mvccpb.KeyValue
	)
	workerList = make([]string, 0)
	if getRes, err = worker.client.Get(context.TODO(), common.WORKER_REGISTER_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	for _, kv = range getRes.Kvs {
		// string(kv.Key)	/cron/worker/169.254.122.64/16
		workerList = append(workerList, common.GetSuffixName(string(kv.Key), common.WORKER_REGISTER_DIR))
	}
	return
}
