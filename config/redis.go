package config

import (
	"fmt"
	"strconv"
)

type (
	RedisParams struct {
		Host            string `json:"redis_host"`
		Port            string `json:"redis_port"`
		PortTls         string `json:"redis_port_tls"`
		User            string `json:"redis_user"`
		Password        string `json:"redis_password"`
		Database        string `json:"redis_database"`
		NumConns        string `json:"redis_conns"`
		RedisAuthString string `json:"redis_auth_string"`
		RedisCAFilePath string `json:"redis_ca_file_path"`
		IsTLSEnabled    string `json:"redis_tls_enabled"`
	}
)

func (rp RedisParams) Validate() error {
	if rp.Host == "" {
		return fmt.Errorf("redis host param cannot be empty")
	}

	if _, err := strconv.ParseUint(rp.Port, 10, 64); err != nil {
		return fmt.Errorf("wrong redis service port %s is provided. Only positive numeric values is acceptable", rp.Port)
	}

	if _, err := strconv.ParseUint(rp.Database, 10, 64); err != nil {
		return fmt.Errorf("wrong redis database %s is provided. Only positive numeric values is acceptable", rp.Database)
	}

	conns, err := strconv.Atoi(rp.NumConns)
	if err != nil {
		return err
	}
	if conns < 3 {
		return fmt.Errorf("redis Number Connections param cannot be less than 3")
	}

	return nil
}
