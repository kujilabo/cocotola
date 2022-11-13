import { FC, ReactElement } from 'react';

import { useAppSelector } from '@/app/hooks';
import { AppBreadcrumbLink, AppDimmer } from '@/components';

import {
  ENGLISH_SENTENCE_STATUS_INIT,
  ENGLISH_SENTENCE_STATUS_QUESTION,
  ENGLISH_SENTENCE_STATUS_ANSWER,
  selectEnglishSentenceStatus,
} from '../../../../features/english_sentence_study';

import { EnglishSentenceMemorizationAnswer } from './EnglishSentenceMemorizationAnswer';
import { EnglishSentenceMemorizationInit } from './EnglishSentenceMemorizationInit';
import { EnglishSentenceMemorizationQuestion } from './EnglishSentenceMemorizationQuestion';

export const EnglishSentenceMemorization: FC<
  EnglishSentenceMemorizationProps
> = (props: EnglishSentenceMemorizationProps): ReactElement => {
  console.log('EnglishSentenceMemorization');
  const status = useAppSelector(selectEnglishSentenceStatus);
  if (status === ENGLISH_SENTENCE_STATUS_INIT) {
    return (
      <EnglishSentenceMemorizationInit
        breadcrumbLinks={props.breadcrumbLinks}
        workbookUrl={props.workbookUrl}
      />
    );
  } else if (status === ENGLISH_SENTENCE_STATUS_QUESTION) {
    return (
      <EnglishSentenceMemorizationQuestion
        breadcrumbLinks={props.breadcrumbLinks}
        workbookUrl={props.workbookUrl}
      />
    );
  } else if (status === ENGLISH_SENTENCE_STATUS_ANSWER) {
    return (
      <EnglishSentenceMemorizationAnswer
        breadcrumbLinks={props.breadcrumbLinks}
        workbookUrl={props.workbookUrl}
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
};
