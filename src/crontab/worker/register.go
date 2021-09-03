package worker

import (
	"context"
	"cronjobs/src/crontab/common"
	"cronjobs/src/crontab/worker/config"
	"net"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var GOL_REGISTER *RegisterWorker

type RegisterWorker struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIp string
}

func InitRegister() (err error) {
	var (
		cfg     clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		localIp string
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
	if localIp, err = getLocalIp(); err != nil {
		return
	}
	GOL_REGISTER = &RegisterWorker{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIp: localIp,
	}
	// 启动注册
	go GOL_REGISTER.register()
	return
}

func (regis *RegisterWorker) register() (err error) {
	var (
		registerKey   string
		grantRes      *clientv3.LeaseGrantResponse
		keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
		keepRes       *clientv3.LeaseKeepAliveResponse
		ctx           context.Context
		cancelFunc    context.CancelFunc
	)
	cancelFunc = nil
	registerKey = common.WORKER_REGISTER_DIR + regis.localIp

	for {
		// 申请租约
		if grantRes, err = regis.lease.Grant(context.TODO(), 10); err != nil {
			// 重新尝试
			goto RETRY
		}
		// 自动续租
		if keepAliveChan, err = regis.lease.KeepAlive(context.TODO(), grantRes.ID); err != nil {
			// 重新尝试
			goto RETRY
		}
		ctx, cancelFunc = context.WithCancel(context.TODO())
		// 注册ip
		if _, err = regis.kv.Put(ctx, registerKey, "", clientv3.WithLease(grantRes.ID)); err != nil {
			goto RETRY
		}
		// 处理续租应答
		for {
			select {
			case keepRes = <-keepAliveChan:
				if keepRes == nil {
					goto RETRY
				}
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			cancelFunc()
		}
	}

}

func getLocalIp() (ipv4 string, err error) {

	var (
		ipNets  []net.Addr
		netAddr net.Addr
		ipNet   *net.IPNet
		isIpNet bool
	)
	//获取所有网卡列表地址
	if ipNets, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 遍历网卡地址列表，需要排除本地回环地址和ipv6地址
	for _, netAddr = range ipNets {
		// 类型转换为ip地址(ipv4或ipv6),并且不是回环地址
		if ipNet, isIpNet = netAddr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过ipv6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.String() // 192.168.10.10
				return
			}
		}
	}
	err = common.NOTFOUND_LOCAL_IP
	return
}
