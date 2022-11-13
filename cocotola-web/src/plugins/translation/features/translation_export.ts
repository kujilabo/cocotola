import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { extractErrorMessage } from '@/features/base';
import { blobRequestConfig } from '@/utils/util';

const baseUrl = `${process.env.REACT_APP_BACKEND}/plugin/translation`;

// Export translation
export type TranslationExportArg = {
  postSuccessProcess: (blog: Blob) => void;
  postFailureProcess: (error: string) => void;
};
type TranslationExportResult = {};

export const exportTranslation = createAsyncThunk<
  TranslationExportResult,
  TranslationExportArg,
  BaseThunkApiConfig
>('translation/export', async (arg: TranslationExportArg, thunkAPI) => {
  const url = `${baseUrl}/export`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, {}, blobRequestConfig(accessToken))
        .then((resp) => {
          arg.postSuccessProcess(new Blob([resp.data]));
          return {} as TranslationExportResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface TranslationExportState {
  loading: boolean;
  failed: boolean;
}

const initialState: TranslationExportState = {
  loading: false,
  failed: false,
};

export const translationExportSlice = createSlice({
  name: 'translation_export',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(exportTranslation.pending, (state) => {
        state.loading = true;
      })
      .addCase(exportTranslation.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
      })
      .addCase(exportTranslation.rejected, (state, action) => {
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectTranslationExportLoading = (state: RootState) =>
  state.translationExport.loading;

export default translationExportSlice.reducer;
