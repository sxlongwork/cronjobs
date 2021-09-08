package common

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
)

/*
任务信息
*/
type Job struct {
	JobName string `json:"jobName"`
	Command string `json:"command"`
	Expr    string `json:"expr"`
}

/*
响应结构
*/
type Response struct {
	Code    int         `json:"code"`
	Meaasge string      `json:"message"`
	Data    interface{} `json:"data"`
}

/*
put delete kill job 事件
*/
type JobEvent struct {
	EventType int  `json:"eventType"`
	Job       *Job `json:"job"`
}

/*
job任务执行计划
*/
type SchduleJobPlan struct {
	Job        *Job                 `json:"job"`
	Expression *cronexpr.Expression `json:"expression"`
	NextTime   time.Time            `json:"nextTime"`
}

/*
job任务执行状态
*/
type JobExecuteState struct {
	Job              *Job
	JobPlanStartTime time.Time
	JobRealStartTime time.Time
	Ctx              context.Context
	CancelFunc       context.CancelFunc
}

/*
任务执行结果
*/
type JobExcuteResult struct {
	JobState  *JobExecuteState
	OutPut    []byte
	Err       error
	StartTime time.Time
	EndTime   time.Time
}

/*
日志结构
*/
type LogRecord struct {
	JobName        string        `json:"jobName" bson:"jobName"`
	Command        string        `json:"command" bson:"command"`
	OutPut         string        `json:"outPut" bson:"outPut"`
	Err            string        `json:"err" bson:"err"`
	JobPlanTime    time.Duration `json:"jobPlanTime" bson:"jobPlanTime"`
	JobSchduleTime time.Duration `json:"jobSchduleTime" bson:"jobSchduleTime"`
	JobStartTime   time.Duration `json:"jobStartTime" bson:"jobStartTime"`
	JobEndTime     time.Duration `json:"jobEndTime" bson:"jobEndTime"`
}

// 批量日志
type LogBatch struct {
	Logs []interface{}
}

// 查询日志过滤参数
type FindByJobName struct {
	JobName string `bson:"jobName"`
}

// 查询日志排序参数
type SortLogByStartTime struct {
	SortOrder int `bson:"jobStartTime"`
}

/*
构建日志记录结构体
*/
func BuildLogRecord(jobResult *JobExcuteResult) (jobRecord *LogRecord) {
	jobRecord = &LogRecord{
		JobName:        jobResult.JobState.Job.JobName,
		Command:        jobResult.JobState.Job.Command,
		OutPut:         string(jobResult.OutPut),
		JobPlanTime:    time.Duration(jobResult.JobState.JobPlanStartTime.UnixNano() / 1000 / 1000),
		JobSchduleTime: time.Duration(jobResult.JobState.JobRealStartTime.UnixNano() / 1000 / 1000),
		JobStartTime:   time.Duration(jobResult.StartTime.UnixNano() / 1000 / 1000),
		JobEndTime:     time.Duration(jobResult.EndTime.UnixNano() / 1000 / 1000),
	}
	return
}

/*
构建job任务执行状态
*/
func BuildJobState(jobplan SchduleJobPlan) (jobState *JobExecuteState) {
	jobState = &JobExecuteState{
		Job:              jobplan.Job,
		JobPlanStartTime: jobplan.NextTime,
		JobRealStartTime: time.Now(),
	}
	jobState.Ctx, jobState.CancelFunc = context.WithCancel(context.TODO())
	return

}

// 构建调度任务计划
func BuildSchduleJobPlan(job *Job) (schduleJobPlan *SchduleJobPlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	if expr, err = cronexpr.Parse(job.Expr); err != nil {
		return
	}
	schduleJobPlan = &SchduleJobPlan{
		Job:        job,
		Expression: expr,
		NextTime:   expr.Next(time.Now()),
	}
	return
}

// 构建响应
func BuildResponse(code int, msg string, data interface{}) (result []byte, err error) {

	var (
		res Response
	)
	res = Response{
		Code:    code,
		Meaasge: msg,
		Data:    data,
	}
	if result, err = json.Marshal(res); err != nil {
		return
	}
	return
}

// watch 事件
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

// 反序列化
func Unmarshal(bytes []byte) (job *Job, err error) {
	job = &Job{}
	if err = json.Unmarshal(bytes, job); err != nil {
		return
	}
	return
}

// 根据key获取去掉指定前缀
func GetSuffixName(key string, prefix string) (name string) {
	return strings.TrimPrefix(key, prefix)
}
