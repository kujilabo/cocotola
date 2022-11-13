import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
// import { TranslationModel } from '../models/translation';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}/plugin/translation`;

// Update translation
export type TranslationUpdateParameter = {
  text: string;
  pos: number;
  translated: string;
  lang2: string;
};
export type TranslationUpdateArg = {
  param: TranslationUpdateParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type TranslationUpdateResult = {
  param: TranslationUpdateParameter;
};
export const updateTranslation = createAsyncThunk<
  TranslationUpdateResult,
  TranslationUpdateArg,
  BaseThunkApiConfig
>('translation/update', async (arg: TranslationUpdateArg, thunkAPI) => {
  const url = `${baseUrl}/text/${arg.param.text}/pos/${arg.param.pos}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .put(url, arg.param, jsonRequestConfig(accessToken))
        .then(() => {
          arg.postSuccessProcess();
          return {
            param: arg.param,
          } as TranslationUpdateResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface TranslationUpdateState {
  loading: boolean;
  failed: boolean;
}
const initialState: TranslationUpdateState = {
  loading: false,
  failed: false,
};

export const translationUpdateSlice = createSlice({
  name: 'translation_update',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(updateTranslation.pending, (state) => {
        state.loading = true;
      })
      .addCase(updateTranslation.fulfilled, (state) => {
        // onsole.log('workbook', action.payload.response);
        state.loading = false;
        state.failed = false;
      })
      .addCase(updateTranslation.rejected, (state) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectTranslationUpdateLoading = (state: RootState): boolean =>
  state.translationUpdate.loading;

export const selectTranslationUpdateFailed = (state: RootState): boolean =>
  state.translationUpdate.failed;

export default translationUpdateSlice.reducer;
