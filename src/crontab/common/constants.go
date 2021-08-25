package common

import "errors"

const (
	JOB_SAVE_DIR string = "/cron/job/"
	JOB_KILL_DIR string = "/cron/kill/"

	// job事件类型
	JOB_PUT_EVENT  int = 1
	JOB_DEL_EVENT  int = 2
	JOB_KILL_EVENT int = 3

	// 任务锁前缀
	JOB_LOCK_PREFIX string = "/cron/lock/"
)

// 上锁失败错误
var TRY_LOCK_ERROR = errors.New("锁已被占用")
