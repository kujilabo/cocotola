import React from 'react';

import { useAppSelector } from 'app/hooks';
import { AppBreadcrumbLink, AppDimmer } from 'components';
import { EnglishWordMemorizationInit } from './EnglishWordMemorizationInit';
import { EnglishWordMemorizationQuestion } from './EnglishWordMemorizationQuestion';
import { EnglishWordMemorizationAnswer } from './EnglishWordMemorizationAnswer';
import {
  ENGLISH_WORD_STATUS_INIT,
  ENGLISH_WORD_STATUS_QUESTION,
  ENGLISH_WORD_STATUS_ANSWER,
  selectEnglishWordStatus,
} from 'plugins/english-word/features/english_word_study';

import 'App.css';

export const EnglishWordMemorization: React.FC<EnglishWordMemorizationProps> = (
  props: EnglishWordMemorizationProps
) => {
  console.log('EnglishWordMemorization');
  const status = useAppSelector(selectEnglishWordStatus);
  if (status === ENGLISH_WORD_STATUS_INIT) {
    return (
      <EnglishWordMemorizationInit
        breadcrumbLinks={props.breadcrumbLinks}
        workbookUrl={props.workbookUrl}
      />
    );
  } else if (status === ENGLISH_WORD_STATUS_QUESTION) {
    return (
      <EnglishWordMemorizationQuestion
        breadcrumbLinks={props.breadcrumbLinks}
        workbookUrl={props.workbookUrl}
      />
    );
  } else if (status === ENGLISH_WORD_STATUS_ANSWER) {
    return (
      <EnglishWordMemorizationAnswer
        breadcrumbLinks={props.breadcrumbLinks}
        workbookUrl={props.workbookUrl}
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
};
