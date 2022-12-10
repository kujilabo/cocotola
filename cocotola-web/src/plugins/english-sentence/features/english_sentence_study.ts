import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import { RootState } from '@/app/store';
import { RecordbookModel } from '@/models/recordbook';

export const ENGLISH_SENTENCE_STATUS_INIT = 1;
export const ENGLISH_SENTENCE_STATUS_QUESTION = 2;
export const ENGLISH_SENTENCE_STATUS_ANSWER = 3;

const removeProblemFromRecordbook = (
  recordbook: RecordbookModel,
  problemId: number
): RecordbookModel => {
  const newRecordbook: RecordbookModel = {
    id: recordbook.id,
    records: [],
  };
  // onsole.log('newRecordbook', newRecordbook);
  for (const result of recordbook.records) {
    if (result.problemId === problemId) {
      continue;
    }
    newRecordbook.records.push({
      problemId: result.problemId,
      level: result.level,
      resultPrev1: result.resultPrev1,
      lastAnsweredAt: result.lastAnsweredAt,
      mastered: result.mastered,
      isReview: result.isReview,
      reviewLevel: result.reviewLevel,
    });
  }
  return newRecordbook;
};

// Set result
export interface EnglishSentenceState {
  status: number;
  ts: string;
  recordbook: RecordbookModel;
  lastResult: boolean;
}
const defaultRecordbook: RecordbookModel = {
  id: 0,
  records: [],
};
const initialState: EnglishSentenceState = {
  status: 0,
  ts: '',
  recordbook: defaultRecordbook,
  lastResult: false,
};
export const englishSentenceSlice = createSlice({
  name: 'english_word',
  initialState: initialState,
  reducers: {
    initEnglishSentenceStatus: (state, action: PayloadAction<string>) => {
      // onsole.log('initEnglishSentenceStatus', action.payload);
      state.status = ENGLISH_SENTENCE_STATUS_INIT;
      state.ts = action.payload;
    },
    setEnglishSentenceRecordbook: (
      state,
      action: PayloadAction<RecordbookModel>
    ) => {
      // onsole.log('setEnglishSentenceRecordbook', action.payload);
      const recordbook = action.payload;
      const records = [];
      for (const record of recordbook.records) {
        if (record.mastered) {
          continue;
        }

        if (record.resultPrev1) {
          const lastAnsweredAt = new Date(record.lastAnsweredAt);

          const daysToAdd = 1;
          const nextDateToAnwswer = new Date(
            lastAnsweredAt.getTime() + daysToAdd * 24 * 60 * 60 * 1000
          );
          console.log('nextDateToAnwswer', nextDateToAnwswer);
          if (new Date().getTime() < nextDateToAnwswer.getTime()) {
            continue;
          }
        }

        records.push({
          problemId: record.problemId,
          level: record.level,
          resultPrev1: record.resultPrev1,
          mastered: false,
          lastAnsweredAt: record.lastAnsweredAt,
          isReview: false,
          reviewLevel: 0,
        });
      }
      state.recordbook.records = records;
      state.status = ENGLISH_SENTENCE_STATUS_QUESTION;
      // onsole.log('setEnglishSentenceRecordbook', state.recordbook);
    },
    setEnglishSentenceStatus: (state, action: PayloadAction<number>) => {
      // onsole.log('setEnglishSentenceStatus', state.status);
      state.status = action.payload;
    },
    nextEnglishSentenceProblem: (state) => {
      // onsole.log('oldRecordbook', state);
      // onsole.log('oldRecordbook', state.recordbook);
      // onsole.log('state.recordbook.id', state.recordbook.id);
      const record = state.recordbook.records[0];
      // onsole.log('result', record);
      // onsole.log('length', state.recordbook.records.length);
      const problemId = record.problemId;
      // onsole.log('problemId', problemId);
      const level = record.level;

      const newRecordbook: RecordbookModel = removeProblemFromRecordbook(
        state.recordbook,
        problemId
      );

      // onsole.log('record.isReview', record.isReview);
      if (record.isReview) {
        const newIndexList = [3, 5, 15, 20, 30];
        if (state.lastResult) {
          record.reviewLevel++;
        } else {
          record.reviewLevel--;
        }
        record.reviewLevel = Math.min(record.reviewLevel, 4);
        record.reviewLevel = Math.max(record.reviewLevel, 0);

        if (record.reviewLevel !== 4) {
          const newIndex = Math.min(
            newIndexList[record.reviewLevel],
            newRecordbook.records.length - 1
          );
          // onsole.log('newIndex', newIndex);
          // Add problem to records
          newRecordbook.records.splice(newIndex, 0, {
            problemId,
            level,
            resultPrev1: state.lastResult,
            mastered: false,
            lastAnsweredAt: record.lastAnsweredAt,
            isReview: true,
            reviewLevel: record.reviewLevel,
          });
        }
      } else {
        if (!state.lastResult) {
          const newIndex = Math.min(3, newRecordbook.records.length - 1);
          // Add problem to records
          newRecordbook.records.splice(newIndex, 0, {
            problemId,
            level,
            resultPrev1: false,
            mastered: false,
            lastAnsweredAt: record.lastAnsweredAt,
            isReview: true,
            reviewLevel: 0,
          });
        }
      }

      state.recordbook = newRecordbook;
      state.status = ENGLISH_SENTENCE_STATUS_QUESTION;
    },
    setEnglishSentenceRecord: (state, action: PayloadAction<boolean>) => {
      state.lastResult = action.payload;
    },
  },
});

export const {
  initEnglishSentenceStatus,
  setEnglishSentenceRecordbook,
  setEnglishSentenceStatus,
  nextEnglishSentenceProblem,
  setEnglishSentenceRecord,
} = englishSentenceSlice.actions;

export const selectEnglishSentenceStatus = (state: RootState): number =>
  state.englishSentence.status;

export const selectEnglishSentenceRecordbook = (
  state: RootState
): RecordbookModel => state.englishSentence.recordbook;

export const selectTs = (state: RootState): string => state.englishSentence.ts;

export default englishSentenceSlice.reducer;
