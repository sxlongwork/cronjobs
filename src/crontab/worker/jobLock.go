package worker

import (
	"context"
	"cronjobs/src/crontab/common"

	"github.com/coreos/etcd/clientv3"
)

type JobLock struct {
	jobName    string
	kv         clientv3.KV
	lease      clientv3.Lease
	leaseID    clientv3.LeaseID
	cancelFunc context.CancelFunc
	isLocked   bool
}

/*
初始化锁
*/
func InitLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (lock *JobLock) {
	lock = &JobLock{
		jobName: jobName,
		kv:      kv,
		lease:   lease,
	}
	return
}

/*
抢锁
*/
func (lock *JobLock) TryLock() (err error) {
	var (
		leaseGrantRes          *clientv3.LeaseGrantResponse
		leasID                 clientv3.LeaseID
		leaseChan              <-chan *clientv3.LeaseKeepAliveResponse
		leaseKeepAliveResponse *clientv3.LeaseKeepAliveResponse
		txn                    clientv3.Txn
		lockKey                string
		txnRes                 *clientv3.TxnResponse
		cxt                    context.Context
		cancelFunc             context.CancelFunc
	)
	// 申请租约
	if leaseGrantRes, err = lock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}
	leasID = leaseGrantRes.ID
	// 创建上下文和取消函数
	cxt, cancelFunc = context.WithCancel(context.TODO())

	//自动续租
	if leaseChan, err = lock.lease.KeepAlive(cxt, leasID); err != nil {
		goto ERR
	}
	// defer cancelFunc()
	// defer lock.lease.Revoke(context.TODO(), lock.leaseID)
	go func() {
		for {
			select {
			case leaseKeepAliveResponse = <-leaseChan:
				if leaseKeepAliveResponse == nil {
					// err = common.TRY_LOCK_ERROR
					goto END
				}
			}
		}
	END:
	}()
	// 创建事务
	lockKey = common.JOB_LOCK_PREFIX + lock.jobName
	txn = lock.kv.Txn(context.TODO())
	txn = txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leasID))).
		Else(clientv3.OpGet(lockKey))
	// 提交事务
	if txnRes, err = txn.Commit(); err != nil {
		goto ERR
	}
	if !txnRes.Succeeded {
		err = common.TRY_LOCK_ERROR
		goto ERR
	}
	// 抢锁成功
	lock.leaseID = leasID
	lock.cancelFunc = cancelFunc
	lock.isLocked = true
	return
ERR:
	cancelFunc()
	lock.lease.Revoke(context.TODO(), lock.leaseID)
	return
}

/*
释放锁
*/
func (lock *JobLock) Unlock() {
	if lock.isLocked {
		lock.cancelFunc()
		lock.lease.Revoke(context.TODO(), lock.leaseID)
	}
}
