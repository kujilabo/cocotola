import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { backendUrl, extractErrorMessage } from '@/features/base';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${backendUrl}/v1/workbook`;

// Remove problem
export type ProblemRemoveParameter = {
  workbookId: number;
  problemId: number;
  version: number;
};
export type ProblemRemoveArg = {
  param: ProblemRemoveParameter;
  postSuccessProcess: (id: number) => void;
  postFailureProcess: (error: string) => void;
};
type ProblemRemoveResponse = {
  id: number;
};
type ProblemRemoveResult = {
  param: ProblemRemoveParameter;
  response: ProblemRemoveResponse;
};

export const removeProblem = createAsyncThunk<
  ProblemRemoveResult,
  ProblemRemoveArg,
  BaseThunkApiConfig
>('problem/remove', async (arg: ProblemRemoveArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.workbookId}/problem/${arg.param.problemId}?version=${arg.param.version}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .delete(url, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as ProblemRemoveResponse;
          arg.postSuccessProcess(response.id);
          return {
            param: arg.param,
            response: response,
          } as ProblemRemoveResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface ProblemRemoveState {
  loading: boolean;
}

const initialState: ProblemRemoveState = {
  loading: false,
};

export const problemRemoveSlice = createSlice({
  name: 'problem_remove',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(removeProblem.pending, (state) => {
        state.loading = true;
      })
      .addCase(removeProblem.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(removeProblem.rejected, (state) => {
        state.loading = false;
      });
  },
});

export const selectProblemEditLoading = (state: RootState): boolean =>
  state.problemRemove.loading;

export default problemRemoveSlice.reducer;
