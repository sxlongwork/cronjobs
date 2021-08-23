package common

import "encoding/json"

type Job struct {
	JobName string `json:"jobName"`
	Command string `json:"command"`
	Expr    string `json:"expr"`
}

type Response struct {
	Code    int         `json:"code"`
	Meaasge string      `json:"message"`
	Data    interface{} `json:"data"`
}

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
