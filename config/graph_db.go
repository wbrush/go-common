package config

import (
	"fmt"
	"strconv"
)

type (
	GraphDbParams struct {
		Host     string `json:"graphdb_host"`
		Port     string `json:"graphdb_port"`
		User     string `json:"graphdb_user"`
		Password string `json:"graphdb_password"`
	}
)

func (gDBParams GraphDbParams) Validate() error {
	if gDBParams.Host == "" {
		return fmt.Errorf("GraphDB Host param cannot be empty")
	}

	if _, err := strconv.ParseUint(gDBParams.Port, 10, 64); err != nil {
		return fmt.Errorf("Wrong service port %s is provided. Only positive numeric values is acceptable", gDBParams.Port)
	}

	if gDBParams.User == "" {
		return fmt.Errorf("GraphDB User param cannot be empty")
	}
	return nil
}
