import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { WorkbookModel } from '@/models/workbook';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}/v1/private/workbook`;

// Find my workbooks
export type MyWorkbookFindParameter = {
  pageNo: number;
  pageSize: number;
  spaceKey: string;
};
export type MyWorkbookFindArg = {
  param: MyWorkbookFindParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type MyWorkbookFindResponse = {
  results: WorkbookModel[];
  totalCount: number;
};
type MyWorkbookFindResult = {
  param: MyWorkbookFindParameter;
  response: MyWorkbookFindResponse;
};

export const findMyWorkbooks = createAsyncThunk<
  MyWorkbookFindResult,
  MyWorkbookFindArg,
  BaseThunkApiConfig
>('private/workbook/find', async (arg: MyWorkbookFindArg, thunkAPI) => {
  const url = `${baseUrl}/search`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      // onsole.log('accessToken1');
      const { accessToken } = thunkAPI.getState().auth;
      // onsole.log('accessToken', accessToken);
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          // onsole.log('the1', resp);
          const response = resp.data as MyWorkbookFindResponse;
          // onsole.log('the2', response);
          arg.postSuccessProcess();
          return {
            param: arg.param,
            response: response,
          } as MyWorkbookFindResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface WorkbookFindState {
  value: number;
  loading: boolean;
  failed: boolean;
  workbooksMap: { [key: string]: WorkbookModel[] };
  workbooksLoadedMap: { [key: string]: boolean };
  workbooks: WorkbookModel[];
  workbooksTotalCount: number;
}

const initialState: WorkbookFindState = {
  value: 0,
  loading: false,
  failed: false,
  workbooksLoadedMap: {},
  workbooksMap: {},
  workbooks: [],
  workbooksTotalCount: 0,
};

export const workbookFindSlice = createSlice({
  name: 'workbook_find',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(findMyWorkbooks.pending, (state) => {
        state.loading = true;
      })
      .addCase(findMyWorkbooks.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
        state.workbooks = action.payload.response.results;
        state.workbooksTotalCount = action.payload.response.totalCount;
        state.workbooksMap = {
          ...state.workbooksMap,
          [action.payload.param.spaceKey]: action.payload.response.results,
        };
        state.workbooksLoadedMap = {
          ...state.workbooksLoadedMap,
          [action.payload.param.spaceKey]: true,
        };
        // onsole.log('workbooksLoadedMap', state.workbooksLoadedMap);
      })
      .addCase(findMyWorkbooks.rejected, (state) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectWorkbookFindLoading = (state: RootState): boolean =>
  state.workbookFind.loading;

export const selectWorkbookFindFailed = (state: RootState): boolean =>
  state.workbookFind.failed;

export const selectWorkbooksLoadedMap = (
  state: RootState
): { [key: string]: boolean } => state.workbookFind.workbooksLoadedMap;

export const selectWorkbooksMap = (
  state: RootState
): { [key: string]: WorkbookModel[] } => state.workbookFind.workbooksMap;

export default workbookFindSlice.reducer;
