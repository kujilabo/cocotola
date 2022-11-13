import { ReactElement } from 'react';

import jwt_decode, { JwtPayload } from 'jwt-decode';
import { parse } from 'query-string';
import { Navigate } from 'react-router';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppDimmer } from '@/components/AppDimmer';
import {
  selectAccessToken,
  selectAuthLoading,
  selectAuthFailed,
  googleAuthorize,
} from '@/features/auth';
import { emptyFunction } from '@/utils/util';

export const LoginCallback = (): ReactElement => {
  const dispatch = useAppDispatch();
  const accessToken = useAppSelector(selectAccessToken);
  let isAccessTokenExpired = true;
  if (accessToken && accessToken != null && accessToken !== '') {
    // onsole.log('decode acc', accessToken);
    const decoded = jwt_decode<JwtPayload>(accessToken) || null;
    if (decoded.exp) {
      isAccessTokenExpired = decoded.exp < new Date().getTime() / 1000;
    }
  }

  const authLoading = useAppSelector(selectAuthLoading);
  const authFailed = useAppSelector(selectAuthFailed);
  const location = window.location.search;
  // onsole.log('Callback', authLoading, isAccessTokenExpired);
  if (authFailed) {
    return <div>Failed</div>;
  } else if (authLoading === false && isAccessTokenExpired) {
    const parsed = parse(location);
    // const code: string '' || parsed.code || '';
    const code = parsed ? String(parsed.code) : '';

    const f = async () => {
      await dispatch(
        googleAuthorize({
          param: {
            organizationName: 'cocotola',
            code: code,
          },
          postSuccessProcess: emptyFunction,
          postFailureProcess: (error: string) => {
            console.log('callback error', error);
            return;
          },
        })
      );
    };
    f().catch(console.error);
    return <AppDimmer />;
  } else if (!isAccessTokenExpired) {
    return <Navigate replace to="/" />;
  } else {
    return <AppDimmer />;
  }
};
