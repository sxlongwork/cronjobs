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
		fmt.Println("connect to etcd fail")
		return
	}

}

func main() {

	var (
		lease         clientv3.Lease
		leasegrantRes *clientv3.LeaseGrantResponse
		leaseID       clientv3.LeaseID
		leaseRes      *clientv3.LeaseKeepAliveResponse
		leaseChan     <-chan *clientv3.LeaseKeepAliveResponse
		kv            clientv3.KV
		putRes        *clientv3.PutResponse
		getRes        *clientv3.GetResponse
		ctx           context.Context
		cancelFunc    context.CancelFunc
	)
	// 创建1个租约
	lease = clientv3.NewLease(cli)

	// 申请租约
	if leasegrantRes, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println("grant a new lease failed")
	} else {
		leaseID = leasegrantRes.ID
	}

	// 自动续租
	// 1、5s后自动超时，取消续租，则key存在时长= 5 + 10(申请租约时的时长) = 15s；即5s后取消续租，key不会立即被删，还有10s存活时间
	// 2、一直续租，key一直存在
	ctx, cancelFunc = context.WithTimeout(context.TODO(), 5*time.Second)
	if leaseChan, err = lease.KeepAlive(ctx, leaseID); err != nil {
		// if leaseChan, err = lease.KeepAlive(context.TODO(), leaseID); err != nil {
		fmt.Println("续租失败")
	}
	go func() {
		for {
			select {
			case leaseRes = <-leaseChan:
				if leaseRes == nil {
					fmt.Println("租约终止")
					goto end
				} else {
					fmt.Println("续租成功")
					fmt.Println("TTL:", leaseRes.TTL)
				}
			}
		}
	end:
	}()

	// KV put data
	kv = clientv3.NewKV(cli)
	// 绑定租约，如果没有续租10s会过期被删除
	if putRes, err = kv.Put(context.TODO(), "job4", "omg", clientv3.WithLease(leaseID)); err != nil {
		fmt.Println("put data fail")
	} else {
		fmt.Println("put data success, revision=", putRes.Header.Revision)
	}
	cancelFunc()
	start := time.Now().Unix()

	for {
		time.Sleep(5 * time.Second)
		if getRes, err = kv.Get(context.TODO(), "job4"); err != nil {
			fmt.Println("get data fail")
		} else {
			if getRes.Kvs == nil {
				fmt.Println("过期")
				end := time.Now().Unix()
				fmt.Println("key存活时长 =", end-start)
				break
			} else {
				for _, v := range getRes.Kvs {
					fmt.Println("value =", string(v.Value))
				}
			}
		}
	}
}
