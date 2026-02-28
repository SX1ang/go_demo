package util

import "errors"

type JsonResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"` // omitempty表示在序列化成 JSON 时，如果字段的值是零值（例如 nil、""、0、false、空切片/空 map），自动省略该字段
}

//func JsonRsp(data interface{}) JsonResult {
//	var rsp JsonResult
//
//	switch v := data.(type) {
//	case CustomizedErr:
//		rsp.Msg = v.Error()
//		if v, ok := data.(CustomizedErr); ok {
//			rsp.Code = v.Code
//		} else {
//			rsp.Code = 10000
//		}
//	default:
//		rsp.Data = v
//		rsp.Msg = "success"
//	}
//	return rsp
//}

func JsonRsp(data any) JsonResult {
	var rsp JsonResult

	// 处理 error
	if err, ok := data.(error); ok && err != nil {
		var ce *CustomizedErr
		// 判断是否是自定义错误 CustomizedErr
		if errors.As(err, &ce) {
			rsp.Code = ce.Code
			rsp.Msg = ce.Msg
			return rsp
		}
		// 非业务错误（系统错误）
		rsp.Code = 10000
		rsp.Msg = err.Error()
		return rsp
	}

	// 2) 处理自定义error CustomizedErr
	if ce, ok := data.(*CustomizedErr); ok && ce != nil {
		rsp.Code = ce.Code
		rsp.Msg = ce.Msg
		return rsp
	}

	// 3) 正常成功返回
	rsp.Code = 200
	rsp.Msg = "success"
	rsp.Data = data
	return rsp
}

// CustomizedErr 自定义的错误类型。包含业务错误码与对应的错误信息。
type CustomizedErr struct {
	Code int
	Msg  string
}

func (e *CustomizedErr) Error() string { return e.Msg }
