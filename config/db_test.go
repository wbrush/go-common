package config

import "testing"

func TestDbParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dp      DbParams
		wantErr bool
	}{
		{
			name: "Validation is ok case",
			dp: DbParams{
				Host:     "somehost",
				Port:     "1234",
				User:     "user",
				Database: "db",
			},
			wantErr: false,
		},
		{
			name: "Wong host",
			dp: DbParams{
				Host:     "",
				Port:     "1234",
				User:     "user",
				Database: "db",
			},
			wantErr: true,
		},
		{
			name: "Wrong port",
			dp: DbParams{
				Host:     "somehost",
				Port:     "wrong",
				User:     "user",
				Database: "db",
			},
			wantErr: true,
		},
		{
			name: "Wrong user",
			dp: DbParams{
				Host:     "somehost",
				Port:     "1234",
				User:     "",
				Database: "db",
			},
			wantErr: true,
		},
		{
			name: "Wrong db",
			dp: DbParams{
				Host:     "somehost",
				Port:     "1234",
				User:     "user",
				Database: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dp.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DbParams.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
