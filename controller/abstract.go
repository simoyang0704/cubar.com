package controller

import (
	"cubar.com/lib/util"
	"net/http"
)

const (
	ERROR_CODE = 404
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func packResult(w http.ResponseWriter, ret interface{}) {

	// 增加统计信息
	// statsd.Incre(callFunc + ".success")
	w.Header().Set("Content-Type", "application/json")

	// 返回结果
	w.Write([]byte(util.StructToJsonString(ret)))
}

func packError(w http.ResponseWriter, err error, code int, msg string) {

	ret := Response{
		Code: code,
		Msg:  msg,
	}

	if ret.Code == 0 {
		ret.Code = ERROR_CODE
	}

	if err != nil {
		ret.Msg = err.Error()
	}

	// 增加statsd
	// statsd.Incre(callFunc + ".error")

	// 返回错误信息
	w.Write([]byte(util.StructToJsonString(ret)))
}
