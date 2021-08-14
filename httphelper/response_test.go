package httphelper

import (
	"fmt"
	"github.com/wbrush/go-common/errorhandler"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJson(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		data interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "correct write test",
			args: args{
				w: httptest.NewRecorder(),
				data: struct {
					a string
				}{
					a: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Json(tt.args.w, tt.args.data)
		})
	}
}

func TestJsonError(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "correct errorhandler.Error write test",
			args: args{
				w:   httptest.NewRecorder(),
				err: errorhandler.NewError(errorhandler.ErrService),
			},
		},
		{
			name: "correct vanilla error write test",
			args: args{
				w:   httptest.NewRecorder(),
				err: fmt.Errorf("test error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			JsonError(tt.args.w, tt.args.err)
		})
	}
}
