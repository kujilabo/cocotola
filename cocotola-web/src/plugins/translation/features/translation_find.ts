import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { jsonRequestConfig } from '@/utils/util';

import { TranslationModel } from '../models/translation';

const baseUrl = `${backendUrl}/plugin/translation`;

// Find translations
export type TranslationFindParameter = {
  letter: string;
};
export type TranslationFindArg = {
  param: TranslationFindParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type TranslationFindResponse = {
  results: TranslationModel[];
};
type TranslationFindResult = {
  param: TranslationFindParameter;
  response: TranslationFindResponse;
};

export const findTranslations = createAsyncThunk<
  TranslationFindResult,
  TranslationFindArg,
  BaseThunkApiConfig
>('translation/find', async (arg: TranslationFindArg, thunkAPI) => {
  const url = `${baseUrl}/find`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as TranslationFindResponse;
          arg.postSuccessProcess();
          return {
            param: arg.param,
            response: response,
          } as TranslationFindResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface TranslationFindState {
  value: number;
  loading: boolean;
  failed: boolean;
  translations: TranslationModel[];
}

const initialState: TranslationFindState = {
  value: 0,
  loading: false,
  failed: false,
  translations: [],
};

export const translationFindSlice = createSlice({
  name: 'translation_find',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(findTranslations.pending, (state) => {
        state.loading = true;
      })
      .addCase(findTranslations.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
        state.translations = action.payload.response.results || [];
      })
      .addCase(findTranslations.rejected, (state) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectTranslationFindLoading = (state: RootState): boolean =>
  state.translationFind.loading;

export const selectTranslationFindFailed = (state: RootState): boolean =>
  state.translationFind.failed;

export const selectTranslations = (state: RootState): TranslationModel[] =>
  state.translationFind.translations;

export default translationFindSlice.reducer;
