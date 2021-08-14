package helpers

import (
	"math/rand"
	"testing"
)

func TestRandStringBytes(t *testing.T) {
	baseSrc := randSrc          //save old randSrc value
	randSrc = rand.NewSource(1) //reset rand source for test
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "5 symbols string generation test",
			args: args{n: 5},
			want: "s1phe", //for seed 1
		},
		{
			name: "more, than 10 symbols string generation test",
			args: args{n: 11},
			want: "pzXpF9GMIHD", //for seed 1
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandStringBytes(tt.args.n); got != tt.want {
				t.Errorf("RandStringBytes() = %v, want %v", got, tt.want)
			}
		})
	}

	randSrc = baseSrc //return old randSrc value
}
