import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { WorkbookModel } from '@/models/workbook';
import { jsonHeaders } from '@/utils/util';

const baseUrl = `${backendUrl}/v1/private/workbook`;

// Get workbook
export type WorkbookGetParameter = {
  id: number;
};
export type WorkbookGetArg = {
  param: WorkbookGetParameter;
  postSuccessProcess: (workbook: WorkbookModel) => void;
  postFailureProcess: (error: string) => void;
};
type WorkbookGetResult = {
  param: WorkbookGetParameter;
  response: WorkbookModel;
};
export const getWorkbook = createAsyncThunk<
  WorkbookGetResult,
  WorkbookGetArg,
  BaseThunkApiConfig
>('workbook/view', async (arg: WorkbookGetArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.id}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .get(url, { headers: jsonHeaders(accessToken), data: {} })
        .then((resp) => {
          const response = resp.data as WorkbookModel;
          arg.postSuccessProcess(response);
          return { param: arg.param, response: response } as WorkbookGetResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface WorkbookGetState {
  loading: boolean;
  failed: boolean;
  workbook: WorkbookModel;
}

const defaultWorkbook = {
  id: 0,
  version: 0,
  problemType: '',
  name: '',
  questionText: '',
  subscribed: false,
};
const initialState: WorkbookGetState = {
  loading: false,
  failed: false,
  workbook: defaultWorkbook,
};

export const workbookGetSlice = createSlice({
  name: 'workbook_get',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(getWorkbook.pending, (state) => {
        state.loading = true;
      })
      .addCase(getWorkbook.fulfilled, (state, action) => {
        // onsole.log('workbook', action.payload.response);
        state.loading = false;
        state.failed = false;
        state.workbook = action.payload.response;
      })
      .addCase(getWorkbook.rejected, (state) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectWorkbookGetLoading = (state: RootState): boolean =>
  state.workbookGet.loading;

export const selectWorkbookGetFailed = (state: RootState): boolean =>
  state.workbookGet.failed;

export const selectWorkbook = (state: RootState) => state.workbookGet.workbook;

export default workbookGetSlice.reducer;
