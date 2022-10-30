package errors

import "testing"

func TestAbc(t *testing.T) {
	tests := []struct {
		name string
		args int
		want int
	}{
		{
			"a",
			2,
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abc(tt.args); got != tt.want {
				t.Errorf("Abc() = %v, want %v", got, tt.want)
			}
		})
	}
}
