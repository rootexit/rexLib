package rexMiddleware

import (
	"context"
	"encoding/json"
	"github.com/rootexit/rexLib/rexCodes"
	"github.com/rootexit/rexLib/rexSony"
	"net/http"
)

func CommonErrResponse(w http.ResponseWriter, r *http.Request, Code int32, v ...any) {

	msg := ""
	if len(v) > 0 {
		msg = v[0].(string)
	} else {
		msg = rexCodes.StatusText(Code)
	}

	requestID := ""
	if r.Context().Value("RequestID") == nil {
		requestID = rexSony.NextId()
	} else {
		requestID = r.Context().Value("RequestID").(string)
	}
	ctx := context.WithValue(r.Context(), "RequestID", requestID)
	r = r.WithContext(ctx)

	resp := rexCodes.CommonResponse{
		Code:      Code,
		Msg:       msg,
		RequestID: requestID,
		Path:      r.RequestURI,
	}
	body, _ := json.Marshal(&resp)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
