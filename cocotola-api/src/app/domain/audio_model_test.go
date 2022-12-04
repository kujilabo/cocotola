package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAudio(t *testing.T) {
	type args struct {
		id           uint
		lang2        Lang2
		text         string
		audioContent string
	}
	tests := []struct {
		name          string
		args          args
		wantID        uint
		wantLang      Lang2
		wantText      string
		wantContent   string
		wantErr       bool
		wantErrDetail error
	}{
		{
			name: "valid",
			args: args{
				id:           1,
				lang2:        Lang2EN,
				text:         "Hello",
				audioContent: "HELLO_CONTENT",
			},
			wantID:      1,
			wantLang:    Lang2EN,
			wantText:    "Hello",
			wantContent: "HELLO_CONTENT",
			wantErr:     false,
		},
		{
			name: "text is empty",
			args: args{
				id:           1,
				lang2:        Lang2EN,
				text:         "",
				audioContent: "HELLO_CONTENT",
			},
			wantErr: true,
		},
		{
			name: "content is empty",
			args: args{
				id:           1,
				lang2:        Lang2EN,
				text:         "Hello",
				audioContent: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAudioModel(tt.args.id, tt.args.lang2, tt.args.text, tt.args.audioContent)
			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, got.GetID())
				assert.Equal(t, tt.wantLang, got.GetLang2())
				assert.Equal(t, tt.wantText, got.GetText())
				assert.Equal(t, tt.wantContent, got.GetContent())
			} else {
				assert.Error(t, err)
				if tt.wantErrDetail != nil && !errors.Is(err, tt.wantErrDetail) {
					t.Errorf("NewAudioModel() err = %v, wantErrDetail %v", err, tt.wantErrDetail)
				}
			}
		})
	}
}
