import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { extractErrorMessage } from '@/features/base';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = import.meta.env.VITE_APP_BACKEND + '/v1/private/workbook';

// Remove workbook
export type WorkbookRemoveParameter = {
  id: number;
  version: number;
};
export type WorkbookRemoveArg = {
  param: WorkbookRemoveParameter;
  postSuccessProcess: (id: number) => void;
  postFailureProcess: (error: string) => void;
};
type WorkbookRemoveResponse = {
  id: number;
};
type WorkbookRemoveResult = {
  param: WorkbookRemoveParameter;
  response: WorkbookRemoveResponse;
};
export const removeWorkbook = createAsyncThunk<
  WorkbookRemoveResult,
  WorkbookRemoveArg,
  BaseThunkApiConfig
>('private/workbook/remove', async (arg: WorkbookRemoveArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.id}?version=${arg.param.version}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .delete(url, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as WorkbookRemoveResponse;
          arg.postSuccessProcess(response.id);
          return {
            param: arg.param,
            response: response,
          } as WorkbookRemoveResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface WorkbookRemoveState {
  loading: boolean;
}

const initialState: WorkbookRemoveState = {
  loading: false,
};

export const workbookRemoveSlice = createSlice({
  name: 'workbook_remove',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(removeWorkbook.pending, (state) => {
        state.loading = true;
      })
      .addCase(removeWorkbook.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(removeWorkbook.rejected, (state, action) => {
        state.loading = false;
      });
  },
});

export const selectWorkbookRemoveLoading = (state: RootState): boolean =>
  state.workbookRemove.loading;

export default workbookRemoveSlice.reducer;
