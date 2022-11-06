import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { extractErrorMessage } from '@/features/base';
import { removeProblem } from '@/features/problem_remove';
import { ProblemModel } from '@/models/problem';
import { jsonRequestConfig } from '@/utils/util';

const baseUrl = `${import.meta.env.VITE_APP_BACKEND}/v1/workbook`;

// Find problems
export type ProblemFindParameter = {
  workbookId: number;
  pageNo: number;
  pageSize: number;
  keyword: string;
};
export type ProblemFindArg = {
  param: ProblemFindParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type ProblemFindResponse = {
  results: ProblemModel[];
  totalCount: number;
};
type ProblemFindResult = {
  param: ProblemFindParameter;
  response: ProblemFindResponse;
};

export const findProblems = createAsyncThunk<
  ProblemFindResult,
  ProblemFindArg,
  BaseThunkApiConfig
>('problem/find', async (arg: ProblemFindArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.workbookId}/problem/find`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as ProblemFindResponse;
          arg.postSuccessProcess();
          return { param: arg.param, response: response } as ProblemFindResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

// Find all problems
export type ProblemFindAllParameter = {
  workbookId: number;
};
export type ProblemFindAllArg = {
  param: ProblemFindAllParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type ProblemFindAllResponse = {
  results: ProblemModel[];
  totalCount: number;
};
type ProblemFindAllResult = {
  param: ProblemFindAllParameter;
  response: ProblemFindAllResponse;
};
export const findAllProblems = createAsyncThunk<
  ProblemFindAllResult,
  ProblemFindAllArg,
  BaseThunkApiConfig
>('problem/find_all', async (arg: ProblemFindAllArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.workbookId}/problem/find_all`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as ProblemFindAllResponse;
          arg.postSuccessProcess();
          return {
            param: arg.param,
            response: response,
          } as ProblemFindAllResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface problemFindState {
  loading: boolean;
  failed: boolean;
  problems: ProblemModel[];
  problemMap: { [key: number]: ProblemModel };
  problemsTotalCount: number;
}

const initialState: problemFindState = {
  loading: false,
  failed: false,
  problems: [],
  problemMap: {},
  problemsTotalCount: 0,
};

export const problemFindSlice = createSlice({
  name: 'problem_find',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      // findProblems
      .addCase(findProblems.pending, (state) => {
        state.loading = true;
      })
      .addCase(findProblems.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
        state.problems = action.payload.response.results;
        state.problemMap = {};
        for (let i = 0; i < state.problems.length; i++) {
          const problem = state.problems[i];
          state.problemMap[problem.id] = problem;
        }
        state.problemsTotalCount = action.payload.response.totalCount;
      })
      .addCase(findProblems.rejected, (state) => {
        state.loading = false;
        state.failed = true;
      })
      // findAllProblems
      .addCase(findAllProblems.pending, (state) => {
        state.loading = true;
      })
      .addCase(findAllProblems.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
        state.problems = action.payload.response.results;
        state.problemMap = {};
        for (let i = 0; i < state.problems.length; i++) {
          const problem = state.problems[i];
          state.problemMap[problem.id] = problem;
        }
        state.problemsTotalCount = action.payload.response.totalCount;
      })
      .addCase(findAllProblems.rejected, (state) => {
        state.loading = false;
        state.failed = true;
      })
      // removeProblem
      .addCase(removeProblem.fulfilled, (state, action) => {
        const removedProblemId = action.payload.param.problemId;
        state.problems = state.problems.filter(
          (n) => n.id !== removedProblemId
        );
        state.problemMap = {};
        for (let i = 0; i < state.problems.length; i++) {
          const problem = state.problems[i];
          state.problemMap[problem.id] = problem;
        }
        state.problemsTotalCount--;
      });
  },
});

export const selectProblemFindLoading = (state: RootState): boolean =>
  state.problemFind.loading;

export const selectProblemFindFailed = (state: RootState): boolean =>
  state.problemFind.failed;

export const selectProblems = (state: RootState): ProblemModel[] =>
  state.problemFind.problems;

export const selectProblemMap = (
  state: RootState
): { [key: number]: ProblemModel } => state.problemFind.problemMap;

export const selectProblemsTotalCount = (state: RootState): number =>
  state.problemFind.problemsTotalCount;

export default problemFindSlice.reducer;
