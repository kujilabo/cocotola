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

type FindSystemOwnerByOrganizationNameFunc func(context.Context, userS.RepositoryFactory, string) (userS.SystemOwner, error)

type GoogleUserUsecase interface {
	RetrieveAccessToken(ctx context.Context, code string) (*service.GoogleAuthResponse, error)

	RetrieveUserInfo(ctx context.Context, GoogleAuthResponse *service.GoogleAuthResponse) (*service.GoogleUserInfo, error)

	RegisterAppUser(ctx context.Context, googleUserInfo *service.GoogleUserInfo, googleAuthResponse *service.GoogleAuthResponse, organizationName string) (*service.TokenSet, error)
}

type googleUserUsecase struct {
	transaction                           service.Transaction
	googleAuthClient                      service.GoogleAuthClient
	authTokenManager                      service.AuthTokenManager
	registerAppUserCallback               func(ctx context.Context, organizationName string, appUser userD.AppUserModel) error
	findSystemOwnerByOrganizationNameFunc FindSystemOwnerByOrganizationNameFunc
}

func NewGoogleUserUsecase(transaction service.Transaction, googleAuthClient service.GoogleAuthClient, authTokenManager service.AuthTokenManager, registerAppUserCallback func(ctx context.Context, organizationName string, appUser userD.AppUserModel) error, findSystemOwnerByOrganizationNameFunc FindSystemOwnerByOrganizationNameFunc) GoogleUserUsecase {
	return &googleUserUsecase{
		transaction:                           transaction,
		googleAuthClient:                      googleAuthClient,
		authTokenManager:                      authTokenManager,
		registerAppUserCallback:               registerAppUserCallback,
		findSystemOwnerByOrganizationNameFunc: findSystemOwnerByOrganizationNameFunc,
	}
}

func (s *googleUserUsecase) RetrieveAccessToken(ctx context.Context, code string) (*service.GoogleAuthResponse, error) {
	resp, err := s.googleAuthClient.RetrieveAccessToken(ctx, code)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return resp, nil
}

func (s *googleUserUsecase) RetrieveUserInfo(ctx context.Context, googleAuthResponse *service.GoogleAuthResponse) (*service.GoogleUserInfo, error) {
	info, err := s.googleAuthClient.RetrieveUserInfo(ctx, googleAuthResponse)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return info, nil
}

func (s *googleUserUsecase) RegisterAppUser(ctx context.Context, googleUserInfo *service.GoogleUserInfo, googleAuthResponse *service.GoogleAuthResponse, organizationName string) (*service.TokenSet, error) {
	var tokenSet *service.TokenSet

	var organization userS.Organization
	var appUser userS.AppUser
	if err := s.transaction.Do(ctx, func(rf userS.RepositoryFactory) error {
		tmpOrganization, tmpAppUser, err := s.registerAppUser(ctx, rf, organizationName, googleUserInfo.Email, googleUserInfo.Name, googleUserInfo.Email, googleAuthResponse.AccessToken, googleAuthResponse.RefreshToken)
		if err != nil && !errors.Is(err, userS.ErrAppUserAlreadyExists) {
			return liberrors.Errorf("s.registerAppUser. err: %w", err)
		}

		organization = tmpOrganization
		appUser = tmpAppUser
		return nil
	}); err != nil {
		return nil, liberrors.Errorf("RegisterAppUser. err: %w", err)
	}

	if err := s.registerAppUserCallback(ctx, organizationName, appUser); err != nil {
		return nil, liberrors.Errorf("registerStudentCallback. err: %w", err)
	}

	tokenSetTmp, err := s.authTokenManager.CreateTokenSet(ctx, appUser, organization)
	if err != nil {
		return nil, liberrors.Errorf("s.authTokenManager.CreateTokenSet. err: %w", err)
	}
	tokenSet = tokenSetTmp
	return tokenSet, nil
}

func (s *googleUserUsecase) registerAppUser(ctx context.Context, rf userS.RepositoryFactory, organizationName string, loginID string, username string,
	providerID, providerAccessToken, providerRefreshToken string) (userS.Organization, userS.AppUser, error) {
	logger := log.FromContext(ctx)

	var organization userS.Organization
	var appUser userS.AppUser

	if err := func() error {
		systemOwner, err := s.findSystemOwnerByOrganizationNameFunc(ctx, rf, organizationName)
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
			return liberrors.Errorf("systemOwner.FindAppUserByLoginID. err: %w", err)
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
			return nil, nil, liberrors.Errorf("registerAppUser. err: %w", err)
		}
	}
	return organization, appUser, nil
}
