package httphelper

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"strconv"
)

type HttpHelperConfig struct {
	client HttpHelperInterface
}

var helperConfig HttpHelperConfig

//  see this for more info on timeouts: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
func getHttpClient() *http.Client {
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	return client
}

func getDefaultCfg() HttpHelperConfig {
	cfg := HttpHelperConfig{

		client: getHttpClient(),
	}

	return cfg
}

func readConfig(cfg HttpHelperConfig) HttpHelperConfig {

	return cfg
}

func init() {
	helperConfig = getDefaultCfg()
	helperConfig = readConfig(helperConfig)
}

//  this doesn't need to create the object, just return it
func Factory() (rv *HttpHelperConfig) {

	return &helperConfig
}

func (this *HttpHelperConfig) Do(req *http.Request) (response *http.Response, err error) {
	logrus.Tracef("Sending request to %s", req.URL.Path)

	response, err = this.client.Do(req)
	if err != nil || response.StatusCode == http.StatusNotFound {
		if err == nil {
			err = errors.New("HTTP Code: " + strconv.Itoa(response.StatusCode) + "; HTTP Status:" + response.Status)
		}

		return
	}

	return
}
