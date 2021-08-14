package config

import (
	"fmt"
	"strconv"
)

type (
	DbParams struct {
		Host     string `json:"db_host"`
		Port     string `json:"db_port"`
		User     string `json:"db_user"`
		Password string `json:"db_password"`
		Database string `json:"db_database"`
		NumConns string `json:"db_conns"`
	}
)

func (dp DbParams) Validate() error {
	if dp.Host == "" {
		return fmt.Errorf("DB Host param cannot be empty")
	}

	if _, err := strconv.ParseUint(dp.Port, 10, 64); err != nil {
		return fmt.Errorf("Wrong service port %s is provided. Only positive numeric values is acceptable", dp.Port)
	}

	if dp.User == "" {
		return fmt.Errorf("DB User param cannot be empty")
	}

	//Password doesn't need any validation

	if dp.Database == "" {
		return fmt.Errorf("DB Database param cannot be empty")
	}

	conns, err := strconv.Atoi(dp.NumConns)
	if err != nil {
		return err
	}
	if conns < 3 {
		return fmt.Errorf("DB Number Connections param cannot be less than 3")
	}

	return nil
}
