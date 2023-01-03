package domain

import (
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type TranslationAddParameter interface {
	GetText() string
	GetPos() WordPos
	GetLang2() Lang2
	GetTranslated() string
}

type translationAddParameter struct {
	Text       string `validate:"required"`
	Pos        WordPos
	Lang2      Lang2
	Translated string
}

func NewTranslationAddParameter(text string, pos WordPos, lang2 Lang2, translated string) (TranslationAddParameter, error) {
	m := &translationAddParameter{
		Text:       text,
		Pos:        pos,
		Lang2:      lang2,
		Translated: translated,
	}

	return m, libD.Validator.Struct(m)
}

func (p *translationAddParameter) GetText() string {
	return p.Text
}

func (p *translationAddParameter) GetPos() WordPos {
	return p.Pos
}

func (p *translationAddParameter) GetLang2() Lang2 {
	return p.Lang2
}

func (p *translationAddParameter) GetTranslated() string {
	return p.Translated
}

type TranslationUpdateParameter interface {
	GetTranslated() string
}

type translationUpdateParameter struct {
	Translated string `validate:"required"`
}

func NewTranslationUpdateParameter(translated string) (TranslationUpdateParameter, error) {
	m := &translationUpdateParameter{
		Translated: translated,
	}

	return m, libD.Validator.Struct(m)
}

func (p *translationUpdateParameter) GetTranslated() string {
	return p.Translated
}
