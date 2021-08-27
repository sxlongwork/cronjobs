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
	JobName        string        `json:"jobName"`
	Command        string        `json:"command"`
	OutPut         string        `json:"outPut"`
	Err            string        `json:"err"`
	JobPlanTime    time.Duration `json:"jobPlanTime"`
	JobSchduleTime time.Duration `josn:"jobSchduleTime"`
	JobExcuteTime  time.Duration `json:"jobExcuteTime"`
	JobEndTime     time.Duration `json:"jobEndTime"`
}

// 批量日志
type LogBatch struct {
	Logs []interface{}
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
		JobExcuteTime:  time.Duration(jobResult.StartTime.UnixNano() / 1000 / 1000),
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

// 根据key获取任务名
func GetJobSaveName(key string) (name string) {
	return strings.TrimLeft(key, JOB_SAVE_DIR)
}

// 根据key获取任务名
func GetJobKillName(key string) (name string) {
	return strings.TrimLeft(key, JOB_KILL_DIR)
}
