package config

import (
	"fmt"
	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/wbrush/go-common/helpers"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	EnvironmentTypeLocal string = "local"
	EnvironmentTypeDev   string = "dev"
	EnvironmentTypeTest  string = "test"
	EnvironmentTypeDemo  string = "demo"
	EnvironmentTypeProd  string = "prod"
)

// A constant exposing all environment types
var AllEnvironmentTypes = []string{
	EnvironmentTypeLocal,
	EnvironmentTypeDev,
	EnvironmentTypeTest,
	EnvironmentTypeProd,
}

// A constant exposing all environment types
var GCPEnvironments2 = []string{
	EnvironmentTypeLocal,
	EnvironmentTypeDev,
	EnvironmentTypeTest,
	EnvironmentTypeProd,
}

// these environments are really running in the cloud
// this marker will be used to make decisions in certain cases
// can be updated if we decide to run more environments in the future
var GCPEnvironments = []string{
	EnvironmentTypeDev,
	EnvironmentTypeTest,
	EnvironmentTypeProd,
}

type ServiceParams struct {
	Version string
	BuiltAt string

	Environment  string `json:"environment"`
	GlobalRegion string `json:"global_region"`

	Host     string `json:"host"`
	BaseUri  string `json:"base_uri"`
	Port     string `json:"port"`
	LogLevel string `json:"log_level"`

	isLoaded bool
}

/*
scanConfigVars scans all ENV vars recursively to found proper values to fill,
using given in "json"-tag names (UPPERCASED)
*/
func (sp *ServiceParams) scanConfigVars(target interface{}) (output map[string]interface{}) {
	output = make(map[string]interface{}, 0)
	s := structs.New(target)

	for _, f := range s.Fields() {
		if f.IsExported() {
			switch f.Kind() {
			case reflect.Struct:
				output[f.Name()] = sp.scanConfigVars(f.Value())
			default:

				envVal, isPresent := os.LookupEnv(helpers.UndescoreUppercased(f.Tag("json")))
				if isPresent {
					output[f.Name()] = envVal
				}
			}
		}
	}

	return output
}

func (sp *ServiceParams) LoadEnvVariables(target interface{}, commit, builtAt string) error {
	sp.Version = commit
	sp.BuiltAt = builtAt

	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("target is not a pointer")
	}
	if reflect.TypeOf(target).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target is not a pointer to struct")
	}

	configsMap := sp.scanConfigVars(target)
	err := mapstructure.Decode(configsMap, target)
	if err != nil {
		return fmt.Errorf("can't decode config to struct: %s", err.Error())
	}

	// set isLoaded to true when finished loading
	sp.isLoaded = true

	return nil
}

func (sp ServiceParams) IsLoaded() bool {
	return sp.isLoaded
}

func (sp ServiceParams) ConfigureLogger(serviceName ...string) error {
	//configure logrus
	// will find out if running in real environment (dev/test/prod/etc) and setup stackdriver formatter
	// else, skip to the basics
	thisServiceName := "service-name-undefined"

	if len(serviceName) > 0 {
		thisServiceName = serviceName[0]
	}

	if _, found := Find(GCPEnvironments, sp.Environment); found {
		formatter := stackdriver.NewFormatter(
			stackdriver.WithService(thisServiceName),
			stackdriver.WithVersion(sp.Version),
		)
		logrus.SetFormatter(formatter)
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05.000000",
		})
	}

	//set the debug level from config, if it present
	logrus.SetLevel(logrus.InfoLevel)
	if sp.LogLevel != "" {
		if lvl, err := logrus.ParseLevel(sp.LogLevel); err == nil {
			logrus.Infof("Log level was set to %s", sp.LogLevel)
			logrus.SetLevel(lvl)
		} else {
			return fmt.Errorf("log level %s is wrong. Available options are: %v", sp.LogLevel, logrus.AllLevels)
		}
	}

	return nil
}

func IsEnviromentValid(env string) bool {
	s := strings.ToLower(env)

	if strings.HasPrefix(env, "dev") {
		return true
	}

	switch s {
	case EnvironmentTypeLocal, EnvironmentTypeProd, EnvironmentTypeDev, EnvironmentTypeTest, EnvironmentTypeDemo:
		return true
	default:
		return false
	}
}

func (sp ServiceParams) Validate() error {
	if !IsEnviromentValid(sp.Environment) {
		return fmt.Errorf("wrong environment %s. Available options are: %v", sp.Environment, AllEnvironmentTypes)
	}

	//sp.Host doesn't need any validation

	if sp.Port != "" { //port is provided
		if _, err := strconv.ParseUint(sp.Port, 10, 64); err != nil {
			return fmt.Errorf("wrong service port %s is provied. Only positive numeric values is acceptable", sp.Port)
		}
	}

	//validate LogLevel
	if sp.LogLevel != "" { //log level is provided
		if _, err := logrus.ParseLevel(sp.LogLevel); err != nil {
			return fmt.Errorf("log level %s is wrong. Available options are: %v", sp.LogLevel, logrus.AllLevels)
		}
	}

	return nil
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
