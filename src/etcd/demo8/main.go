package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var (
	config clientv3.Config
	cli    *clientv3.Client
	err    error
)

func init() {
	config = clientv3.Config{
		Endpoints:   []string{"47.93.208.134:2379"},
		DialTimeout: 3 * time.Second,
	}
	if cli, err = clientv3.New(config); err != nil {
		fmt.Println("connect to etcd fail.")
		return
	}
}

func main() {
	var (
		lease         clientv3.Lease
		leaseGrantRes *clientv3.LeaseGrantResponse
		leaseID       clientv3.LeaseID
		leaseChan     <-chan *clientv3.LeaseKeepAliveResponse
		leaseRes      *clientv3.LeaseKeepAliveResponse
		ctx           context.Context
		cancelFunc    context.CancelFunc
		kv            clientv3.KV
		txn           clientv3.Txn
		txnRes        *clientv3.TxnResponse
	)

	// 1、创建租约,自动续租，拿着租约去抢key
	lease = clientv3.NewLease(cli)
	if leaseGrantRes, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println("grant new lease error")
		return
	}
	leaseID = leaseGrantRes.ID

	ctx, cancelFunc = context.WithCancel(context.TODO())
	if leaseChan, err = lease.KeepAlive(ctx, leaseID); err != nil {
		fmt.Println("自动续租失败")
		return
	}
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseID)

	go func() {
		for {
			select {
			case leaseRes = <-leaseChan:
				if leaseChan == nil {
					fmt.Println("租约终止")
					goto end
				} else {
					fmt.Println("续租成功", leaseRes.ID)
				}
			}
		}
	end:
	}()
	// 2、创建事务，抢到key，执行put操作；没抢到
	kv = clientv3.NewKV(cli)
	txn = kv.Txn(context.TODO())
	// 如果key的创建revision为0，说明key还没有被创建，抢到锁。可以执行put操作
	if txnRes, err = txn.If(clientv3.Compare(clientv3.CreateRevision("job3"), "=", 0)).
		Then(clientv3.OpPut("job3", "xxx", clientv3.WithLease(leaseID))).
		Else(clientv3.OpGet("job3")).
		Commit(); err != nil {
		fmt.Println("txn error")
		return
	}
	// 判断是否抢锁成功
	if !txnRes.Succeeded {
		fmt.Println("抢锁失败", string(txnRes.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	// 抢锁成功, 处理任务
	fmt.Println("处理任务")

	time.Sleep(5 * time.Second)
	// 3、释放资源
	// defer
}
