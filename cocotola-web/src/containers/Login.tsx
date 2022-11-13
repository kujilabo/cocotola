import { ReactElement } from 'react';

import { Button, Container, Menu } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { selectRedirectUrl, redirectTo } from '@/features/auth';
import { clientId, frontendUrl } from '@/features/base';

import { emptyFunction } from '../../../bazel-cocotola/cocotola-web/src/utils/util';

export const Login = (): ReactElement => {
  const dispatch = useAppDispatch();
  const redirectUrl = useAppSelector(selectRedirectUrl);
  const googleAuth = () => {
    let url = 'https://accounts.google.com/o/oauth2/auth';
    url += '?client_id=';
    url += clientId;
    url += '&redirect_uri=';
    url += frontendUrl;
    url += '/app/callback';
    url += '&scope=profile email';
    url += '&response_type=';
    url += 'code';
    url += '&access_type=';
    url += 'offline';
    console.log(url);
    dispatch(redirectTo({ url: url }));
  };
  const guestAuth = () => emptyFunction;

  if (redirectUrl && redirectUrl !== '') {
    console.log('redirect');
    window.location.href = redirectUrl;
    // return <Navigate replace to={redirectUrl} />;
  }

  return (
    <div>
      <Menu inverted></Menu>
      <Container fluid>
        <Button basic color="teal" onClick={googleAuth}>
          Sign in with Google
        </Button>
        <Button basic color="teal" onClick={guestAuth}>
          Sign in as Guest
        </Button>
      </Container>
    </div>
  );
};
