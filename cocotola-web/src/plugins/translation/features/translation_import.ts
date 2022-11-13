import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}`;

// Import translation
export type TranslationImportArg = {
  param: FormData;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type TranslationImportResult = Record<string, never>;

export const importTranslation = createAsyncThunk<
  TranslationImportResult,
  TranslationImportArg,
  BaseThunkApiConfig
>('translation/import', async (arg: TranslationImportArg, thunkAPI) => {
  const url = `${baseUrl}/plugin/translation/import`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          arg.postSuccessProcess();
          return {} as TranslationImportResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface TranslationImportState {
  loading: boolean;
  failed: boolean;
}

const initialState: TranslationImportState = {
  loading: false,
  failed: false,
};

export const translationImportSlice = createSlice({
  name: 'translation_import',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(importTranslation.pending, (state) => {
        state.loading = true;
      })
      .addCase(importTranslation.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
      })
      .addCase(importTranslation.rejected, (state, action) => {
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectTranslationImportLoading = (state: RootState) =>
  state.translationImport.loading;

export default translationImportSlice.reducer;
