package main

import "fmt"

/*

golang没有构造函数，所以一般会定义NewXXX函数来初始化相关类
NewXXX函数返回接口时就是简单工厂模式，也就是说Golang的一般推荐做法就是简单工厂

*/

// 发消息
type SendMsg interface {
	Send(msg string) string
}

type QQ struct {
	id string
}

type Wechat struct {
	name string
}

// QQ实现SendMsg
func (qq *QQ) Send(msg string) string {
	message := fmt.Sprintf("qq send message [%s]", msg)
	return message
}

// Wechat实现SendMsg
func (we *Wechat) Send(msg string) string {
	message := fmt.Sprintf("wechat send message [%s]", msg)
	return message
}

// 传入不同设备，返回不同的类
func NewSend(device string) SendMsg {
	switch device {
	case "qq":
		return &QQ{}

	case "wechat":
		return &Wechat{}
	}
	return nil
}

func main() {
	qq := NewSend("qq")
	if qq != nil {
		fmt.Println(qq.Send("I am qq"))
	}

	wechat := NewSend("wechat")
	if wechat != nil {
		fmt.Println(wechat.Send("I am wechat"))
	}
}
