import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
// import { TranslationModel } from '../models/translation';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}/plugin/translation`;

// Remove translation
export type TranslationRemoveParameter = {
  text: string;
  pos: number;
};
export type TranslationRemoveArg = {
  param: TranslationRemoveParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type TranslationRemoveResult = {
  param: TranslationRemoveParameter;
};
export const removeTranslation = createAsyncThunk<
  TranslationRemoveResult,
  TranslationRemoveArg,
  BaseThunkApiConfig
>('translation/remove', async (arg: TranslationRemoveArg, thunkAPI) => {
  const url = `${baseUrl}/text/${arg.param.text}/pos/${arg.param.pos}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .delete(url, jsonRequestConfig(accessToken))
        .then((resp) => {
          arg.postSuccessProcess();
          return {
            param: arg.param,
          } as TranslationRemoveResult;
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

export const translationRemoveSlice = createSlice({
  name: 'translation_remove',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(removeTranslation.pending, (state) => {
        state.loading = true;
      })
      .addCase(removeTranslation.fulfilled, (state, action) => {
        // onsole.log('workbook', action.payload.response);
        state.loading = false;
        state.failed = false;
      })
      .addCase(removeTranslation.rejected, (state, action) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectTranslationRemoveLoading = (state: RootState): boolean =>
  state.translationRemove.loading;

export const selectTranslationRemoveFailed = (state: RootState): boolean =>
  state.translationRemove.failed;

export default translationRemoveSlice.reducer;
