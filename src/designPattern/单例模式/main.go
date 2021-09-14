package main

import (
	"fmt"
	"sync"
)

/*
使用懒惰模式的单例模式，使用双重检查加锁保证线程安全
*/
var singleton *Singleton
var once sync.Once

type Singleton struct{}

func NewSingleton() *Singleton {
	// func (o *Once) Do(f func()) Do方法当且仅当第一次被调用时才执行函数f。
	once.Do(func() {
		singleton = &Singleton{}
	})
	return singleton
}

func main() {
	single := NewSingleton()
	if single == nil {
		fmt.Println("error")
	}
}
