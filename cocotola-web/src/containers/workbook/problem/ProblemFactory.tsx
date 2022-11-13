import { ReactElement } from 'react';

import {
  ActionCreatorWithPayload,
  PayloadAction,
  createSlice,
} from '@reduxjs/toolkit';

import { CustomProblem } from '@/containers/workbook/problem/CustomProblem';
import { EnglishWordProblemTypeId, ProblemModel } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';

interface EmptyState {
  status: number;
}
const initialState: EmptyState = {
  status: 0,
};
export const emptySlice = createSlice({
  name: 'empty',
  initialState: initialState,
  reducers: {
    nothing: (state) => {
      state.status = 0;
    },
    stringPayload: (state, action: PayloadAction<string>) => {
      state.status = 0;
    },
  },
});
export class ProblemFactory {
  private problemMap: { [problemType: string]: CustomProblem };

  constructor(problemMap: { [problemType: string]: CustomProblem }) {
    this.problemMap = problemMap;

    const problem = this.problemMap[EnglishWordProblemTypeId];
    if (!problem) {
      console.log('problem not found', this.problemMap);
    }
  }

  createMenu(
    problemType: string,
    init: (s: string) => void,
    workbook: WorkbookModel
  ): ReactElement {
    if (problemType === '') {
      return <div></div>;
    }
    const problemCreator = this.problemMap[problemType];
    if (!problemCreator) {
      console.log('problemCreator not found', problemType, this.problemMap);
      return <div>problemCreator not found</div>;
    }
    return problemCreator.createMenu(init, workbook);
  }

  createReadOnlyProblem(
    problemType: string,
    id: number,
    workbookId: number,
    problem: ProblemModel
    // getAudio: (
    //   id: number,
    //   timestamp: string,
    //   postFunc: (value: string) => void
    // ) => void
  ): ReactElement {
    const problemCreator = this.problemMap[problemType];
    if (!problemCreator) {
      console.log('problemCreator not found', problemType, this.problemMap);
      return <div></div>;
    }
    return problemCreator.createReadOnlyProblem(id, workbookId, problem);
  }

  createReadWriteProblem(
    problemType: string,
    id: number,
    workbookId: number,
    problem: ProblemModel
    // getAudio: (
    //   id: number,
    //   timestamp: string,
    //   postFunc: (value: string) => void
    // ) => void
  ): ReactElement {
    const problemCreator = this.problemMap[problemType];
    if (!problemCreator) {
      console.log('problemCreator not found', problemType, this.problemMap);
      return <div></div>;
    }
    return problemCreator.createReadWriteProblem(id, workbookId, problem);
  }

  createProblemNew(problemType: string, workbook: WorkbookModel): ReactElement {
    const problemCreator = this.problemMap[problemType];
    if (!problemCreator) {
      console.log('problemCreator not found', problemType, this.problemMap);
      return <div></div>;
    }
    return problemCreator.createProblemNew(workbook);
  }

  createProblemEdit(
    problemType: string,
    workbook: WorkbookModel,
    problem: ProblemModel
  ): ReactElement {
    const problemCreator = this.problemMap[problemType];
    if (!problemCreator) {
      console.log('problemCreator not found', problemType, this.problemMap);
      return <div></div>;
    }
    return problemCreator.createProblemEdit(workbook, problem);
  }

  createProblemStudy(problemType: string, studyType: string): ReactElement {
    console.log('createProblemStudy 1');
    const problemCreator = this.problemMap[problemType];
    if (!problemCreator) {
      console.log('problemCreator not found', problemType, this.problemMap);
      return <div></div>;
    }
    return problemCreator.createProblemStudy(studyType);
  }

  initProblemStudy(
    problemType: string
  ): ActionCreatorWithPayload<string, string> {
    // onsole.log('initProblemStudy 1');

    const problemCreator = this.problemMap[problemType];
    if (!problemCreator) {
      console.log('problemCreator not found', problemType, this.problemMap);
      return emptySlice.actions.stringPayload;
    }
    // onsole.log('initProblemStudy 2');
    return problemCreator.initProblemStudy();
  }
}
