package config_test

import (
	"testing"

	"github.com/kujilabo/cocotola/cocotola-api-translator/src/config"
)

func TestAbc(t *testing.T) {
	tests := []struct {
		name string
		args int
		want int
	}{
		{
			"a",
			2,
			6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.Abc(tt.args); got != tt.want {
				t.Errorf("Abc() = %v, want %v", got, tt.want)
			}
		})
	}
}
