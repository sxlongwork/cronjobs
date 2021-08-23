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
		kv        clientv3.KV
		putRes    *clientv3.PutResponse
		delRes    *clientv3.DeleteResponse
		watcher   clientv3.Watcher
		watchChan clientv3.WatchChan
		watchRes  clientv3.WatchResponse
	)
	// 创建kv
	kv = clientv3.NewKV(cli)
	// 创建gorouting一直put, delete数据
	go func() {
		for {
			// 插入数据
			if putRes, err = kv.Put(context.TODO(), "/cron/jobs/job3", "go to ..."); err != nil {
				fmt.Println("put data error")
				return
			} else {
				fmt.Println("插入成功, Revision=", putRes.Header.Revision)
			}
			// 删除数据
			if delRes, err = kv.Delete(context.TODO(), "/cron/jobs/job3"); err != nil {
				fmt.Println("del data error")
				return
			} else {
				fmt.Println("删除成功, Revision=", delRes.Header.Revision)
			}
			time.Sleep(2 * time.Second)
		}

	}()

	// 使用watch接口监听/cron/jobs/job3的变化
	watcher = clientv3.NewWatcher(cli)
	watchChan = watcher.Watch(context.TODO(), "/cron/jobs/job3", clientv3.WithPrevKV())
	for {
		select {
		case watchRes = <-watchChan:
			if len(watchRes.Events) != 0 {
				for _, event := range watchRes.Events {
					// event.Type 事件类型，只支持put,delete
					// fmt.Println(event.Type, event.Kv.CreateRevision)
					switch event.Type {
					case clientv3.EventTypePut:
						fmt.Println("插入操作, Revision=", event.Kv.CreateRevision)
					case clientv3.EventTypeDelete:
						fmt.Println("删除操作, Revision=", event.Kv.ModRevision)
					}
				}
			}

		}
	}

}
