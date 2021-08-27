package worker

import (
	"context"
	"cronjobs/src/crontab/common"
	"cronjobs/src/crontab/worker/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GOL_LOGRECOEDMGR *LogMgr

type LogMgr struct {
	logRecords chan *common.LogRecord
	client     *mongo.Client
	collection *mongo.Collection
	autoCommit chan *common.LogBatch
}

func (logMgr *LogMgr) saveBatchLogs(logs *common.LogBatch) {
	// 批量写入mongodb
	logMgr.collection.InsertMany(context.TODO(), logs.Logs)
}

/*
写日志协程
*/
func (logMgr *LogMgr) scanLogRecord() {
	var (
		logRecord   *common.LogRecord
		logs        *common.LogBatch
		timer       *time.Timer
		oldLogBatch *common.LogBatch
	)
	for {
		select {
		case logRecord = <-logMgr.logRecords:
			// 一条条写太慢，优化为批量写入
			// logs为空，则为第一条日志，先初始化
			if logs == nil {
				logs = &common.LogBatch{}
				// 为了防止日志提交很慢，一直在累积，用户看不到任务执行的日志记录，这里做一个定时器，2s后若logBatch还没满就自动提交
				// timer = time.NewTimer(time.Duration(config.GOL_CONFIG.AutoCommitLogTime) * time.Millisecond)
				timer = time.AfterFunc(time.Duration(config.GOL_CONFIG.AutoCommitLogTime)*time.Millisecond, func(logs *common.LogBatch) func() {
					return func() {
						logMgr.autoCommit <- logs
					}
				}(logs))
			}
			// 加入批量切片中
			logs.Logs = append(logs.Logs, logRecord)
			// 如果切片放满则提交
			if len(logs.Logs) >= config.GOL_CONFIG.LogBatchCount {
				logMgr.saveBatchLogs(logs)
				logs = nil
				timer.Stop()
			}
		case oldLogBatch = <-logMgr.autoCommit: // 提交过期的logbatch
			if oldLogBatch != logs {
				// 如果已经提交了，则跳过
				continue
			}
			logMgr.saveBatchLogs(oldLogBatch)
			logs = nil

		}
	}
}

/*
初始化实例
*/
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
	GOL_LOGRECOEDMGR = &LogMgr{
		logRecords: make(chan *common.LogRecord, 1000),
		client:     client,
		collection: col,
		autoCommit: make(chan *common.LogBatch, 1000),
	}

	go GOL_LOGRECOEDMGR.scanLogRecord()
	return
}

/*
将日志记录放入管道中
*/
func (logMgr *LogMgr) PutJobLog(logRecord *common.LogRecord) {
	logMgr.logRecords <- logRecord
}
