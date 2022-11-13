import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { extractErrorMessage } from '@/features/base';
import { jsonHeaders } from '@/utils/util';

import { TranslationModel } from '../models/translation';

const baseUrl = process.env.REACT_APP_BACKEND + '/plugin/translation';

// Get translation
export type TranslationGetParameter = {
  text: string;
  pos: number;
};
export type TranslationGetArg = {
  param: TranslationGetParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type TranslationGetResult = {
  param: TranslationGetParameter;
  response: TranslationModel;
};
export const getTranslation = createAsyncThunk<
  TranslationGetResult,
  TranslationGetArg,
  BaseThunkApiConfig
>('translation/get', async (arg: TranslationGetArg, thunkAPI) => {
  const url = `${baseUrl}/text/${arg.param.text}/pos/${arg.param.pos}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .get(url, { headers: jsonHeaders(accessToken), data: {} })
        .then((resp) => {
          const response = resp.data as TranslationModel;
          arg.postSuccessProcess();
          return {
            param: arg.param,
            response: response,
          } as TranslationGetResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface TranslationGetState {
  loading: boolean;
  failed: boolean;
  translation: TranslationModel;
}
const defaultTranslation = {
  id: 0,
  version: 0,
  updatedAt: '',
  text: '',
  pos: 0,
  translated: '',
  lang2: '',
  provider: '',
};
const initialState: TranslationGetState = {
  loading: false,
  failed: false,
  translation: defaultTranslation,
};

export const translationGetSlice = createSlice({
  name: 'translation_get',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(getTranslation.pending, (state) => {
        state.loading = true;
      })
      .addCase(getTranslation.fulfilled, (state, action) => {
        // onsole.log('workbook', action.payload.response);
        state.loading = false;
        state.failed = false;
        state.translation = action.payload.response;
      })
      .addCase(getTranslation.rejected, (state, action) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectTranslationGetLoading = (state: RootState): boolean =>
  state.translationGet.loading;

export const selectTranslationGetFailed = (state: RootState): boolean =>
  state.translationGet.failed;

export const selectTranslation = (state: RootState) =>
  state.translationGet.translation;

export default translationGetSlice.reducer;
