package db

import (
	"fmt"
	"testing"
)

func TestCheckIfDuplicateError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is duplicate error",
			err:  fmt.Errorf("ERROR #23505 duplicate key value violates unique constraint \"some_code_uindex\""),
			want: true,
		},
		{
			name: "is no duplicate error",
			err:  fmt.Errorf("ERROR #1234 some other error"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckIfDuplicateError(tt.err); got != tt.want {
				t.Errorf("CheckIfDupplicateError() = %v, want %v", got, tt.want)
			}
		})
	}
}
