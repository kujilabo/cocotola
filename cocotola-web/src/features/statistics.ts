import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { StatHistoryModel } from '@/models/stat';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}/v1/stat`;

// Get stat
export type StatGetParameter = {
  pageNo: number;
};
export type StatGetArg = {
  param: StatGetParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type StatGetResponse = {
  history: StatHistoryModel;
};
type StatGetResult = {
  param: StatGetParameter;
  response: StatGetResponse;
};

export const getStat = createAsyncThunk<
  StatGetResult,
  StatGetArg,
  BaseThunkApiConfig
>('stat', async (arg: StatGetArg, thunkAPI) => {
  const url = `${baseUrl}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      // onsole.log('accessToken1');
      const { accessToken } = thunkAPI.getState().auth;
      // onsole.log('accessToken', accessToken);
      return axios
        .get(url, jsonRequestConfig(accessToken))
        .then((resp) => {
          // onsole.log('the1', resp);
          const response = resp.data as StatGetResponse;
          console.log('the2', response);
          arg.postSuccessProcess();
          return {
            param: arg.param,
            response: response,
          } as StatGetResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

const defaulHistory = {
  results: [],
};
export interface StatState {
  loading: boolean;
  failed: boolean;
  loaded: boolean;
  history: StatHistoryModel;
}

const initialState: StatState = {
  loading: false,
  failed: false,
  loaded: false,
  history: defaulHistory,
};

export const statSlice = createSlice({
  name: 'stat',
  initialState: initialState,
  reducers: {
    init: (state) => {
      state.loaded = false;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(getStat.pending, (state) => {
        state.loading = true;
        console.log('');
      })
      .addCase(getStat.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
        state.history = action.payload.response.history;
        console.log('history', state.history);
      })
      .addCase(getStat.rejected, (state) => {
        state.loading = false;
        state.failed = true;
        console.log('');
      });
  },
});

export const selectStatLoading = (state: RootState): boolean =>
  state.stat.loading;

export const selectStatLoaded = (state: RootState): boolean =>
  state.stat.loaded;

export const selectStatFailed = (state: RootState): boolean =>
  state.stat.failed;

export const selectStatHistory = (state: RootState): StatHistoryModel =>
  state.stat.history;

export default statSlice.reducer;
