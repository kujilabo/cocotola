package usecase

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/auth/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type GuestUserUsecase interface {
	RetrieveGuestToken(ctx context.Context, organizationName string) (*service.TokenSet, error)
}

type guestUserUsecase struct {
	transaction      service.Transaction
	authTokenManager service.AuthTokenManager
}

func NewGuestUserUsecase(transaction service.Transaction, authTokenManager service.AuthTokenManager) GuestUserUsecase {
	return &guestUserUsecase{
		transaction:      transaction,
		authTokenManager: authTokenManager,
	}
}

func (s *guestUserUsecase) RetrieveGuestToken(ctx context.Context, organizationName string) (*service.TokenSet, error) {
	var tokenSet *service.TokenSet

	if err := s.transaction.Do(ctx, func(rf userS.RepositoryFactory) error {
		systemAdmin, err := userS.NewSystemAdmin(ctx, rf)
		if err != nil {
			return liberrors.Errorf(" userS.NewSystemAdmin. err: %w", err)
		}

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, organizationName)
		if err != nil {
			return liberrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		// guest, err := systemOwner.FindAppUserByLoginID(ctx, "guest")
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to FindAppUserByLoginID. err: %w", err)
		// }

		organization, err := systemOwner.GetOrganization(ctx)
		if err != nil {
			return liberrors.Errorf("failed to GetOrganization. err: %w", err)
		}

		model, err := userD.NewModel(0, 1, time.Now(), time.Now(), 0, 0)
		if err != nil {
			return liberrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
		}

		guest, err := userD.NewAppUserModel(model, userD.OrganizationID(organization.GetID()), "guest", "Guest", []string{}, map[string]string{})
		if err != nil {
			return liberrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
		}

		tokenSetTmp, err := s.authTokenManager.CreateTokenSet(ctx, guest, organization)
		if err != nil {
			return liberrors.Errorf("s.authTokenManager.CreateTokenSet. err: %w", err)
		}

		tokenSet = tokenSetTmp
		return nil
	}); err != nil {
		return nil, liberrors.Errorf("RetrieveGuestToken. err: %w", err)
	}
	return tokenSet, nil
}
