package deferfunc_test

import (
	"context"
	"testing"

	"github.com/yushafro/effective-mobile-tz/pkg/deferfunc"
	"github.com/yushafro/effective-mobile-tz/pkg/test"
)

func TestClose(t *testing.T) {
	type args func() error
	tests := []test.TestNoExpected[args]{
		{
			Name: "success",
			Args: func() error { return nil },
		},
		{
			Name: "error",
			Args: func() error { return nil },
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			deferfunc.Close(context.Background(), tt.Args, "")
		})
	}
}
