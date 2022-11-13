import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { extractErrorMessage } from '@/features/base';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${process.env.REACT_APP_BACKEND}/plugin/tatoeba`;

// Import tatoeba sentence
export type TatoebaImportArg = {
  param: FormData;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type TatoebaImportResult = {};

export const importTatoebaSentence = createAsyncThunk<
  TatoebaImportResult,
  TatoebaImportArg,
  BaseThunkApiConfig
>('tatoeba/sentence/import', async (arg: TatoebaImportArg, thunkAPI) => {
  const url = `${baseUrl}/sentence/import`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          arg.postSuccessProcess();
          return {} as TatoebaImportResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export const importTatoebaLink = createAsyncThunk<
  TatoebaImportResult,
  TatoebaImportArg,
  BaseThunkApiConfig
>('tatoeba/link/import', async (arg: TatoebaImportArg, thunkAPI) => {
  const url = `${baseUrl}/link/import`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          arg.postSuccessProcess();
          return {} as TatoebaImportResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface TatoebaImportState {
  loading: boolean;
  failed: boolean;
}

const initialState: TatoebaImportState = {
  loading: false,
  failed: false,
};

export const tatoebaImportSlice = createSlice({
  name: 'tatoeba_import',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(importTatoebaSentence.pending, (state) => {
        state.loading = true;
      })
      .addCase(importTatoebaSentence.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
      })
      .addCase(importTatoebaSentence.rejected, (state, action) => {
        state.loading = false;
        state.failed = true;
      })
      .addCase(importTatoebaLink.pending, (state) => {
        state.loading = true;
      })
      .addCase(importTatoebaLink.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
      })
      .addCase(importTatoebaLink.rejected, (state, action) => {
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectTatoebaImportLoading = (state: RootState) =>
  state.tatoebaImport.loading;

export default tatoebaImportSlice.reducer;
