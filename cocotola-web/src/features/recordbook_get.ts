import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { RecordbookModel } from '@/models/recordbook';
import { jsonHeaders } from '@/utils/util';

const baseUrl = `${backendUrl}/v1/study/workbook`;

// Get recordbook
export type RecordbookGetParameter = {
  workbookId: number;
  studyType: string;
};
export type RecordbookGetArg = {
  param: RecordbookGetParameter;
  postSuccessProcess: (recordbook: RecordbookModel) => void;
  postFailureProcess: (error: string) => void;
};
type RecordbookGetResult = {
  param: RecordbookGetParameter;
  response: RecordbookModel;
};
export const getRecordbook = createAsyncThunk<
  RecordbookGetResult,
  RecordbookGetArg,
  BaseThunkApiConfig
>('recordbook/get', async (arg: RecordbookGetArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.workbookId}/study_type/${arg.param.studyType}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .get(url, { headers: jsonHeaders(accessToken), data: {} })
        .then((resp) => {
          const response = resp.data as RecordbookModel;
          arg.postSuccessProcess(response);
          return { response: response } as RecordbookGetResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

// Get completion rate
export type CompletionRateGetParameter = {
  workbookId: number;
};
export type CompletionRateGetArg = {
  param: CompletionRateGetParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type CompletionRateGetResult = {
  param: CompletionRateGetParameter;
  response: { [key: string]: number };
};
export const getCompletionRate = createAsyncThunk<
  CompletionRateGetResult,
  CompletionRateGetArg,
  BaseThunkApiConfig
>(
  'recordbook/complietion_rate_get',
  async (arg: CompletionRateGetArg, thunkAPI) => {
    const url = `${baseUrl}/${arg.param.workbookId}/completion_rate`;
    const { refreshToken } = thunkAPI.getState().auth;
    return await thunkAPI
      .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
      .then(() => {
        const { accessToken } = thunkAPI.getState().auth;
        return axios
          .get(url, { headers: jsonHeaders(accessToken), data: {} })
          .then((resp) => {
            const response = resp.data as { [key: string]: number };
            arg.postSuccessProcess();
            return { response: response } as CompletionRateGetResult;
          })
          .catch((err: Error) => {
            const errorMessage = extractErrorMessage(err);
            arg.postFailureProcess(errorMessage);
            return thunkAPI.rejectWithValue(errorMessage);
          });
      });
  }
);

export interface recordbookGetState {
  loading: boolean;
  failed: boolean;
  recordbook: RecordbookModel;
  completionRateMap: { [key: string]: number };
}

const defaultRecordbook: RecordbookModel = {
  id: 0,
  records: [],
};
const initialState: recordbookGetState = {
  loading: false,
  failed: false,
  // recordbook: [new RecordbookModel(0, [])],
  recordbook: defaultRecordbook,
  completionRateMap: {},
};

export const recordbookGetSlice = createSlice({
  name: 'recordbook_get',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      // get recordbook
      .addCase(getRecordbook.pending, (state) => {
        state.loading = true;
      })
      .addCase(getRecordbook.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
        state.recordbook = action.payload.response;
      })
      .addCase(getRecordbook.rejected, (state) => {
        state.loading = false;
        state.failed = true;
      })
      // get completion rate
      .addCase(getCompletionRate.pending, (state) => {
        state.loading = true;
      })
      .addCase(getCompletionRate.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
        state.completionRateMap = action.payload.response;
      })
      .addCase(getCompletionRate.rejected, (state) => {
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectRecordbookGetLoading = (state: RootState) =>
  state.recordbookGet.loading;

export const selectRecordbookGetFailed = (state: RootState) =>
  state.recordbookGet.failed;

export const selectRecordbook = (state: RootState) =>
  state.recordbookGet.recordbook;

export const selectRecordbookCompletionRateMap = (state: RootState) =>
  state.recordbookGet.completionRateMap;

export default recordbookGetSlice.reducer;
