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
		fmt.Println("connect to etcd failed")
		return
	}
	fmt.Println("connect to etcd success")
}

func main() {
	// 创建kv，操作数据
	var (
		kv     clientv3.KV
		putRes *clientv3.PutResponse
		getRes *clientv3.GetResponse
	)
	kv = clientv3.NewKV(cli)
	// put存放数据，并获取当前
	if putRes, err = kv.Put(context.TODO(), "job2", "java", clientv3.WithPrevKV()); err != nil {
		fmt.Println("put data failed")
	} else {
		// fmt.Println(string(putRes.PrevKv.Value))
		fmt.Println(putRes.Header.Revision)
		if putRes.PrevKv != nil {
			fmt.Println(putRes)
		}
	}
	// 获取job为前缀的所有键值对
	if getRes, err = kv.Get(context.TODO(), "job", clientv3.WithPrefix()); err != nil {
		fmt.Println("get data failed.")
	} else {
		fmt.Println(getRes.Header.Revision)
		for _, v := range getRes.Kvs {
			fmt.Println(string(v.Key), string(v.Value))
		}
	}
}
