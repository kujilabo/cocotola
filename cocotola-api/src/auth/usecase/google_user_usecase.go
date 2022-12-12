package usecase

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/auth/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type GoogleUserUsecase interface {
	RetrieveAccessToken(ctx context.Context, code string) (*service.GoogleAuthResponse, error)

	RetrieveUserInfo(ctx context.Context, GoogleAuthResponse *service.GoogleAuthResponse) (*service.GoogleUserInfo, error)

	RegisterAppUser(ctx context.Context, googleUserInfo *service.GoogleUserInfo, googleAuthResponse *service.GoogleAuthResponse, organizationName string) (*service.TokenSet, error)
}

type googleUserUsecase struct {
	transaction             service.Transaction
	googleAuthClient        service.GoogleAuthClient
	authTokenManager        service.AuthTokenManager
	registerAppUserCallback func(ctx context.Context, organizationName string, appUser userD.AppUserModel) error
}

func NewGoogleUserUsecase(transaction service.Transaction, googleAuthClient service.GoogleAuthClient, authTokenManager service.AuthTokenManager, registerAppUserCallback func(ctx context.Context, organizationName string, appUser userD.AppUserModel) error) GoogleUserUsecase {
	return &googleUserUsecase{
		transaction:             transaction,
		googleAuthClient:        googleAuthClient,
		authTokenManager:        authTokenManager,
		registerAppUserCallback: registerAppUserCallback,
	}
}

func (s *googleUserUsecase) RetrieveAccessToken(ctx context.Context, code string) (*service.GoogleAuthResponse, error) {
	return s.googleAuthClient.RetrieveAccessToken(ctx, code)
}

func (s *googleUserUsecase) RetrieveUserInfo(ctx context.Context, googleAuthResponse *service.GoogleAuthResponse) (*service.GoogleUserInfo, error) {
	return s.googleAuthClient.RetrieveUserInfo(ctx, googleAuthResponse)
}

func (s *googleUserUsecase) RegisterAppUser(ctx context.Context, googleUserInfo *service.GoogleUserInfo, googleAuthResponse *service.GoogleAuthResponse, organizationName string) (*service.TokenSet, error) {
	var tokenSet *service.TokenSet

	var organization userS.Organization
	var appUser userS.AppUser
	if err := s.transaction.Do(ctx, func(rf userS.RepositoryFactory) error {
		systemAdmin := userS.NewSystemAdmin(rf)

		tmpOrganization, tmpAppUser, err := s.registerAppUser(ctx, systemAdmin, organizationName, googleUserInfo.Email, googleUserInfo.Name, googleUserInfo.Email, googleAuthResponse.AccessToken, googleAuthResponse.RefreshToken)
		if err != nil && !errors.Is(err, userS.ErrAppUserAlreadyExists) {
			return err
		}

		organization = tmpOrganization
		appUser = tmpAppUser
		return nil
	}); err != nil {
		return nil, err
	}

	if err := s.registerAppUserCallback(ctx, organizationName, appUser); err != nil {
		return nil, liberrors.Errorf("failed to registerStudentCallback. err: %w", err)
	}

	tokenSetTmp, err := s.authTokenManager.CreateTokenSet(ctx, appUser, organization)
	if err != nil {
		return nil, err
	}
	tokenSet = tokenSetTmp
	return tokenSet, nil
}

func (s *googleUserUsecase) registerAppUser(ctx context.Context, systemAdmin userS.SystemAdmin, organizationName string, loginID string, username string,
	providerID, providerAccessToken, providerRefreshToken string) (userS.Organization, userS.AppUser, error) {
	logger := log.FromContext(ctx)

	var organization userS.Organization
	var appUser userS.AppUser

	if err := func() error {
		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, organizationName)
		if err != nil {
			return liberrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		tmpOrganization, err := systemOwner.GetOrganization(ctx)
		if err != nil {
			return liberrors.Errorf("failed to FindOrganization. err: %w", err)
		}

		appUser1, err := systemOwner.FindAppUserByLoginID(ctx, loginID)
		if err == nil {
			organization = tmpOrganization
			appUser = appUser1
			return userS.ErrAppUserAlreadyExists
		}

		if !errors.Is(err, userS.ErrAppUserNotFound) {
			logger.Infof("Unsupported %v", err)
			return err
		}

		logger.Infof("Add student. %+v", appUser1)
		parameter, err := userS.NewAppUserAddParameter(
			loginID,  //googleUserInfo.Email,
			username, //googleUserInfo.Name,
			[]string{""},
			map[string]string{
				"password":             "----",
				"provider":             "google",
				"providerId":           providerID,           //googleUserInfo.Email,
				"providerAccessToken":  providerAccessToken,  // googleAuthResponse.AccessToken,
				"providerRefreshToken": providerRefreshToken, //googleAuthResponse.RefreshToken,
			},
		)
		if err != nil {
			return liberrors.Errorf("invalid AppUserAddParameter. err: %w", err)
		}

		studentID, err := systemOwner.AddAppUser(ctx, parameter)
		if err != nil {
			return liberrors.Errorf("failed to AddStudent. err: %w", err)
		}

		appUser2, err := systemOwner.FindAppUserByID(ctx, studentID)
		if err != nil {
			return liberrors.Errorf("failed to FindStudentByID. err: %w", err)
		}

		appUser = appUser2
		return nil
	}(); err != nil {
		if errors.Is(err, userS.ErrAppUserAlreadyExists) {
			return organization, appUser, nil
		} else {
			return nil, nil, err
		}
	}
	return organization, appUser, nil
}
