package errors

import (
	"reflect"
	"testing"
)

func TestFromError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromError() = %v, want %v", got, tt.want)
			}
		})
	}
}
