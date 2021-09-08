package worker

import (
	"cronjobs/src/crontab/common"
	"log"
	"time"
)

var (
	GOL_SCHDULE *SchduleJob
)

type SchduleJob struct {
	JobEvents chan *common.JobEvent
	JobPlans  map[string]*common.SchduleJobPlan
	JobStates map[string]*common.JobExecuteState
	JobResult chan *common.JobExcuteResult
}

/*
统一处理jobevent
*/
func (schdule *SchduleJob) handleJobEvent(jobevent *common.JobEvent) {
	var (
		schduleJobPlan *common.SchduleJobPlan
		err            error
		jobPlanExists  bool
		jobExcing      bool
		jobState       *common.JobExecuteState
	)
	switch jobevent.EventType {
	case common.JOB_PUT_EVENT: // 保存/修改任务事件
		// 生成任务的调度计划
		if schduleJobPlan, err = common.BuildSchduleJobPlan(jobevent.Job); err != nil {
			return
		}
		// 将任务添加到调度计划表中
		schdule.JobPlans[jobevent.Job.JobName] = schduleJobPlan
		log.Printf("add job %s to jobPlans map.\n", jobevent.Job.JobName)
	case common.JOB_DEL_EVENT: // 删除任务事件
		// 判断任务是否在调度计划表中，在的话就删除它
		if schduleJobPlan, jobPlanExists = schdule.JobPlans[jobevent.Job.JobName]; jobPlanExists {
			delete(schdule.JobPlans, jobevent.Job.JobName)
			log.Printf("delete job %s from jobPlans map.\n", jobevent.Job.JobName)
		}
	case common.JOB_KILL_EVENT: // 强杀任务事件
		// 判断任务是否在执行中，是就杀掉它
		if jobState, jobExcing = schdule.JobStates[jobevent.Job.JobName]; jobExcing {
			jobState.CancelFunc()
			log.Printf("killed job %s which is running now.\n", jobevent.Job.JobName)
			return
		}

	}
}

/*
尝试启动调度
*/
func (schdule *SchduleJob) tryStartSchdule(jobPlan *common.SchduleJobPlan) {

	var (
		jobState  *common.JobExecuteState
		jobExcing bool
	)
	// 对每一个任务生成对应的执行状态
	jobState = common.BuildJobState(*jobPlan)
	// 判断任务执行状态，任务在执行中跳过
	if _, jobExcing = schdule.JobStates[jobPlan.Job.JobName]; jobExcing {
		// fmt.Println(jobState.Job.JobName, "正在执行")
		log.Printf("%s is already running.\n", jobPlan.Job.JobName)
		return
	}
	// 如任务没有在执行中状态，就加它入之中执行map中
	schdule.JobStates[jobPlan.Job.JobName] = jobState
	// 开始执行任务
	go GOL_EXCUTOR.ExcuteJob(jobState)
	// fmt.Println("开始执行任务：", jobState.Job.JobName, jobState.JobPlanStartTime, jobState.JobRealStartTime)
}

/*
尝试调度任务，并返回距离执行下一次执行任务的间隔时间
*/
func (schdule *SchduleJob) trySchdule() time.Duration {
	var (
		jobPlan   *common.SchduleJobPlan
		curTime   time.Time
		nearTime  *time.Time
		timeAfter time.Duration
	)

	// 初始化时，任务计划表为空，设置等待1s
	if len(schdule.JobPlans) == 0 {
		timeAfter = time.Duration(1)
		return timeAfter
	}
	curTime = time.Now()
	// 遍历调度计划表所有任务
	for _, jobPlan = range schdule.JobPlans {

		if jobPlan.NextTime.Before(curTime) || jobPlan.NextTime.Equal(curTime) {
			// 尝试执行任务
			// fmt.Println("开始执行任务：", jobPlan.Job.JobName, jobPlan.Job.Command)
			schdule.tryStartSchdule(jobPlan)
			// 重新计算下次执行时间
			jobPlan.NextTime = jobPlan.Expression.Next(curTime)
		}
		// 最近一次执行任务的时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}

	}
	// 距离下一次执行任务需要等待多长时间
	timeAfter = nearTime.Sub(curTime)
	return timeAfter
}

/*
处理任务执行结果
*/
func (schdule *SchduleJob) handleJobResult(jobResult *common.JobExcuteResult) {
	var (
		jobRecord *common.LogRecord
	)
	// 任务执行完成从执行map中删除
	delete(schdule.JobStates, jobResult.JobState.Job.JobName)
	log.Println(jobResult.JobState.Job.JobName, "执行完成，执行结果=", string(jobResult.OutPut), "执行错误=", jobResult.Err.Error())
	// 将任务执行结果生成对应日志记录
	if jobResult.Err != common.TRY_LOCK_ERROR {
		jobRecord = common.BuildLogRecord(jobResult)
	}
	if jobResult.Err != nil {
		jobRecord.Err = jobResult.Err.Error()
	} else {
		jobRecord.Err = ""
	}

	// 将日志记录推送给日志协程处理，写入mongodb
	GOL_LOGRECOEDMGR.PutJobLog(jobRecord)

}

/*
循环监听jobEvent，对每一个jobEvent进行解析，生成执行计划，
*/
func (schdule *SchduleJob) scanJobsLoop() {
	var (
		jobEvent  *common.JobEvent
		waitTime  time.Duration
		timer     *time.Timer
		jobResult *common.JobExcuteResult
	)
	//初始化尝试调度任务，获取下次距离下次调度的等待时间
	waitTime = schdule.trySchdule()

	timer = time.NewTimer(waitTime)
	for {
		select {
		case jobEvent = <-schdule.JobEvents:
			// 统一处理处理jobevent事件
			schdule.handleJobEvent(jobEvent)
		case <-timer.C:
		case jobResult = <-schdule.JobResult:
			schdule.handleJobResult(jobResult)
		}

		// 重新计算距离下一次执行任务的间隔
		waitTime = schdule.trySchdule()
		// 重置timer
		timer.Reset(waitTime)
	}

}

/*
调度管理器初始化
*/
func InitSchduleMgr() (err error) {
	var (
		jobEvent  chan *common.JobEvent
		jobPlans  map[string]*common.SchduleJobPlan
		jobStarts map[string]*common.JobExecuteState
		jobResult chan *common.JobExcuteResult
	)
	jobEvent = make(chan *common.JobEvent, 1000)
	jobPlans = make(map[string]*common.SchduleJobPlan)
	jobStarts = make(map[string]*common.JobExecuteState)
	jobResult = make(chan *common.JobExcuteResult, 1000)
	GOL_SCHDULE = &SchduleJob{
		JobEvents: jobEvent,
		JobPlans:  jobPlans,
		JobStates: jobStarts,
		JobResult: jobResult,
	}
	// 启动调度协程
	go GOL_SCHDULE.scanJobsLoop()
	log.Println("job schdule server has started.")
	return
}

/*
保存任务事件(新增/修改/删除)
*/
func (schdule *SchduleJob) PutJobEvent(jobevent *common.JobEvent) {
	schdule.JobEvents <- jobevent
}

/*
保存任务执行结果
*/
func (schdule *SchduleJob) PutJobExcuteResult(jobResult *common.JobExcuteResult) {
	schdule.JobResult <- jobResult
}
