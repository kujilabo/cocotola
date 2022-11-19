import { ReactElement, FC, useEffect } from 'react';

import jwt_decode, { JwtPayload } from 'jwt-decode';
import { Navigate } from 'react-router-dom';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import {
  selectAuthFailed,
  selectAccessToken,
  selectRefreshToken,
  selectAuthLoading,
  refreshAccessToken,
} from '@/features/auth';

type PrivateRouteProps = {
  element: JSX.Element;
};

export const PrivateRoute: FC<PrivateRouteProps> = (
  props: PrivateRouteProps
): ReactElement => {
  const dispatch = useAppDispatch();
  const loading = useAppSelector(selectAuthLoading);
  const failed = useAppSelector(selectAuthFailed);
  const accessToken = useAppSelector(selectAccessToken);
  const refreshToken = useAppSelector(selectRefreshToken);

  let isAccessTokenExpired = true;
  if (accessToken && accessToken != null && accessToken !== '') {
    const decoded = jwt_decode<JwtPayload>(accessToken) || null;
    if (decoded.exp) {
      isAccessTokenExpired = decoded.exp < new Date().getTime() / 1000;
    }
  }

  let isRefreshTokenExpired = true;
  if (refreshToken && refreshToken != null && refreshToken !== '') {
    const decoded = jwt_decode<JwtPayload>(refreshToken) || null;
    if (decoded.exp) {
      isRefreshTokenExpired = decoded.exp < new Date().getTime() / 1000;
    }
  }

  useEffect(() => {
    if (!failed && !loading && isAccessTokenExpired && !isRefreshTokenExpired) {
      // onsole.log('xxx refreshAccessToken');
      const f = async () => {
        await dispatch(
          refreshAccessToken({
            refreshToken: refreshToken,
          })
        );
      };
      f().catch(console.error);
    }
  }, [loading, refreshToken, isAccessTokenExpired, isRefreshTokenExpired]); // eslint-disable-line react-hooks/exhaustive-deps

  if (failed) {
    return <div>Authentication Failure</div>;
  } else if (!loading && isAccessTokenExpired && !isRefreshTokenExpired) {
    return <div>Refreshing...</div>;
  } else if (isRefreshTokenExpired) {
    return <Navigate replace to={`/app/login`} />;
  } else {
    console.log('children');
    return <>{props.element}</>;
  }
};
