package worker

import (
	"cronjobs/src/crontab/common"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

var (
	GOL_EXCUTOR *JobExcutor
)

type JobExcutor struct {
}

/*
执行任务
*/
func (excutor *JobExcutor) ExcuteJob(jobExcuteState *common.JobExecuteState) {
	var (
		cmd             *exec.Cmd
		outPut          []byte
		err             error
		jobExcuteResult *common.JobExcuteResult = &common.JobExcuteResult{}
		lock            *JobLock
	)
	jobExcuteResult.JobState = jobExcuteState

	// 为了防止1个任务被多个节点重复重复执行，需要实现分布式锁
	// 获取分布式锁
	lock = GOL_JOBMGR.CreateLock(jobExcuteState.Job.JobName)

	// 为了让各个节点均匀的执行任务，上锁前随机sleep 0-1s
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	// 尝试上锁
	err = lock.TryLock()
	defer lock.Unlock() // 解锁

	jobExcuteResult.StartTime = time.Now()
	if err != nil {
		jobExcuteResult.Err = err
		jobExcuteResult.OutPut = nil
		jobExcuteResult.EndTime = time.Now()
	} else {
		log.Printf("start to execute job %s.\n", jobExcuteState.Job.JobName)
		// 执行命令并捕获结果
		cmd = exec.CommandContext(jobExcuteState.Ctx, "/bin/bash", "-c", jobExcuteState.Job.Command)
		if outPut, err = cmd.CombinedOutput(); err != nil {
			jobExcuteResult.Err = err
			jobExcuteResult.OutPut = nil
		} else {
			jobExcuteResult.OutPut = outPut
			jobExcuteResult.Err = err
		}
		jobExcuteResult.EndTime = time.Now()
	}

	// 执行完成，推送结果
	// fmt.Println(jobExcuteResult.Job.JobName, "执行完成")
	GOL_SCHDULE.PutJobExcuteResult(jobExcuteResult)

}

func InitJobExcutor() (err error) {
	GOL_EXCUTOR = &JobExcutor{}
	return
}
