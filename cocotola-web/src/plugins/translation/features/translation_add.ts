import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { extractErrorMessage } from '@/features/base';
// import { TranslationModel } from '../models/translation';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = process.env.REACT_APP_BACKEND + '/plugin/translation';

// Add translation
export type TranslationAddParameter = {
  text: string;
  pos: number;
  translated: string;
  lang2: string;
};
export type TranslationAddArg = {
  param: TranslationAddParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type TranslationAddResult = {
  param: TranslationAddParameter;
};
export const addTranslation = createAsyncThunk<
  TranslationAddResult,
  TranslationAddArg,
  BaseThunkApiConfig
>('translation/add', async (arg: TranslationAddArg, thunkAPI) => {
  const url = `${baseUrl}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          arg.postSuccessProcess();
          return {
            param: arg.param,
          } as TranslationAddResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface TranslationAddState {
  loading: boolean;
  failed: boolean;
}
const initialState: TranslationAddState = {
  loading: false,
  failed: false,
};

export const translationAddSlice = createSlice({
  name: 'translation_add',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(addTranslation.pending, (state) => {
        state.loading = true;
      })
      .addCase(addTranslation.fulfilled, (state, action) => {
        // onsole.log('workbook', action.payload.response);
        state.loading = false;
        state.failed = false;
      })
      .addCase(addTranslation.rejected, (state, action) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectTranslationAddLoading = (state: RootState): boolean =>
  state.translationAdd.loading;

export const selectTranslationAddFailed = (state: RootState): boolean =>
  state.translationAdd.failed;

export default translationAddSlice.reducer;
