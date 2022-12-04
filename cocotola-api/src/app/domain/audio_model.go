//go:generate mockery --output mock --name AudioModel
package domain

import libD "github.com/kujilabo/cocotola/lib/domain"

type AudioID uint

type AudioModel interface {
	GetID() uint
	GetLang2() Lang2
	GetText() string
	GetContent() string
}

type audioModel struct {
	ID      uint   `validate:"required"`
	Lang2   Lang2  `validate:"required"`
	Text    string `validate:"required"`
	Content string `validate:"required"`
}

func NewAudioModel(id uint, lang2 Lang2, text, content string) (AudioModel, error) {
	m := &audioModel{
		ID:      id,
		Lang2:   lang2,
		Text:    text,
		Content: content,
	}

	return m, libD.Validator.Struct(m)
}

func (m *audioModel) GetID() uint {
	return m.ID
}

func (m *audioModel) GetLang2() Lang2 {
	return m.Lang2
}

func (m *audioModel) GetText() string {
	return m.Text
}

func (m *audioModel) GetContent() string {
	return m.Content
}
