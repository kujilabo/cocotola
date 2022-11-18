import { ReactElement } from 'react';

import { ActionCreatorWithPayload, Reducer } from '@reduxjs/toolkit';

import { CustomProblem } from '@/containers/workbook/problem/CustomProblem';
import { ProblemModel } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';

import {
  englishSentenceSlice,
  initEnglishSentenceStatus,
} from '../../..//features/english_sentence_study';
import { EnglishSentenceProblemReadOnly } from '../../../components/workbook/problem/EnglishSentenceProblemReadOnly';
import { EnglishSentenceProblemReadWrite } from '../../../components/workbook/problem/EnglishSentenceProblemReadWrite';
import { EnglishSentenceMemorization } from '../../../components/workbook/study/memorization/EnglishSentenceMemorization';

import { EnglishSentenceProblemEdit } from './EnglishSentenceProblemEdit';
import { EnglishSentenceProblemMenu } from './EnglishSentenceProblemMenu';
import { EnglishSentenceProblemNew } from './EnglishSentenceProblemNew';

export class EnglishSentenceProblem extends CustomProblem {
  getName(): string {
    return 'english_sentence';
  }

  getReducer(): Reducer {
    return englishSentenceSlice.reducer;
  }

  createMenu(init: (s: string) => void, workbook: WorkbookModel): ReactElement {
    return <EnglishSentenceProblemMenu initStudy={init} workbook={workbook} />;
  }

  createReadOnlyProblem(
    id: number,
    workbookId: number,
    problem: ProblemModel
  ): ReactElement {
    return (
      <EnglishSentenceProblemReadOnly
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
  ): ReactElement {
    return (
      <EnglishSentenceProblemReadWrite
        key={id}
        workbookId={workbookId}
        problem={problem}
        baseWorkbookPath={`/app/private/workbook/${workbookId}`}
      />
    );
  }

  createProblemNew(workbook: WorkbookModel): ReactElement {
    return <EnglishSentenceProblemNew workbook={workbook} />;
  }

  createProblemEdit(
    workbook: WorkbookModel,
    problem: ProblemModel
  ): ReactElement {
    return <EnglishSentenceProblemEdit workbook={workbook} problem={problem} />;
  }

  createProblemStudy(studyType: string): ReactElement {
    if (studyType === 'memorization') {
      console.log('createProblemStudy.memorization');
      return (
        <EnglishSentenceMemorization
          breadcrumbLinks={[
            { url: '/app/private/workbook', text: 'My Workbooks' },
          ]}
          workbookUrl={'/app/private/workbook/'}
          studyType={studyType}
        />
      );
    } else {
      // return <EnglishSentenceDictation />;
      return <div>xxx</div>;
    }
  }

  initProblemStudy(): ActionCreatorWithPayload<string> {
    // onsole.log('eng initProblemStudy');
    return initEnglishSentenceStatus;
  }
}
