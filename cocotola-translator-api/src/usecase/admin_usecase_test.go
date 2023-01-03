//go:generate mockery --output mock --name AdminUsecase
package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/domain"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/service"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/usecase"
)

func matchErrorFunc(expected error) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, actual error, args ...interface{}) bool {
		if errors.Is(actual, expected) {
			return true
		}
		return assert.Fail(t, fmt.Sprintf("error type is mismatch. expected: %v, actual: %v", expected, actual))
	}
}

func Test_adminUsecase_RemoveTranslation(t *testing.T) {
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		// logrus.SetLevel(logrus.DebugLevel)
		customRepo := ts.rf.NewCustomTranslationRepository(ctx)
		adminUsecase := usecase.NewAdminUsecase(ctx, ts.transaction)

		// given
		// - customRepo has one data
		param, err := domain.NewTranslationAddParameter("apple", domain.PosNoun, domain.Lang2JA, "りんご")
		assert.NoError(t, err)
		err = customRepo.Add(ctx, param)
		assert.NoError(t, err)
		type args struct {
			lang2 domain.Lang2
			text  string
			pos   domain.WordPos
		}
		tests := []struct {
			name      string
			args      args
			assertion assert.ErrorAssertionFunc
		}{
			{
				"word is registered",
				args{domain.Lang2JA, "apple", domain.PosNoun},
				assert.NoError,
			},
			{
				"word is not registered",
				args{domain.Lang2JA, "orange", domain.PosNoun},
				matchErrorFunc(service.ErrTranslationNotFound)},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// when
				err := adminUsecase.RemoveTranslation(ctx, tt.args.lang2, tt.args.text, tt.args.pos)

				// then
				tt.assertion(t, err)
			})
		}
	}
	test(t, fn)
}

// func Test_adminUsecase_RemoveTranslation(t *testing.T) {
// 	ctx := context.Background()

// 	// given
// 	customRepo := new(service_mock.CustomTranslationRepository)
// 	customRepo.On("Remove", anythingOfContext, domain.Lang2JA, "apple", domain.PosNoun).Return(nil)
// 	customRepo.On("Remove", anythingOfContext, domain.Lang2JA, "orange", domain.PosNoun).Return(service.ErrTranslationNotFound)
// 	rf := new(service_mock.RepositoryFactory)
// 	rf.On("NewCustomTranslationRepository", anythingOfContext).Return(customRepo, nil)
// 	transaction := newTransaction(rf)
// 	adminUsecase := usecase.NewAdminUsecase(ctx, transaction)
// 	type args struct {
// 		lang2 domain.Lang2
// 		text  string
// 		pos   domain.WordPos
// 	}
// 	tests := []struct {
// 		name      string
// 		args      args
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		{"word is registered", args{domain.Lang2JA, "apple", domain.PosNoun}, assert.NoError},
// 		{"word is not registered", args{domain.Lang2JA, "orange", domain.PosNoun}, matchErrorFunc(service.ErrTranslationNotFound)},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// when
// 			err := adminUsecase.RemoveTranslation(ctx, tt.args.lang2, tt.args.text, tt.args.pos)

// 			// then
// 			tt.assertion(t, err)
// 		})
// 	}
// }
