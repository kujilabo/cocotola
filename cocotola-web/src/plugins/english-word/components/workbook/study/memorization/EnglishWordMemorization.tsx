import { FC, ReactElement } from 'react';

import { Container, Divider } from 'semantic-ui-react';

import { useAppSelector } from '@/app/hooks';
import { AppBreadcrumbLink, AppDimmer } from '@/components';
import { selectWorkbook } from '@/features/workbook_get';
import { WorkbookModel } from '@/models/workbook';
import {
  ENGLISH_WORD_STATUS_INIT,
  ENGLISH_WORD_STATUS_QUESTION,
  ENGLISH_WORD_STATUS_ANSWER,
  selectEnglishWordStatus,
} from '@/plugins/english-word/features/english_word_study';

import { EnglishWordMemorizationAnswer } from './EnglishWordMemorizationAnswer';
import { EnglishWordMemorizationBreadcrumb } from './EnglishWordMemorizationBreadcrumb';
import { EnglishWordMemorizationInit } from './EnglishWordMemorizationInit';
import { EnglishWordMemorizationQuestion } from './EnglishWordMemorizationQuestion';

const getMain = (
  status: number,
  workbook: WorkbookModel,
  studyType: string
): ReactElement => {
  if (status === ENGLISH_WORD_STATUS_INIT) {
    return <EnglishWordMemorizationInit />;
  } else if (status === ENGLISH_WORD_STATUS_QUESTION) {
    return (
      <EnglishWordMemorizationQuestion
        workbook={workbook}
        studyType={studyType}
      />
    );
  } else if (status === ENGLISH_WORD_STATUS_ANSWER) {
    return (
      <EnglishWordMemorizationAnswer
        workbook={workbook}
        studyType={studyType}
      />
    );
  } else {
    console.log('status', status);
    return <AppDimmer />;
  }
};

type EnglishWordMemorizationProps = {
  breadcrumbLinks: AppBreadcrumbLink[];
  workbookUrl: string;
  studyType: string;
};

export const EnglishWordMemorization: FC<EnglishWordMemorizationProps> = (
  props: EnglishWordMemorizationProps
): ReactElement => {
  console.log('EnglishWordMemorization');
  const status = useAppSelector(selectEnglishWordStatus);
  const workbook = useAppSelector(selectWorkbook);

  const breadcrumb = (
    <EnglishWordMemorizationBreadcrumb
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
