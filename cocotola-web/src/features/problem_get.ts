import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { extractErrorMessage } from '@/features/base';
import { ProblemModel } from '@/models/problem';
import { jsonHeaders } from '@/utils/util';

const baseUrl = import.meta.env.VITE_APP_BACKEND + '/v1/workbook';

// Get problem
export type ProblemGetParameter = {
  workbookId: number;
  problemId: number;
};
export type ProblemGetArg = {
  param: ProblemGetParameter;
  postSuccessProcess: (problem: ProblemModel) => void;
  postFailureProcess: (error: string) => void;
};

type ProblemGetReult = {
  param: ProblemGetParameter;
  response: ProblemModel;
};

export const getProblem = createAsyncThunk<
  ProblemGetReult,
  ProblemGetArg,
  BaseThunkApiConfig
>('problem/get', async (arg: ProblemGetArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.workbookId}/problem/${arg.param.problemId}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then((resp) => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .get(url, { headers: jsonHeaders(accessToken), data: {} })
        .then((resp) => {
          const response = resp.data as ProblemModel;
          arg.postSuccessProcess(response);
          return { param: arg.param, response: response } as ProblemGetReult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface ProblemGetState {
  value: number;
  loading: boolean;
  failed: boolean;
  problem: ProblemModel;
  loadedProblemId: string;
}

const defaultProblem = {
  id: 0,
  version: 0,
  updatedAt: '',
  number: 0,
  problemType: '',
  properties: {},
};
const initialState: ProblemGetState = {
  value: 0,
  loading: false,
  failed: false,
  problem: defaultProblem,
  loadedProblemId: '',
};

export const problemGetSlice = createSlice({
  name: 'problem_get',
  initialState: initialState,
  reducers: {
    resetLoadedProblemId: (state, action: PayloadAction<string>) => {
      console.log('state.loadedProblemId =', state.loadedProblemId);
      state.loadedProblemId = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(getProblem.pending, (state) => {
        state.loading = true;
      })
      .addCase(getProblem.fulfilled, (state, action) => {
        // onsole.log('problem', action.payload.response);
        state.loading = false;
        state.failed = false;
        state.problem = action.payload.response;
      })
      .addCase(getProblem.rejected, (state, action) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
      });
  },
});

export const { resetLoadedProblemId } = problemGetSlice.actions;

export const selectProblemGetLoading = (state: RootState) =>
  state.problemGet.loading;

export const selectProblem = (state: RootState) => state.problemGet.problem;

export const selectLoadedProblemId = (state: RootState) =>
  state.problemGet.loadedProblemId;

export default problemGetSlice.reducer;
