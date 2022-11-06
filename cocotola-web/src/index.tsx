import { StrictMode } from 'react';

import { createRoot } from 'react-dom/client';
import { Provider } from 'react-redux';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { persistStore } from 'redux-persist';
import { PersistGate } from 'redux-persist/integration/react';

import '@/index.css';
import { store } from '@/app/store';
import { PrivateRoute } from '@/components/PrivateRoute';
import { App } from '@/containers/App';
import { HealthCheck } from '@/containers/HealthCheck';
import { Login } from '@/containers/Login';
import { LoginCallback } from '@/containers/LoginCallback';

const persistor = persistStore(store);

const Index = () => (
  <Provider store={store}>
    <PersistGate loading={null} persistor={persistor}>
      <BrowserRouter>
        <Routes>
          <Route path={`/healthcheck`} element={<HealthCheck />} />
          <Route path={`/app/login`} element={<Login />} />
          <Route path={`/app/callback`} element={<LoginCallback />} />
          <Route path={`*`} element={<PrivateRoute element={<App />} />} />
        </Routes>
      </BrowserRouter>
    </PersistGate>
  </Provider>
);

createRoot(document.getElementById('root') as HTMLElement).render(
  <StrictMode>
    <Index />
  </StrictMode>
);
