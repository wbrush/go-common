package errorhandler

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var testErr = errors.New("some error message")

func TestError_Error(t *testing.T) {
	type fields struct {
		error       error
		Code        ErrorCode
		Value       string
		Description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "show proper error message",
			fields: fields{
				error: testErr,
				Code:  ErrService,
			},
			want: fmt.Sprintf("[%s] %s", string(ErrService), testErr.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				error:       tt.fields.error,
				Code:        tt.fields.Code,
				Value:       tt.fields.Value,
				Description: tt.fields.Description,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_ErrorCode(t *testing.T) {
	type fields struct {
		error       error
		Code        ErrorCode
		Value       string
		Description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "return code in string representation",
			fields: fields{
				Code: ErrNotFound,
			},
			want: string(ErrNotFound),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				error:       tt.fields.error,
				Code:        tt.fields.Code,
				Value:       tt.fields.Value,
				Description: tt.fields.Description,
			}
			if got := e.ErrorCode(); got != tt.want {
				t.Errorf("Error.ErrorCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_ToMap(t *testing.T) {
	type fields struct {
		error       error
		Code        ErrorCode
		Value       string
		Description string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "return map representation",
			fields: fields{
				Code:        ErrNotFound,
				Value:       "testValue",
				Description: "desc",
			},
			want: map[string]interface{}{
				"code":        string(ErrNotFound),
				"value":       "testValue",
				"description": "desc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				error:       tt.fields.error,
				Code:        tt.fields.Code,
				Value:       tt.fields.Value,
				Description: tt.fields.Description,
			}
			if got := e.ToMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		code  ErrorCode
		value []string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "no value was provided case",
			args: args{
				code: ErrNotFound,
			},
			want: &Error{Code: ErrNotFound},
		},
		{
			name: "value was provided case",
			args: args{
				code:  ErrNotFound,
				value: []string{"test"},
			},
			want: &Error{Code: ErrNotFound, Value: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.code, tt.args.value...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewErrorWithDesc(t *testing.T) {
	type args struct {
		code  ErrorCode
		desc  string
		value []string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "no value was provided case",
			args: args{
				code: ErrNotFound,
				desc: "test desc",
			},
			want: &Error{Code: ErrNotFound, Description: "test desc"},
		},
		{
			name: "value was provided case",
			args: args{
				code:  ErrNotFound,
				desc:  "test desc",
				value: []string{"test"},
			},
			want: &Error{Code: ErrNotFound, Description: "test desc", Value: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErrorWithDesc(tt.args.code, tt.args.desc, tt.args.value...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithDesc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromVanillaError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "service error creation from golang error",
			args: args{
				err: testErr,
			},
			want: &Error{
				Code:  ErrService,
				Value: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromVanillaError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_GetHttpCode(t *testing.T) {
	tests := []struct {
		name   string
		fields *Error
		want   int
	}{
		{
			name:   "service error case",
			fields: NewError(ErrService),
			want:   http.StatusInternalServerError,
		},
		{
			name:   "other error case",
			fields: NewError(ErrBadParam),
			want:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.GetHttpCode(); got != tt.want {
				t.Errorf("Error.GetHttpCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
