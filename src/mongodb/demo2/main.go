package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TP struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}
type LogRecord struct {
	JobName  string `bson:"jobName"`
	Command  string `bson:"command"`
	Content  string `bson:"content"`
	ExecTime TP     `bson:"execTime"`
}

func main() {
	var (
		client        *mongo.Client
		err           error
		db            *mongo.Database
		mycol         *mongo.Collection
		record        *LogRecord
		insertOneRes  *mongo.InsertOneResult
		insertManyRes *mongo.InsertManyResult
		records       []interface{}
	)
	//mongo admin -u admin -p EBkJJJe5rcF7 --host <host> --port <port>
	// db.createUser({user:"root",pwd:"123456",roles:[{role:"root",db:"admin"}]}) db.auth("root","123456") root/123456
	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:123456@47.93.208.134:27017").SetConnectTimeout(3*time.Second)); err != nil {
		fmt.Println("connect to mongodb error")
		return
	}
	db = client.Database("mydb")
	mycol = db.Collection("mycollection")

	record = &LogRecord{
		JobName: "job1",
		Command: "echo it is job1",
		Content: "it is job1",
		ExecTime: TP{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}
	// 插入一条数据
	if insertOneRes, err = mycol.InsertOne(context.TODO(), record); err != nil {
		fmt.Println("insert data fail.", err)
		return
	}
	fmt.Println(insertOneRes.InsertedID)
	// 插入多条数据
	records = []interface{}{*record, *record, *record}
	if insertManyRes, err = mycol.InsertMany(context.TODO(), records); err != nil {
		fmt.Println("insert many data error")
		return
	}
	for _, value := range insertManyRes.InsertedIDs {
		id := value.(primitive.ObjectID)
		fmt.Println(id.Hex())
	}

}
