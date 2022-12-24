package domain

import (
	"errors"
	"reflect"
	"testing"

	libD "github.com/kujilabo/cocotola/lib/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewLang2(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		args          string
		want          Lang2
		wantErr       bool
		wantErrDetail error
	}{
		{
			name:    "en",
			args:    "en",
			want:    Lang2EN,
			wantErr: false,
		},
		{
			name:    "ja",
			args:    "ja",
			want:    Lang2JA,
			wantErr: false,
		},
		{
			name:          "empty string",
			args:          "",
			wantErr:       true,
			wantErrDetail: libD.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewLang2(tt.args)
			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if tt.wantErrDetail != nil && !errors.Is(err, tt.wantErrDetail) {
					t.Errorf("NewAudioModel() err = %v, wantErrDetail %v", err, tt.wantErrDetail)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLang2() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestNewLang3(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		args    string
// 		want    Lang3
// 		wantErr bool
// 	}{
// 		{
// 			name:    "eng",
// 			args:    "eng",
// 			want:    Lang3ENG,
// 			wantErr: false,
// 		},
// 		{
// 			name:    "jpn",
// 			args:    "jpn",
// 			want:    Lang3JPN,
// 			wantErr: false,
// 		},
// 		{
// 			name:    "empty string",
// 			args:    "",
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := NewLang3(tt.args)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("NewLang3() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewLang3() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
