//go:generate mockery --output mock --name AudioModel
package domain

import libD "github.com/kujilabo/cocotola/lib/domain"

type AudioID uint

type AudioModel interface {
	GetID() uint
	GetLang5() Lang5
	GetText() string
	GetContent() string
}

type audioModel struct {
	ID      uint `validate:"required"`
	Lang5   Lang5
	Text    string `validate:"required"`
	Content string `validate:"required"`
}

func NewAudioModel(id uint, lang5 Lang5, text, content string) (AudioModel, error) {
	m := &audioModel{
		ID:      id,
		Lang5:   lang5,
		Text:    text,
		Content: content,
	}

	return m, libD.Validator.Struct(m)
}

func (a *audioModel) GetID() uint {
	return a.ID
}

func (a *audioModel) GetLang5() Lang5 {
	return a.Lang5
}

func (a *audioModel) GetText() string {
	return a.Text
}

func (a *audioModel) GetContent() string {
	return a.Content
}
