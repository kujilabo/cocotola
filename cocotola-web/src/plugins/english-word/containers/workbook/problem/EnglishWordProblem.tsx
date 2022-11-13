import React from 'react';
import { ActionCreatorWithPayload, Reducer } from '@reduxjs/toolkit';

import { CustomProblem } from 'containers/workbook/problem/CustomProblem';
import { ProblemModel } from 'models/problem';
import { WorkbookModel } from 'models/workbook';
import { EnglishWordProblemNew } from './EnglishWordProblemNew';
import { EnglishWordProblemEdit } from './EnglishWordProblemEdit';
import { EnglishWordProblemMenu } from './EnglishWordProblemMenu';

import { EnglishWordProblemReadOnly } from '../../../components/workbook/problem/EnglishWordProblemReadOnly';
import { EnglishWordProblemReadWrite } from '../../../components/workbook/problem/EnglishWordProblemReadWrite';
import { EnglishWordMemorization } from '../../../components/workbook/study/memorization/EnglishWordMemorization';
// import EnglishWordDictation from '../../../components/workbook/study/dictation/EnglishWordDictation';
import {
  englishWordSlice,
  initEnglishWordStatus,
} from '../../../features/english_word_study';

export class EnglishWordProblem extends CustomProblem {
  getName(): string {
    return 'english_word';
  }

  getReducer(): Reducer {
    return englishWordSlice.reducer;
  }

  createMenu(
    init: (s: string) => void,
    workbook: WorkbookModel
  ): React.ReactElement {
    return <EnglishWordProblemMenu initStudy={init} workbook={workbook} />;
  }

  createReadOnlyProblem(
    id: number,
    workbookId: number,
    problem: ProblemModel
  ): React.ReactElement {
    return (
      <EnglishWordProblemReadOnly
        key={id}
        workbookId={workbookId}
        problem={problem}
      />
    );
  }

  createReadWriteProblem(
    id: number,
    workbookId: number,
    problem: ProblemModel
  ): React.ReactElement {
    return (
      <EnglishWordProblemReadWrite
        key={id}
        workbookId={workbookId}
        problem={problem}
        baseWorkbookPath={`/app/private/workbook/${workbookId}`}
      />
    );
  }

  createProblemNew(workbook: WorkbookModel): React.ReactElement {
    return <EnglishWordProblemNew workbook={workbook} />;
  }

  createProblemEdit(
    workbook: WorkbookModel,
    problem: ProblemModel
  ): React.ReactElement {
    return <EnglishWordProblemEdit workbook={workbook} problem={problem} />;
  }

  createProblemStudy(studyType: string): React.ReactElement {
    if (studyType === 'memorization') {
      // onsole.log('createProblemStudy.memorization');
      return (
        <EnglishWordMemorization
          breadcrumbLinks={[
            { url: '/app/private/workbook', text: 'My Workbooks' },
          ]}
          workbookUrl={'/app/private/workbook/'}
        />
      );
    } else {
      // return <EnglishWordDictation />;
      return <div>xxx</div>;
    }
  }

  // createBreadcrumbs(studyType: string): React.ReactElement {
  //   if (studyType == 'memorization') {

  //   } else {

  //   }
  // }

  initProblemStudy(): ActionCreatorWithPayload<string> {
    // onsole.log('eng initProblemStudy');
    return initEnglishWordStatus;
  }
}
