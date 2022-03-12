package writer

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/firdasafridi/parkinglot/lib/common/commonerr"
	"github.com/firdasafridi/parkinglot/lib/common/log"
)

var (
	ok = []byte(`{"data":"ok"}`)
)

// WriteOK is a helper around Response.Write with status OK
func WriteOK(ctx context.Context, w http.ResponseWriter, data interface{}) {
	write(ctx, w, data, http.StatusOK)
}

// WriteStrOK is a helper around Response.Write with status OK
func WriteStrOK(ctx context.Context, w http.ResponseWriter) {
	set(ctx, w, ok, http.StatusOK)
}

// WriteJSONAPIError is a helper
func WriteJSONAPIError(ctx context.Context, w http.ResponseWriter, err error) {
	switch errCause := errors.Cause(err).(type) {
	case *commonerr.ErrorMessage:
		write(ctx, w, errCause, errCause.Code)
	default:
		write(ctx, w, commonerr.ErrorMessage{
			ErrorList: commonerr.SetNewInternalError().GetListError(),
		}, http.StatusInternalServerError)
	}
}

func write(ctx context.Context, w http.ResponseWriter, data interface{}, status int) {
	datab, err := json.Marshal(data)
	if err != nil {
		datab = []byte(`{"error_list":[{"error_name": "internal", "error_description": "Internal Server Error"}]}`)
	}
	set(ctx, w, datab, status)
}

func set(ctx context.Context, w http.ResponseWriter, datab []byte, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(datab)
	if err != nil {
		log.Errorln("[HTTP]", err)
	}
}

type Format struct {
	Data interface{} `json:"data"`
}

// SetOKWithData set http ok status
func SetOKWithData(ctx context.Context, w http.ResponseWriter, data interface{}) {
	datab, _ := json.Marshal(Format{
		Data: data,
	})
	set(ctx, w, datab, http.StatusOK)
}
