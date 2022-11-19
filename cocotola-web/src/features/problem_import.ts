import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}/v1/workbook`;

// Import problem
export type ProblemImportArg = {
  workbookId: number;
  param: FormData;
  postSuccessProcess: (id: number) => void;
  postFailureProcess: (error: string) => void;
};
type ProblemImportResponse = {
  id: number;
};
type ProblemImportResult = {
  response: ProblemImportResponse;
};

export const importProblem = createAsyncThunk<
  ProblemImportResult,
  ProblemImportArg,
  BaseThunkApiConfig
>('problem/import', async (arg: ProblemImportArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.workbookId}/problem/import`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as ProblemImportResponse;
          arg.postSuccessProcess(response.id);
          return { response: response } as ProblemImportResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface ProblemImportState {
  value: number;
  loading: boolean;
}

const initialState: ProblemImportState = {
  value: 0,
  loading: false,
};

export const problemImportSlice = createSlice({
  name: 'problem_import',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(importProblem.pending, (state) => {
        state.loading = true;
      })
      .addCase(importProblem.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(importProblem.rejected, (state) => {
        state.loading = false;
      });
  },
});

export const selectProblemImportLoading = (state: RootState) =>
  state.problemImport.loading;

export default problemImportSlice.reducer;
