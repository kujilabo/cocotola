import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
// import { resetLoadedProblemId } from '@/features/problem_get';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}/v1/workbook`;

// Update problem
export type ProblemUpdateParameter = {
  workbookId: number;
  problemId: number;
  version: number;
  problemType: string;
  properties: { [key: string]: string };
};
export type ProblemUpdateArg = {
  param: ProblemUpdateParameter;
  postSuccessProcess: (id: number) => void;
  postFailureProcess: (error: string) => void;
};
type ProblemUpdateResponse = {
  id: number;
};
type ProblemUpdateResult = {
  param: ProblemUpdateParameter;
  response: ProblemUpdateResponse;
};

export const updateProblem = createAsyncThunk<
  ProblemUpdateResult,
  ProblemUpdateArg,
  BaseThunkApiConfig
>('problem/update', async (arg: ProblemUpdateArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.workbookId}/problem/${arg.param.problemId}?version=${arg.param.version}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;

      return axios
        .put(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as ProblemUpdateResponse;
          arg.postSuccessProcess(response.id);
          // thunkAPI.dispatch(
          //   resetLoadedProblemId(`${arg.param.problemId}:${arg.param.version}`)
          // );
          return {
            param: arg.param,
            response: response,
          } as ProblemUpdateResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

// Update problem property
export type ProblemPropertyUpdateParameter = {
  workbookId: number;
  problemId: number;
  version: number;
  problemType: string;
  properties: { [key: string]: string };
};
export type ProblemPropertyUpdateArg = {
  param: ProblemPropertyUpdateParameter;
  postSuccessProcess: (id: number) => void;
  postFailureProcess: (error: string) => void;
};
type ProblemPropertyUpdateResponse = {
  id: number;
};
type ProblemPropertyUpdateResult = {
  param: ProblemPropertyUpdateParameter;
  response: ProblemPropertyUpdateResponse;
};

export const updateProblemProperty = createAsyncThunk<
  ProblemPropertyUpdateResult,
  ProblemPropertyUpdateArg,
  BaseThunkApiConfig
>('problem/update', async (arg: ProblemPropertyUpdateArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.workbookId}/problem/${arg.param.problemId}/property?version=${arg.param.version}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;

      return axios
        .put(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as ProblemPropertyUpdateResponse;
          arg.postSuccessProcess(response.id);
          // thunkAPI.dispatch(
          //   resetLoadedProblemId(`${arg.param.problemId}:${arg.param.version}`)
          // );
          return {
            param: arg.param,
            response: response,
          } as ProblemPropertyUpdateResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface ProblemUpdateState {
  value: number;
  loading: boolean;
}

const initialState: ProblemUpdateState = {
  value: 0,
  loading: false,
};

export const problemUpdateSlice = createSlice({
  name: 'problem_update',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(updateProblem.pending, (state) => {
        state.loading = true;
      })
      .addCase(updateProblem.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(updateProblem.rejected, (state) => {
        state.loading = false;
      });
  },
});

export const selectProblemEditLoading = (state: RootState) =>
  state.problemUpdate.loading;

export const selectProblemUpdateLoading = (state: RootState) =>
  state.problemUpdate.loading;

export default problemUpdateSlice.reducer;
