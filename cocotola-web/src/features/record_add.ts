import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}/v1/study/workbook`;

// Add record
export type RecordAddParameter = {
  workbookId: number;
  studyType: string;
  problemId: number;
  result: boolean;
  mastered: boolean;
};
export type RecordAddArg = {
  param: RecordAddParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type RecordAddResult = {
  param: RecordAddParameter;
};
export const addRecord = createAsyncThunk<
  RecordAddResult,
  RecordAddArg,
  BaseThunkApiConfig
>('record/add', async (arg: RecordAddArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.workbookId}/study_type/${arg.param.studyType}/problem/${arg.param.problemId}/record`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then(() => {
          arg.postSuccessProcess();
          return { param: arg.param } as RecordAddResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface recordAddState {
  loading: boolean;
  failed: boolean;
}

const initialState: recordAddState = {
  loading: false,
  failed: false,
};

export const recordAddSlice = createSlice({
  name: 'record_add',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(addRecord.pending, (state) => {
        state.loading = true;
      })
      .addCase(addRecord.fulfilled, (state) => {
        state.loading = false;
        state.failed = false;
      })
      .addCase(addRecord.rejected, (state) => {
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectRecordAddoading = (state: RootState) =>
  state.recordAdd.loading;

export const selectRecordAddFailed = (state: RootState) =>
  state.recordAdd.failed;

export default recordAddSlice.reducer;
