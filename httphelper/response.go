package httphelper

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/wbrush/go-common/errorhandler"
	"net/http"
)

// Json writes to ResponseWriter a single JSON-object
func Json(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		logrus.Errorf("write to ResponseWriter error: %s", err.Error())
		JsonError(w, errorhandler.NewError(errorhandler.ErrService, "ResponseWriter"))
		return
	}
}

// JsonError writes to ResponseWriter error
func JsonError(w http.ResponseWriter, err error) {
	var e *errorhandler.Error
	var ok bool

	if e, ok = err.(*errorhandler.Error); !ok {
		e = errorhandler.FromVanillaError(err)
	}

	js, _ := json.Marshal(e.ToMap())
	w.WriteHeader(e.GetHttpCode())
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		logrus.Errorf("write to ResponseWriter error: %s", err.Error())
		//JsonError(w, errorhandler.NewError(errorhandler.ErrService, "ResponseWriter"))        //  I would think this should cause a recursive loop
		return
	}
}
