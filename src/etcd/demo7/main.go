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
		kv    clientv3.KV
		op    clientv3.Op
		opRes clientv3.OpResponse
	)
	kv = clientv3.NewKV(cli)

	// put data
	op = clientv3.OpPut("job4", "op operation")
	if opRes, err = kv.Do(context.TODO(), op); err != nil {
		fmt.Println("put data error")
		return
	} else {
		fmt.Println("put success, Revision=", opRes.Put().Header.Revision)
	}
	// get data
	op = clientv3.OpGet("job4")
	if opRes, err = kv.Do(context.TODO(), op); err != nil {
		fmt.Println("get data error")
		return
	} else {
		fmt.Println("get success, Value=", string(opRes.Get().Kvs[0].Value))
	}
	// delete data
	op = clientv3.OpDelete("job4")
	if opRes, err = kv.Do(context.TODO(), op); err != nil {
		fmt.Println("del data error")
		return
	} else {
		fmt.Println("del success, Revision=", opRes.Del().Header.Revision)
	}
}
