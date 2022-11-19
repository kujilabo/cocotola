import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

import { RootState, BaseThunkApiConfig } from '@/app/store';
import { refreshAccessToken } from '@/features/auth';
import { AudioModel } from '@/models/audio';
import { jsonHeaders } from '@/utils/util';

import { backendUrl, extractErrorMessage } from './base';

const baseUrl = `${backendUrl}/v1/workbook`;

// Find audio
export type AudioViewParameter = {
  workbookId: number;
  problemId: number;
  audioId: number;
  updatedAt: string;
};
export type AudioViewArg = {
  param: AudioViewParameter;
  postFunc: (value: string) => void;
  postSuccessProcess: () => void;
  postFailureProcess: (error: string) => void;
};

type AudioViewResult = {
  param: AudioViewParameter;
  response: AudioModel;
};

export const getAudio = createAsyncThunk<
  AudioViewResult,
  AudioViewArg,
  BaseThunkApiConfig
>('audio/get', async (arg: AudioViewArg, thunkAPI) => {
  const url = `${baseUrl}/${arg.param.workbookId}/problem/${arg.param.problemId}/audio/${arg.param.audioId}`;
  const { refreshToken } = thunkAPI.getState().auth;
  return await thunkAPI
    .dispatch(refreshAccessToken({ refreshToken: refreshToken }))
    .then(() => {
      const { accessToken } = thunkAPI.getState().auth;
      return axios
        .get(url, { headers: jsonHeaders(accessToken), data: {} })
        .then((resp) => {
          const response = resp.data as AudioModel;
          arg.postSuccessProcess();
          arg.postFunc(response.content);
          return { param: arg.param, response: response } as AudioViewResult;
        })
        .catch((err: Error) => {
          const errorMessage = extractErrorMessage(err);
          arg.postFailureProcess(errorMessage);
          return thunkAPI.rejectWithValue(errorMessage);
        });
    });
});

export interface AudioViewState {
  loading: boolean;
  failed: boolean;
  audio: AudioModel;
}
const defaultAudio: AudioModel = {
  id: 0,
  lang2: '',
  text: '',
  content: '',
};
const initialState: AudioViewState = {
  loading: false,
  failed: false,
  audio: defaultAudio,
};

export const audioSlice = createSlice({
  name: 'audio_get',
  initialState: initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(getAudio.pending, (state) => {
        state.loading = true;
      })
      .addCase(getAudio.fulfilled, (state, action) => {
        state.loading = false;
        state.failed = false;
        state.audio = action.payload.response;
      })
      .addCase(getAudio.rejected, (state) => {
        // onsole.log('rejected', action);
        state.loading = false;
        state.failed = true;
      });
  },
});

export const selectAudioViewLoading = (state: RootState) => state.audio.loading;

export const selectAudioListFailed = (state: RootState) => state.audio.failed;

export const selectAudio = (state: RootState) => state.audio.audio;

export default audioSlice.reducer;
