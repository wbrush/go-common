package config

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestEnvironmentType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		env  EnvironmentType
		want bool
	}{
		{
			name: "valid env case",
			env:  EnvironmentTypeLocal,
			want: true,
		},
		{
			name: "wrong env case",
			env:  EnvironmentType("wrong"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.env.IsValid(); got != tt.want {
				t.Errorf("EnvironmentType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceParams_LoadEnvVariables(t *testing.T) {
	type args struct {
		commit  string
		builtAt string
	}
	tests := []struct {
		name    string
		sp      ServiceParams
		args    args
		wantErr bool
	}{
		{
			name: "Env loading",
			sp:   ServiceParams{},
			args: args{
				commit:  "12345abcd",
				builtAt: "12321323",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sp.LoadEnvVariables(&tt.sp, tt.args.commit, tt.args.builtAt); (err != nil) != tt.wantErr {
				t.Errorf("ServiceParams.LoadEnvVariables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceParams_IsLoaded(t *testing.T) {
	tests := []struct {
		name string
		sp   ServiceParams
		want bool
	}{
		{
			name: "config loaded case",
			sp: ServiceParams{
				isLoaded: true,
			},
			want: true,
		},
		{
			name: "config is not loaded case",
			sp: ServiceParams{
				isLoaded: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sp.IsLoaded(); got != tt.want {
				t.Errorf("ServiceParams.IsLoaded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceParams_ConfigureLogger(t *testing.T) {
	tests := []struct {
		name         string
		sp           ServiceParams
		wantLogLevel logrus.Level
		wantErr      bool
	}{
		{
			name: "proper loglevel",
			sp: ServiceParams{
				LogLevel: "debug",
			},
			wantLogLevel: logrus.DebugLevel,
			wantErr:      false,
		},
		{
			name:         "no loglevel",
			sp:           ServiceParams{},
			wantLogLevel: logrus.InfoLevel,
			wantErr:      false,
		},
		{
			name: "bad loglevel",
			sp: ServiceParams{
				LogLevel: "bad",
			},
			wantLogLevel: logrus.InfoLevel,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sp.ConfigureLogger(); (err != nil) != tt.wantErr {
				t.Errorf("ServiceParams.ConfigureLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
			if logrus.GetLevel() != tt.wantLogLevel {
				t.Errorf("ServiceParams.ConfigureLogger() wrong level = %v, wantLogLevel %v", logrus.GetLevel(), tt.wantLogLevel)
			}
		})
	}
}

func TestServiceParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sp      ServiceParams
		wantErr bool
	}{
		{
			name: "Validation is ok case",
			sp: ServiceParams{
				Port:        "1234",
				LogLevel:    "debug",
				Environment: EnvironmentTypeLocal,
			},
			wantErr: false,
		},
		{
			name: "Wrong environment case",
			sp: ServiceParams{
				Environment: EnvironmentType("wrong"),
			},
			wantErr: true,
		},
		{
			name: "Wrong port case",
			sp: ServiceParams{
				Environment: EnvironmentTypeLocal,
				Port:        "wrong",
			},
			wantErr: true,
		},
		{
			name: "Wrong log level case",
			sp: ServiceParams{
				Environment: EnvironmentTypeLocal,
				LogLevel:    "wrong",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sp.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ServiceParams.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
