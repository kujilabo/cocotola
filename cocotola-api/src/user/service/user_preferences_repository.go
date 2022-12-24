package service

import "context"

type UserPreferenceRepository interface {
	SetPreference(ctx context.Context)
}
