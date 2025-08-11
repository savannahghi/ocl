package ocl

import (
	"net/http"
	"testing"
)

func TestIsDuplicateMnemonicError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "Duplicate mnemonic error",
			err: &APIError{
				StatusCode: http.StatusBadRequest,
				Mnemonic:   "already exists",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDuplicateMnemonicError(tt.err); got != tt.want {
				t.Errorf("IsDuplicateMnemonicError() = %v, want %v", got, tt.want)
			}
		})
	}
}
