package httphelper

import "net/http"

//go:generate mockgen -destination ./test/mock_httphelper.go -package mock_httphelper github.com/KWRI/go-common/httphelper HttpHelperInterface

type HttpHelperInterface interface {
	Do(req *http.Request) (response *http.Response, err error)
}
