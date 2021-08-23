package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TP struct {
	StartTime int64 `bson:"endTime"`
	EndTime   int64 `bson:"endTime"`
}
type LogRecord struct {
	JobName  string `bson:"jobName"`
	Command  string `bson:"command"`
	Content  string `bson:"content"`
	ExecTime TP     `bson:"execTime"`
}

type FindByJobName struct {
	JobName string `bson:"jobName"`
}

type DeleteByTime struct {
	StartTime LtTime `bson:"execTime.startTime"`
}
type LtTime struct {
	before int64 `bson:"$lt"`
}

func main() {
	var (
		client *mongo.Client
		err    error
		db     *mongo.Database
		mycol  *mongo.Collection
		record *LogRecord
		cursor *mongo.Cursor
		delRes *mongo.DeleteResult
	)
	//mongo admin -u admin -p EBkJJJe5rcF7 --host <host> --port <port>
	// db.createUser({user:"root",pwd:"123456",roles:[{role:"root",db:"admin"}]}) db.auth("root","123456") root/123456
	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:123456@47.93.208.134:27017").SetConnectTimeout(3*time.Second)); err != nil {
		fmt.Println("connect to mongodb error")
		return
	}
	db = client.Database("mydb")
	mycol = db.Collection("mycollection")

	// 根据jobname查找,可以限制查找个数
	findByName := &FindByJobName{JobName: "job1"}
	if cursor, err = mycol.Find(context.TODO(), findByName /*, options.Find().SetSkip(0), options.Find().SetLimit(1)*/); err != nil {
		fmt.Println("find data error")
		return
	}
	defer cursor.Close(context.TODO())
	// 遍历cursor得到数据
	for cursor.Next(context.TODO()) {
		record = &LogRecord{}
		if err := cursor.Decode(record); err != nil {
			fmt.Println("decode data faile, ERROR:", err)
			return
		}
		//打印获取的数据
		fmt.Println(record)
	}

	// 删除数据
	delByName := &DeleteByTime{
		StartTime: LtTime{
			before: time.Now().Unix(),
		},
	}
	if delRes, err = mycol.DeleteMany(context.TODO(), delByName); err != nil {
		fmt.Println("del data error")
		return
	}
	fmt.Println("删除条数", delRes.DeletedCount)

}
