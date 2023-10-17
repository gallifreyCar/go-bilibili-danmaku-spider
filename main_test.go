package main

import "testing"

func Test_secondsToMinuteSecond(t *testing.T) {
	type args struct {
		seconds float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{seconds: 91.35700},
			want: "00:01:31.36",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := secondsToHourMinuteSecond(tt.args.seconds); got != tt.want {
				t.Errorf("secondsToMinuteSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}
