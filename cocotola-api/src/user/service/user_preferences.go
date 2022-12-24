package service

import "context"

type UserPreferences interface {
}

type userPreferences struct {
}

func NewUserPreferences() UserPreferences {
	return &userPreferences{}
}

func (m *userPreferences) SetPreference(ctx context.Context, key string, value interface{}) error {
	return nil
}
