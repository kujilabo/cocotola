import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { extractErrorMessage } from '@/features/base';
import { jsonRequestConfig } from '@/utils/util';

import { TatoebaSentencePairModel } from '../models/tatoeba';

const baseUrl = process.env.REACT_APP_BACKEND + '/plugin/tatoeba';

// Find tatoeba sentence
export type TatoebaFindParameter = {
  pageNo: number;
  pageSize: number;
  keyword: string;
  random: boolean;
};
export type TatoebaFindArg = {
  param: TatoebaFindParameter;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};
type TatoebaFindResponse = {
  results: TatoebaSentencePairModel[];
};
type TatoebaFindResult = {
  param: TatoebaFindParameter;
  response: TatoebaFindResponse;
};

export const findTatoebaSentences = createAsyncThunk<
  TatoebaFindResult,
  TatoebaFindArg,
  BaseThunkApiConfig
>('tatoeba/find', async (arg: TatoebaFindArg, thunkAPI) => {
  const url = `${baseUrl}/find`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .post(url, arg.param, jsonRequestConfig(accessToken))
        .then((resp) => {
          const response = resp.data as TatoebaFindResponse;
          arg.postSuccessProcess();
          return {
            param: arg.param,
            response: response,
          } as TatoebaFindResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface TatoebaFindState {
  value: number;
  loading: boolean;
  failed: boolean;
  sentences: TatoebaSentencePairModel[];
}

const initialState: TatoebaFindState = {
  value: 0,
  loading: false,
  failed: false,
  sentences: [],
};

export const tatoebaSentenceFindSlice = createSlice({
  name: 'translation_find',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(findTatoebaSentences.pending, (state) => {
        state.loading = true;
      })
      .addCase(findTatoebaSentences.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
        state.sentences = action.payload.response.results || [];
      })
      .addCase(findTatoebaSentences.rejected, (state, action) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
        state.sentences = [];
      });
  },
});

export const selectTatoebaFindLoading = (state: RootState): boolean =>
  state.tatoebaSentenceFindSlice.loading;

export const selectTatoebaFindFailed = (state: RootState): boolean =>
  state.tatoebaSentenceFindSlice.failed;

export const selectTatoebaSentences = (
  state: RootState
): TatoebaSentencePairModel[] => state.tatoebaSentenceFindSlice.sentences;

export default tatoebaSentenceFindSlice.reducer;
