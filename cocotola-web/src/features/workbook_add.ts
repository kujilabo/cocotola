import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}/v1/private/workbook`;

// Add workbook
export type WorkbookAddParameter = {
  name: string;
  lang2: string;
  questionText: string;
  spaceKey: string;
};
export type WorkbookAddArg = {
  param: WorkbookAddParameter;
  postSuccessProcess: (id: number) => void;
  postFailureProcess: (error: string) => void;
};
type WorkbookAddResponse = {
  id: number;
};
type WorkbookAddResult = {
  param: WorkbookAddParameter;
  response: WorkbookAddResponse;
};

export const addWorkbook = createAsyncThunk<
  WorkbookAddResult,
  WorkbookAddArg,
  BaseThunkApiConfig
>('private/workbook/new', async (arg: WorkbookAddArg, thunkAPI) => {
  const url = baseUrl;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as WorkbookAddResponse;
          arg.postSuccessProcess(response.id);
          return { param: arg.param, response: response } as WorkbookAddResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface WorkbookNewState {
  value: number;
  loading: boolean;
}

const initialState: WorkbookNewState = {
  value: 0,
  loading: false,
};

export const workbookNewSlice = createSlice({
  name: 'workbook_new',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(addWorkbook.pending, (state) => {
        state.loading = true;
      })
      .addCase(addWorkbook.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(addWorkbook.rejected, (state, action) => {
        state.loading = false;
      });
  },
});

export const selectWorkbookAddLoading = (state: RootState): boolean =>
  state.workbookAdd.loading;

export default workbookNewSlice.reducer;
