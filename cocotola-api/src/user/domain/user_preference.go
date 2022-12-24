package domain

type UserPreference interface {
	GetPreferences() map[string]interface{}
}

type userPreference struct {
	Preferences map[string]interface{}
}

func NewUserPreference(preferences map[string]interface{}) UserPreference {
	return &userPreference{
		Preferences: preferences,
	}
}

func (m *userPreference) GetPreferences() map[string]interface{} {
	return m.Preferences
}
