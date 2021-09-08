package master

import (
	"context"
	"cronjobs/src/crontab/common"
	"cronjobs/src/crontab/master/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GOL_LOGMGR *LogMgr

type LogMgr struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func InitLogMgr() (err error) {
	var (
		client *mongo.Client
		col    *mongo.Collection
	)
	if client, err = mongo.Connect(context.TODO(),
		options.Client().ApplyURI(config.GOL_CONFIG.MongodbUrl).
			SetConnectTimeout(time.Duration(config.GOL_CONFIG.MongodbTimeout)*time.Millisecond)); err != nil {

	}
	col = client.Database(config.GOL_CONFIG.MongodbName).Collection(config.GOL_CONFIG.MongodbCollectionName)
	GOL_LOGMGR = &LogMgr{
		client:     client,
		collection: col,
	}
	return
}

func (logMgr *LogMgr) FindByName(name string, start int, limit int) (logArr []*common.LogRecord, err error) {
	var (
		findByName *common.FindByJobName
		logOrder   *common.SortLogByStartTime
		cursor     *mongo.Cursor
		logRecord  *common.LogRecord
	)
	// fmt.Println(name, start, limit)
	logArr = make([]*common.LogRecord, 0)

	//  日志名称过滤参数
	findByName = &common.FindByJobName{
		JobName: name,
	}
	// 日志排序参数
	logOrder = &common.SortLogByStartTime{
		SortOrder: -1,
	}
	if cursor, err = logMgr.collection.Find(context.TODO(), findByName, options.Find().SetSort(logOrder), options.Find().SetSkip(int64(start)), options.Find().SetLimit(int64(limit))); err != nil {
		return
	}
	// 延迟释放游标
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		logRecord = &common.LogRecord{}
		if err = cursor.Decode(logRecord); err != nil {
			continue
		}

		logArr = append(logArr, logRecord)
	}

	return
}

func (logMgr *LogMgr) ClearJobLogs(name string) (err error) {
	var (
		findByName *common.FindByJobName
		delRes     *mongo.DeleteResult
	)

	//  日志名称过滤参数
	findByName = &common.FindByJobName{
		JobName: name,
	}

	if delRes, err = logMgr.collection.DeleteMany(context.TODO(), findByName); err != nil {
		log.Println("delete job logs error.", err)
		return
	}
	log.Printf("delete job %s logs success,total %d logs record.\n", name, delRes.DeletedCount)

	return
}
