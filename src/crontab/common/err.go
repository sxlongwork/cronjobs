package common

import "errors"

// 上锁失败错误
var TRY_LOCK_ERROR = errors.New("锁已被占用")
