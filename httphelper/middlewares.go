package httphelper

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func MWUserInfoHeader(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//  do stuff here if needed

	next.ServeHTTP(w, r.WithContext(ctx))
	return
}
