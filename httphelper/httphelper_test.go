package httphelper

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	mock "github.com/wbrush/go-common/httphelper/test"
)

func TestGetDefaultCfg(t *testing.T) {
	cfg := Factory()

	if cfg == nil || cfg.client == nil {
		t.Error("Invalid client given in httphelper")
	}
}

func TestReadConfig(t *testing.T) {
	// cfg := GetDefaultCfg()

	// // set some env vars to get different values
	// os.Setenv("ENABLE_HYSTRIX", "true")
	// os.Setenv("SYNC_HYSTRIX", "false")
	// os.Setenv("ENABLE_AUTHORIZATION", "false")
	// os.Setenv("AUTHORIZATION_TYPE", "NO_LOGIN")
	// cfg = ReadConfig(cfg)

	// if cfg.HystrixEnabled != true {
	// 	t.Error("HystrixEnabled is not settable to false")
	// }
	// if cfg.HystrixSync != false {
	// 	t.Error("HystrixSync is not settable to false")
	// }
	// if cfg.AuthorizationEnabled != false {
	// 	t.Error("HystrixEnabled is not settable to false")
	// }
	// if cfg.AuthorizationType != NO_LOGIN {
	// 	t.Error("AuthorizationType is not settable to NO_LOGIN")
	// }

	// // set some env vars to get different values
	// os.Setenv("ENABLE_HYSTRIX", "false")
	// os.Setenv("SYNC_HYSTRIX", "true")
	// os.Setenv("ENABLE_AUTHORIZATION", "true")
	// os.Setenv("AUTHORIZATION_TYPE", "ONE_LOGIN")
	// cfg = ReadConfig(cfg)

	// if cfg.HystrixEnabled != false {
	// 	t.Error("HystrixEnabled is not settable to false")
	// }
	// if cfg.HystrixSync != true {
	// 	t.Error("HystrixSync is not settable to true")
	// }
	// if cfg.AuthorizationEnabled != true {
	// 	t.Error("AuthorizationEnabled is not settable to true")
	// }
	// if cfg.AuthorizationType != ONE_LOGIN {
	// 	t.Error("AuthorizationType is not settable to ONE_LOGIN")
	// }

}

//  test basic operation
func TestDo(t *testing.T) {
	//  initialize the configuration
	testHelper := Factory()

	//  set up mocked function
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockObj := mock.NewMockHttpHelperInterface(mockCtrl)

	//  establish input and output for mocked function
	request, _ := http.NewRequest("GET", "URL", nil)
	expected_response := &http.Response{
		StatusCode: http.StatusOK,
		Status:     "Ok",
		Header:     http.Header{},
	}
	//expected_response.Header.Set("location", "http://somewhere.com/new/path/id/1234")
	mockObj.EXPECT().Do(request).Return(expected_response, nil)
	helperConfig.client = mockObj

	//  perform function & test results
	response, err := testHelper.Do(request)
	if err != nil {
		t.Error("Received Error From Do()")
	}
	if response == nil {
		t.Error("Received No Response From Do()")
	}
}
