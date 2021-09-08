package master

import (
	"cronjobs/src/crontab/common"
	"cronjobs/src/crontab/master/config"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var GOL_HTTPSERVER *ApiServer

/*
保存任务：新增或修改
*/
func handleJobSave(res http.ResponseWriter, req *http.Request) {

	var (
		reqData string
		job     *common.Job
		err     error
		result  []byte
	)
	// 解析post表单，请求数据
	if err = req.ParseForm(); err != nil {
		goto Err
	}
	reqData = req.PostForm.Get("job")

	//反序列化解析为job
	job = &common.Job{}
	if err = json.Unmarshal([]byte(reqData), job); err != nil {
		goto Err
	}
	// fmt.Println("save:", *job)
	// 保存job
	if err = GOL_JOBMGR.SaveJob(job); err != nil {
		log.Printf("save job %s ERROR %v.\n", job.JobName, err)
		goto Err
	}
	log.Printf("save job %s success.\n", job.JobName)

	// 构造响应
	if result, err = common.BuildResponse(200, "success", job); err == nil {
		res.Write(result)
	}
	return
Err:
	if result, err = common.BuildResponse(-1, err.Error(), job); err == nil {
		res.Write(result)
	}
}

/*
删除任务
*/
func handleJobDel(res http.ResponseWriter, req *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		result []byte
	)
	if err = req.ParseForm(); err != nil {
		goto Err
	}
	//获取job名称
	name = req.PostForm.Get("jobName")
	// fmt.Println("del:", name)

	// 删除job
	if oldJob, err = GOL_JOBMGR.DelJob(name); err != nil {
		log.Printf("delete job %s ERROR %v.\n", name, err)
		goto Err
	}
	log.Printf("delete job %s success.\n", name)

	// 构造响应
	if result, err = common.BuildResponse(200, "success", oldJob); err == nil {
		res.Write(result)
	}

	return
Err:
	if result, err = common.BuildResponse(-1, err.Error(), oldJob); err == nil {
		res.Write(result)
	}
}

/*
查询任务列表
*/
func handleJobList(res http.ResponseWriter, req *http.Request) {
	var (
		err    error
		jobs   []*common.Job
		result []byte
	)

	if jobs, err = GOL_JOBMGR.ListJob(); err != nil {
		goto Err
	}
	// 构造响应
	if result, err = common.BuildResponse(200, "success", jobs); err == nil {
		res.Write(result)
	}
	return

Err:
	if result, err = common.BuildResponse(-1, err.Error(), jobs); err == nil {
		res.Write(result)
	}
}

/*
杀死任务，终止正在执行的任务
*/
func handleJobKill(res http.ResponseWriter, req *http.Request) {
	var (
		err    error
		name   string
		result []byte
	)
	if err = req.ParseForm(); err != nil {
		goto Err
	}
	//获取job名称
	name = req.PostForm.Get("jobName")
	// fmt.Println("kill:", name)

	//杀死任务
	if err = GOL_JOBMGR.KillJob(name); err != nil {
		goto Err
	}

	log.Printf("kill job %s success.\n", name)
	// 构造响应
	if result, err = common.BuildResponse(200, "success", nil); err == nil {
		res.Write(result)
	}
	return
Err:
	// 构造响应
	if result, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		res.Write(result)
	}
}

/*
查询任务执行日志
*/
func handleJobLog(res http.ResponseWriter, req *http.Request) {
	var (
		err        error
		name       string
		limitParam string
		startParam string
		start      int
		limit      int
		logs       []*common.LogRecord
		result     []byte
	)
	if err = req.ParseForm(); err != nil {
		goto Err
	}
	name = req.Form.Get("jobName")
	startParam = req.Form.Get("start")
	limitParam = req.Form.Get("limit")
	// fmt.Println("log:", name, startParam, limitParam)

	// 字符串转换位置整形,转换出错则使用默认值
	if start, err = strconv.Atoi(startParam); err != nil {
		start = 1
	}
	if limit, err = strconv.Atoi(limitParam); err != nil {
		limit = 10
	}
	//
	if logs, err = GOL_LOGMGR.FindByName(name, start, limit); err != nil {
		goto Err
	}
	// 构造响应
	if result, err = common.BuildResponse(200, "success", logs); err == nil {
		res.Write(result)
	}
	return

Err:
	// 构造响应
	if result, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		res.Write(result)
	}
	return
}

/*
查询注册的worker节点
*/
func handleWorkerList(res http.ResponseWriter, req *http.Request) {
	var (
		err     error
		workers []string
		result  []byte
	)
	// 获取worker列表
	if workers, err = GOL_WORKERMGR.GetWorkerList(); err != nil {
		goto Err
	}
	// 获取成功，构造成功响应
	if result, err = common.BuildResponse(200, "success", workers); err == nil {
		res.Write(result)
	}
	return
Err:
	// 获取失败，构造失败响应
	if result, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		res.Write(result)
	}
}

/*
清理任务日志
*/
func handleJobLogClear(res http.ResponseWriter, req *http.Request) {
	var (
		err    error
		name   string
		result []byte
		count  int64
	)
	if err = req.ParseForm(); err != nil {
		goto Err
	}
	//获取job名称
	name = req.PostForm.Get("jobName")
	// fmt.Println("kill:", name)

	//清除任务日志
	if err, count = GOL_LOGMGR.ClearJobLogs(name); err != nil {
		goto Err
	}

	log.Printf("clear job %s logs success, total %d log records\n", name, count)
	// 构造响应
	if result, err = common.BuildResponse(200, "success", nil); err == nil {
		res.Write(result)
	}
	return
Err:
	// 构造响应
	if result, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		res.Write(result)
	}
}

func InitApiServer() (err error) {
	var (
		serverMux        *http.ServeMux
		httpServer       *http.Server
		staticDir        http.Dir
		staticFileHandle http.Handler
	)
	// 定义一个多路转换器，相当于路由
	serverMux = http.NewServeMux()

	// 注册模式,向路由器中注册处理函数
	serverMux.HandleFunc("/cron/job/save", handleJobSave)
	serverMux.HandleFunc("/cron/job/delete", handleJobDel)
	serverMux.HandleFunc("/cron/job/list", handleJobList)
	serverMux.HandleFunc("/cron/job/kill", handleJobKill)
	serverMux.HandleFunc("/cron/job/log", handleJobLog)
	serverMux.HandleFunc("/cron/worker/list", handleWorkerList)
	serverMux.HandleFunc("/cron/job/log/clear", handleJobLogClear)

	//静态文件目录
	staticDir = http.Dir("./webroot")
	staticFileHandle = http.FileServer(staticDir)
	serverMux.Handle("/", http.StripPrefix("/", staticFileHandle))

	// 定义1个http server
	httpServer = &http.Server{
		Addr:         ":" + strconv.Itoa(config.GOL_CONFIG.ServerPort),
		Handler:      serverMux,
		ReadTimeout:  time.Duration(config.GOL_CONFIG.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(config.GOL_CONFIG.WriteTimeout) * time.Millisecond,
	}

	GOL_HTTPSERVER = &ApiServer{httpServer: httpServer}
	//启动server并监听
	go func() {
		err = GOL_HTTPSERVER.httpServer.ListenAndServe()
	}()
	log.Printf("start listening server in :%d\n", config.GOL_CONFIG.ServerPort)
	return
}
