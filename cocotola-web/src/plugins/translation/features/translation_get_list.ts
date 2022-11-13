import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { extractErrorMessage } from '@/features/base';
import { jsonHeaders } from '@/utils/util';

import { TranslationModel } from '../models/translation';

const baseUrl = process.env.REACT_APP_BACKEND + '/plugin/translation';

// Get translations
export type TranslationGetListParameter = {
  text: string;
};
export type TranslationGetListArg = {
  param: TranslationGetListParameter;
  postSuccessProcess: (translations: TranslationModel[]) => void;
  postFailureProcess: (error: string) => void;
};
type TranslationGetListResult = {
  param: TranslationGetListParameter;
  response: TranslationModel[];
};
export const getTranslations = createAsyncThunk<
  TranslationGetListResult,
  TranslationGetListArg,
  BaseThunkApiConfig
>('translation/get', async (arg: TranslationGetListArg, thunkAPI) => {
  const url = `${baseUrl}/text/${arg.param.text}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .get(url, { headers: jsonHeaders(accessToken), data: {} })
        .then((resp) => {
          const response = resp.data as TranslationModel[];
          arg.postSuccessProcess(response);
          return {
            param: arg.param,
            response: response,
          } as TranslationGetListResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface TranslationGetListState {
  loading: boolean;
  failed: boolean;
  translations: TranslationModel[];
}
const initialState: TranslationGetListState = {
  loading: false,
  failed: false,
  translations: [],
};

export const translationGetListSlice = createSlice({
  name: 'translation_get_list',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(getTranslations.pending, (state) => {
        state.loading = true;
      })
      .addCase(getTranslations.fulfilled, (state, action) => {
        console.log('workbook', action.payload.response);
        state.loading = false;
        state.failed = false;
        state.translations = action.payload.response;
      })
      .addCase(getTranslations.rejected, (state, action) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectTranslationGetListLoading = (state: RootState): boolean =>
  state.translationGetList.loading;

export const selectTranslationGetListFailed = (state: RootState): boolean =>
  state.translationGetList.failed;

export const selectTranslations = (state: RootState) =>
  state.translationGetList.translations;

export default translationGetListSlice.reducer;
