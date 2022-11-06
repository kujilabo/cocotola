import { ReactElement } from 'react';

import jwt_decode, { JwtPayload } from 'jwt-decode';
import { Route, Routes, Link } from 'react-router-dom';
import { Menu, Dropdown } from 'semantic-ui-react';

import '@/containers/App.css';
import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { initI18n } from '@/app/i18n';
import { PrivateWorkbookEdit } from '@/containers/private_workbook/PrivateWorkbookEdit';
import { PrivateWorkbookList } from '@/containers/private_workbook/PrivateWorkbookList';
import { PrivateWorkbookNew } from '@/containers/private_workbook/PrivateWorkbookNew';
import { PrivateWorkbookView } from '@/containers/private_workbook/PrivateWorkbookView';
import { logout, selectAccessToken } from '@/features/auth';

export interface AppJwtPayload extends JwtPayload {
  username: string;
  role: string;
}

initI18n();

export const App = (): ReactElement => {
  const dispatch = useAppDispatch();
  const accessToken = useAppSelector(selectAccessToken);
  // const redirectPath = useAppSelector(selectRedirectPath);

  // if (redirectPath === '/app/login') {
  //   // onsole.log('aaa');
  //   history.push('/app/login');
  //   return <div></div>;
  // } else if (redirectPath !== '') {
  //   // onsole.log('bbb');
  //   history.push(redirectPath);
  //   return <div></div>;
  // }

  const decoded =
    accessToken && accessToken != null && accessToken !== ''
      ? jwt_decode<AppJwtPayload>(accessToken) || null
      : null;
  const username = decoded ? decoded.username : '';
  const role = decoded ? decoded.role : '';
  // onsole.log('decoded', decoded);
  // onsole.log('role', role);

  return (
    <div>
      <Menu>
        <Menu.Item>
          <Link to={'/app/private/workbook'}>Private space</Link>
        </Menu.Item>
        <Menu.Item>
          <Link to={'/app/space/1/workbook'}>Public space</Link>
        </Menu.Item>
        {role == 'Owner' ? (
          <Dropdown item text="Plugin">
            <Dropdown.Menu>
              <Dropdown.Item>
                <Link to={'/plugin/translation/list'}>Translation</Link>
              </Dropdown.Item>
              <Dropdown.Item>
                <Link to={'/plugin/tatoeba/list'}>Tatoeba</Link>
              </Dropdown.Item>
            </Dropdown.Menu>
          </Dropdown>
        ) : (
          <div />
        )}

        <Menu.Menu position="right">
          <Dropdown item text="" icon="bars">
            <Dropdown.Menu>
              <Dropdown.Item>{username}</Dropdown.Item>
              <Dropdown.Item onClick={() => dispatch(logout())}>
                Sign out
              </Dropdown.Item>
            </Dropdown.Menu>
          </Dropdown>
        </Menu.Menu>
      </Menu>

      <Routes>
        <Route
          path={`/app/private/workbook`}
          element={<PrivateWorkbookList />}
        />
        <Route path="/app/private/workbook/new">
          element={<PrivateWorkbookNew />}
        </Route>
        <Route path="/app/private/workbook/:_workbookId">
          element={<PrivateWorkbookView />}
        </Route>
        <Route path="/app/private/workbook/:_workbookId/edit">
          element={<PrivateWorkbookEdit />}
        </Route>
      </Routes>
    </div>
  );
};
