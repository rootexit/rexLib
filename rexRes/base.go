package rexRes

import (
	"context"
	"github.com/aws/smithy-go/encoding/xml"
	"github.com/google/uuid"
	"github.com/rootexit/rexLib/rexCodes"
	"github.com/rootexit/rexLib/rexCommonHeader"
	"github.com/rootexit/rexLib/rexErrors"
	"github.com/rootexit/rexLib/rexMiddleware"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

// BaseResponse is the base response struct.
type BaseResponse[T any] struct {
	// Code represents the business code, not the http status code.
	Code int32 `json:"code" xml:"code"`
	// Msg represents the business message, if Code = BusinessCodeOK,
	// and Msg is empty, then the Msg will be set to BusinessMsgOk.
	Msg string `json:"msg" xml:"msg"`
	// Code represents the business code, not the http status code.
	Path string `json:"path" xml:"path"`
	// Msg represents the business message, if Code = BusinessCodeOK,
	// and Msg is empty, then the Msg will be set to BusinessMsgOk.
	RequestId string `json:"request_id" xml:"request_id"`
	// Data represents the business data.
	Data T `json:"data,omitempty" xml:"data,omitempty"`
}

type baseXmlResponse[T any] struct {
	XMLName  xml.Name `xml:"xml"`
	Version  string   `xml:"version,attr"`
	Encoding string   `xml:"encoding,attr"`
	BaseResponse[T]
}

// JsonBaseResponse writes v into w with http.StatusOK.
func JsonBaseResponse(w http.ResponseWriter, r *http.Request, res any, err any) {
	JsonBaseResponseCtx(context.Background(), w, r, res, err)
}

// JsonBaseResponseCtx writes v into w with http.StatusOK.
func JsonBaseResponseCtx(ctx context.Context, w http.ResponseWriter, r *http.Request, res any, err any) {
	if strings.Contains(w.Header().Get(ContentType), ContentTypeHtml) {
		// note: 因为大部分返回html的时候都是模板渲染，所以不需要写入
	} else {
		httpx.OkJsonCtx(ctx, w, wrapBaseResponse(ctx, r, res, err))
	}
}

// XmlBaseResponse writes v into w with http.StatusOK.
func XmlBaseResponse(w http.ResponseWriter, r *http.Request, res any, err any) {
	OkXml(w, wrapXmlBaseResponse(context.Background(), r, res, err))
}

// XmlBaseResponseCtx writes v into w with http.StatusOK.
func XmlBaseResponseCtx(ctx context.Context, w http.ResponseWriter, r *http.Request, res any, err any) {
	OkXmlCtx(ctx, w, wrapXmlBaseResponse(ctx, r, res, err))
}

func wrapXmlBaseResponse(ctx context.Context, r *http.Request, res any, err any) baseXmlResponse[any] {
	base := wrapBaseResponse(ctx, r, res, err)
	return baseXmlResponse[any]{
		Version:      xmlVersion,
		Encoding:     xmlEncoding,
		BaseResponse: base,
	}
}

func wrapBaseResponse(ctx context.Context, r *http.Request, res any, err any) BaseResponse[any] {
	path := r.URL.Path
	// note: 先从请求中获取
	requestID := ""
	xRequestIDFor := r.Header.Get(rexCommonHeader.HeaderXRequestIDFor)
	if xRequestIDFor == "" {
		// note: 再从上下文中获取
		ctxRequestId := ctx.Value(rexMiddleware.CtxRequestID)
		if ctxRequestId == nil {
			requestID = uuid.NewString()
		} else {
			requestID = ctxRequestId.(string)
		}
	} else {
		requestID = xRequestIDFor
	}
	var resp BaseResponse[any]
	if err == nil {
		resp.Code = rexCodes.EngineStatusOK
		resp.Msg = rexCodes.StatusText(rexCodes.EngineStatusOK, rexCodes.LangEnUS)
		resp.RequestId = requestID
		resp.Path = path
		resp.Data = res
	} else {
		switch data := err.(type) {
		case *rexErrors.CodeMsg:
			resp.Code = data.Code
			resp.Msg = data.Msg
			resp.RequestId = requestID
			resp.Path = path
			resp.Data = res
		case rexErrors.CodeMsg:
			resp.Code = data.Code
			resp.Msg = data.Msg
			resp.RequestId = requestID
			resp.Path = path
			resp.Data = res
		case *status.Status:
			resp.Code = int32(data.Code())
			resp.Msg = data.Message()
			resp.RequestId = requestID
			resp.Path = path
			resp.Data = res
		case interface{ GRPCStatus() *status.Status }:
			resp.Code = int32(data.GRPCStatus().Code())
			resp.Msg = data.GRPCStatus().Message()
			resp.RequestId = requestID
			resp.Path = path
			resp.Data = res
		case error:
			resp.Code = rexCodes.EngineStatusBadRequest
			resp.Msg = data.Error()
			resp.RequestId = requestID
			resp.Path = path
			resp.Data = res
		default:
			resp.Code = rexCodes.EngineStatusOK
			resp.Msg = rexCodes.StatusText(rexCodes.EngineStatusOK, rexCodes.LangEnUS)
			resp.RequestId = requestID
			resp.Path = path
			resp.Data = res
		}
	}

	return resp
}
