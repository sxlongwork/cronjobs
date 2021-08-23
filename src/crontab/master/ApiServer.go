package master

import (
	"cronjobs/src/crontab/common"
	"cronjobs/src/crontab/master/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var GOL_HTTPSERVER *ApiServer

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
	// 保存job
	if job, err = GOL_JOBMGR.SaveJob(job); err != nil {
		fmt.Println("ERROR", err)
		goto Err
	}

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
	name = req.PostForm.Get("name")

	// 删除job
	if oldJob, err = GOL_JOBMGR.DelJob(name); err != nil {
		goto Err
	}
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
	name = req.PostForm.Get("name")

	//杀死任务
	if err = GOL_JOBMGR.KillJob(name); err != nil {
		goto Err
	}
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
		serverMux  *http.ServeMux
		httpServer *http.Server
	)
	// 定义一个多路转换器，相当于路由
	serverMux = http.NewServeMux()

	// 注册模式,向路由器中注册处理函数
	serverMux.HandleFunc("/cron/job/save", handleJobSave)
	serverMux.HandleFunc("/cron/job/delete", handleJobDel)
	serverMux.HandleFunc("/cron/job/list", handleJobList)
	serverMux.HandleFunc("/cron/job/kill", handleJobKill)

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
	return
}