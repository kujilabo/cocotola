import { ReactElement } from 'react';

import { EnglishSentenceProblemEdit } from './EnglishSentenceProblemEdit';
import { EnglishSentenceProblemMenu } from './EnglishSentenceProblemMenu';
import { EnglishSentenceProblemNew } from './EnglishSentenceProblemNew';

import { EnglishSentenceProblemReadOnly } from '../../../components/workbook/problem/EnglishSentenceProblemReadOnly';
import { EnglishSentenceProblemReadWrite } from '../../../components/workbook/problem/EnglishSentenceProblemReadWrite';

// import EnglishSentenceDictation from '../../../components/workbook/study/dictation/EnglishSentenceDictation';

import { ActionCreatorWithPayload, Reducer } from '@reduxjs/toolkit';

import { CustomProblem } from '@/containers/workbook/problem/CustomProblem';
import { ProblemModel } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';
import { EnglishSentenceMemorization } from '@/plugins/english-sentence/components/workbook/study/memorization/EnglishSentenceMemorization';
import {
  englishSentenceSlice,
  initEnglishSentenceStatus,
} from '@/plugins/english-sentence/features/english_sentence_study';

export class EnglishSentenceProblem extends CustomProblem {
  getName(): string {
    return 'english_word';
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
      // onsole.log('createProblemStudy.memorization');
      return (
        <EnglishSentenceMemorization
          breadcrumbLinks={[
            { url: '/app/private/workbook', text: 'My Workbooks' },
          ]}
          workbookUrl={'/app/private/workbook/'}
        />
      );
    } else {
      // return <EnglishSentenceDictation />;
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
    return initEnglishSentenceStatus;
  }
}
