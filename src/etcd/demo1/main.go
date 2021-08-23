package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	var (
		cfg    clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		putres *clientv3.PutResponse
		err    error
	)
	cfg = clientv3.Config{
		Endpoints:   []string{"47.93.208.134:2379"},
		DialTimeout: 3 * time.Second,
	}
	// 创建一个客户端
	if client, err = clientv3.New(cfg); err != nil {
		fmt.Println("connect to etcd error")
		return
	}
	fmt.Println("connect to etcd success")
	// client = client
	defer client.Close()
	// 创建一个KV操作用于对数据操作
	kv = clientv3.NewKV(client)
	if putres, err = kv.Put(context.TODO(), "job1", "golang1"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(putres.Header.Revision)
	}

}
