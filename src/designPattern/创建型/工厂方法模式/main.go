package main

import "fmt"

/*

工厂方法模式使用子类的方法延迟生成对象到子类中实现。
Go中不存在继承，所以使用匿名组合来实现

*/
type Device struct {
	name string
}
type SendMsg interface {
	SetName(name string)
	Send(msg string)
}

type DeviceFactory interface {
	Create() Device
}

type QQ struct {
	Device
}

type Wechat struct {
	Device
}

type QQFactory struct{}
type WechatFactory struct{}

func (qq *QQ) SetName(name string) {
	qq.name = name
}

func (wechat *Wechat) SetName(name string) {
	wechat.name = name
}

func (qq *QQ) Send(msg string) {
	fmt.Printf("%s send message [%s]\n", qq.name, msg)
}
func (we *Wechat) Send(msg string) {
	fmt.Printf("%s send message [%s]\n", we.name, msg)
}

func (qqSend *QQFactory) Create() SendMsg {
	return &QQ{}
}
func (weSend *WechatFactory) Create() SendMsg {
	return &Wechat{}
}

func main() {
	qqSend := &QQFactory{}
	qq := qqSend.Create()
	qq.SetName("qqA")
	qq.Send("I am qqA")

	weSend := &WechatFactory{}
	wechat := weSend.Create()
	wechat.SetName("wechatA")
	wechat.Send("I am wechatA")

}
