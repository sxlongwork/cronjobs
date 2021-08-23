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
	fmt.Println("connect etcd success")
}

func main() {
	var (
		lease    clientv3.Lease
		kv       clientv3.KV
		leaseRes *clientv3.LeaseGrantResponse
		leaseid  clientv3.LeaseID
		putRes   *clientv3.PutResponse
		getRes   *clientv3.GetResponse
	)
	// 创建Lease
	lease = clientv3.NewLease(cli)

	// 申请一个10s的租约
	if leaseRes, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println("grant new lease failed.")
	}
	// 获取租约id
	leaseid = leaseRes.ID

	kv = clientv3.NewKV(cli)
	// put数据，将键值对与租约绑定，到期后，自动删除该键值对
	if putRes, err = kv.Put(context.TODO(), "job3", "doin", clientv3.WithLease(leaseid)); err != nil {
		fmt.Println("put data fail.")
	} else {
		fmt.Println("put data success, Revision:", putRes.Header.Revision)
	}
	// 定义1个for循环，检查10后key是否过期了
	for {
		time.Sleep(2 * time.Second)
		if getRes, err = kv.Get(context.TODO(), "job3"); err != nil {
			fmt.Println("get data fail.")
		} else {
			if getRes.Kvs != nil {
				for _, v := range getRes.Kvs {
					fmt.Println("还未获取，值为", string(v.Value))
				}
			} else {
				fmt.Println("key过期删除了")
				break
			}

		}

	}

}
