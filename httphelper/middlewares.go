package httphelper

import (
	"context"
	"net/http"
)

func MWUserInfoHeader(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//  do stuff here if needed

	next.ServeHTTP(w, r.WithContext(context.Background()))
	return
}
