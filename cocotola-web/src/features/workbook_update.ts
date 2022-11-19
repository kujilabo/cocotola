import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}/v1/private/workbook`;

// Update workbook
export type WorkbookUpdateParameter = {
  id: number;
  version: number;
  name: string;
  questionText: string;
};
export type WorkbookUpdateArg = {
  param: WorkbookUpdateParameter;
  postSuccessProcess: (id: number) => void;
  postFailureProcess: (error: string) => void;
};
type WorkbookUpdateResponse = {
  id: number;
};
type WorkbookUpdateResult = {
  param: WorkbookUpdateParameter;
  response: WorkbookUpdateResponse;
};

export const updateWorkbook = createAsyncThunk<
  WorkbookUpdateResult,
  WorkbookUpdateArg,
  BaseThunkApiConfig
>('private/workbook/update', async (arg: WorkbookUpdateArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.id}?version=${arg.param.version}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .put(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as WorkbookUpdateResponse;
          arg.postSuccessProcess(response.id);
          return {
            param: arg.param,
            response: response,
          } as WorkbookUpdateResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface WorkbookUpdateState {
  loading: boolean;
}

const initialState: WorkbookUpdateState = {
  loading: false,
};

export const workbookUpdateSlice = createSlice({
  name: 'workbook_update',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(updateWorkbook.pending, (state) => {
        state.loading = true;
      })
      .addCase(updateWorkbook.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(updateWorkbook.rejected, (state) => {
        state.loading = false;
      });
  },
});

export const selectWorkbookUpdateLoading = (state: RootState): boolean =>
  state.workbookUpdate.loading;

export default workbookUpdateSlice.reducer;
