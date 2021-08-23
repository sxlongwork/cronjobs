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
		fmt.Println("connect to etcd failed.")
		return
	}
	fmt.Println("connect to etcd success.")
}

func main() {
	var (
		kv     clientv3.KV
		putRes *clientv3.PutResponse
		delRes *clientv3.DeleteResponse
	)
	kv = clientv3.NewKV(cli)
	if putRes, err = kv.Put(context.TODO(), "aa", "about"); err != nil {
		fmt.Println("put data fail.")
	} else {
		fmt.Println("Revision:", putRes.Header.Revision)
	}
	if delRes, err = kv.Delete(context.TODO(), "aa", clientv3.WithPrefix()); err != nil {
		fmt.Println("del datat fail.")
	} else {
		fmt.Println(delRes.Deleted, delRes.Header.Revision)
	}
}
