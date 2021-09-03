package common

import "errors"

// 上锁失败错误
var TRY_LOCK_ERROR = errors.New("锁已被占用")

// 没有获取到本地ip
var NOTFOUND_LOCAL_IP = errors.New("没有找到本地IP")
