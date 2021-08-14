package httphelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mozillazg/request"
	"github.com/sirupsen/logrus"
	"github.com/wbrush/go-common/errorhandler"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const DefaultRequestTimeout = 30

type RequestData struct {
	URL       string
	QueryData map[string][]string
	Data      map[string]string
	Json      interface{}
	Headers   map[string]string
}

func MakeHTTPRequest(r RequestData, dest interface{}, method ...string) error {
	if len(method) == 0 { //accepts only one method, else ignoring
		method = []string{"GET"} //GET is default
	}

	req := request.NewRequest(new(http.Client))
	req.Client.Timeout = time.Duration(DefaultRequestTimeout * time.Second)
	req.Data = r.Data
	req.Json = r.Json
	req.Headers = r.Headers

	fullURL := r.URL
	if len(r.QueryData) > 0 {
		fullURL += "?" + url.Values(r.QueryData).Encode()
	}

	var (
		resp *request.Response
		err  error
	)

	logrus.Tracef("sending %s request to: %s ...", method[0], fullURL)
	switch method[0] {
	case "GET":
		resp, err = req.Get(fullURL)
	case "POST":
		resp, err = req.Post(fullURL)
	case "PUT":
		resp, err = req.Put(fullURL)
	case "PATCH":
		resp, err = req.Patch(fullURL)
	case "DELETE":
		resp, err = req.Delete(fullURL)
	default:
		return fmt.Errorf("method %s is not acceptable", method[0])
	}

	if err != nil {
		return fmt.Errorf("error %s request (req.Post phrase) to %s, message: %s", method[0], r.URL, err.Error())
	}

	var content []byte
	content, err = resp.Content()
	if err != nil {
		return fmt.Errorf("error %s request (resp.Content phrase) to %s, message: %s, response: %s", method[0], r.URL, err.Error(), string(content))
	}

	if resp.StatusCode != http.StatusOK { //check if request error occurs
		var respErr errorhandler.Error
		err = json.Unmarshal(content, &respErr)
		if err != nil || !respErr.Code.IsValid() { //not an known error type, so return whole answer as is to an error
			return errors.New(string(content))
		}

		return respErr
	}

	err = json.Unmarshal(content, dest)
	if err != nil {
		return fmt.Errorf("error %s request (json.Unmarshal phrase) to %s, message: %s, response: %s", method[0], r.URL, err.Error(), string(content))
	}

	return nil
}

func MakeQueryData(input map[string]interface{}) map[string][]string {
	out := make(map[string][]string, 0)

	for key, value := range input {
		switch value.(type) {
		case int:
			out[key] = []string{strconv.Itoa(value.(int))}
		case int16:
			out[key] = []string{strconv.FormatInt(int64(value.(int16)), 10)}
		case int64:
			out[key] = []string{strconv.FormatInt(value.(int64), 10)}
		case float64:
			//TODO: change format, if need
			out[key] = []string{strconv.FormatFloat(value.(float64), 'f', 10, 10)}
		case string:
			out[key] = []string{value.(string)}
		case []string:
			out[key] = value.([]string)
		case bool:
			out[key] = []string{strconv.FormatBool(value.(bool))}
		case time.Time:
			out[key] = []string{value.(time.Time).Format("2006-01-02T15:04:05.999999999Z07:00")}
		case []interface{}:
			o := make([]string, 0)
			//this integrated switch need as data comes as array of interfaces
			for _, v := range value.([]interface{}) {
				switch v.(type) {
				case string:
					o = append(o, v.(string))

				}
			}
			out[key] = o
			//TODO add more cases
		}
	}
	return out
}
