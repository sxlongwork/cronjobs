package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	var (
		connect *mongo.Client
		db      *mongo.Database
		col     *mongo.Collection
		err     error
	)

	// mongo admin -u admin -p 4wo3oPweXHvJ --host <host> --port <port>
	// 建立连接
	if connect, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://47.93.208.134:27017").SetConnectTimeout(3*time.Second)); err != nil {
		fmt.Println("connect to mongodb error")
		return
	}
	// connect = connect
	fmt.Println("connect to mongodb success.")

	// 选择数据库
	db = connect.Database("mydb")

	// 选择表
	col = db.Collection("myCollection")
	col = col

}
