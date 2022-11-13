import { ReactElement } from 'react';

import { ActionCreatorWithPayload, Reducer } from '@reduxjs/toolkit';

import { ProblemModel } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';

export abstract class CustomProblem {
  abstract getName(): string;

  abstract getReducer(): Reducer;

  abstract createMenu(
    init: (s: string) => void,
    workbook: WorkbookModel
  ): ReactElement;

  abstract createReadOnlyProblem(
    id: number,
    workbookId: number,
    problem: ProblemModel
  ): // getAudio: (
  // id: number,
  // timestamp: string,
  // postFunc: (value: string) => void
  // ) => void
  ReactElement;

  abstract createReadWriteProblem(
    id: number,
    workbookId: number,
    problem: ProblemModel
  ): // getAudio: (
  // id: number,
  // timestamp: string,
  // postFunc: (value: string) => void
  // ) => void
  ReactElement;

  abstract createProblemNew(workbook: WorkbookModel): ReactElement;

  abstract createProblemEdit(
    workbook: WorkbookModel,
    problem: ProblemModel
  ): ReactElement;

  abstract createProblemStudy(studyType: string): ReactElement;

  abstract initProblemStudy(): ActionCreatorWithPayload<string>;
}
