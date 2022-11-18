import { FC, ReactElement } from 'react';

import { Container, Divider } from 'semantic-ui-react';

import { useAppSelector } from '@/app/hooks';
import { AppBreadcrumbLink, AppDimmer } from '@/components';
import { selectWorkbook } from '@/features/workbook_get';
import { WorkbookModel } from '@/models/workbook';
import {
  ENGLISH_SENTENCE_STATUS_INIT,
  ENGLISH_SENTENCE_STATUS_QUESTION,
  ENGLISH_SENTENCE_STATUS_ANSWER,
  selectEnglishSentenceStatus,
} from '@/plugins/english-sentence/features/english_sentence_study';

import { EnglishSentenceMemorizationAnswer } from './EnglishSentenceMemorizationAnswer';
import { EnglishSentenceMemorizationBreadcrumb } from './EnglishSentenceMemorizationBreadcrumb';
import { EnglishSentenceMemorizationInit } from './EnglishSentenceMemorizationInit';
import { EnglishSentenceMemorizationQuestion } from './EnglishSentenceMemorizationQuestion';

const getMain = (
  status: number,
  workbook: WorkbookModel,
  studyType: string
): ReactElement => {
  if (status === ENGLISH_SENTENCE_STATUS_INIT) {
    return <EnglishSentenceMemorizationInit />;
  } else if (status === ENGLISH_SENTENCE_STATUS_QUESTION) {
    return (
      <EnglishSentenceMemorizationQuestion
        workbook={workbook}
        studyType={studyType}
      />
    );
  } else if (status === ENGLISH_SENTENCE_STATUS_ANSWER) {
    return (
      <EnglishSentenceMemorizationAnswer
        workbook={workbook}
        studyType={studyType}
      />
    );
  } else {
    console.log('status', status);
    return <AppDimmer />;
  }
};

type EnglishSentenceMemorizationProps = {
  breadcrumbLinks: AppBreadcrumbLink[];
  workbookUrl: string;
  studyType: string;
};

export const EnglishSentenceMemorization: FC<
  EnglishSentenceMemorizationProps
> = (props: EnglishSentenceMemorizationProps): ReactElement => {
  console.log('EnglishSentenceMemorization');
  const status = useAppSelector(selectEnglishSentenceStatus);
  const workbook = useAppSelector(selectWorkbook);

  const breadcrumb = (
    <EnglishSentenceMemorizationBreadcrumb
      breadcrumbLinks={props.breadcrumbLinks}
      workbookUrl={props.workbookUrl}
      name={workbook.name}
      id={workbook.id}
    />
  );
  const main = getMain(status, workbook, props.studyType);

  return (
    <Container fluid>
      {breadcrumb}
      <Divider hidden />
      {main}
    </Container>
  );
};
