import React, { useEffect } from 'react';

import { useParams } from 'react-router-dom';
import { Container } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppBreadcrumbLink } from '@/components';
import { selectProblemMap } from '@/features/problem_find';
import { selectRecordbook } from '@/features/recordbook_get';
import { selectWorkbook } from '@/features/workbook_get';

import {
  setEnglishWordRecordbook,
  selectTs,
  selectEnglishWordRecordbook,
} from '../../../../features/english_word_study';

import { EnglishWordMemorizationBreadcrumb } from './EnglishWordMemorizationBreadcrumb';

type ParamTypes = {
  _workbookId: string;
  _studyType: string;
};
export const EnglishWordMemorizationInit: React.FC<
  EnglishWordMemorizationInitProps
> = (props: EnglishWordMemorizationInitProps) => {
  const { _workbookId } = useParams<ParamTypes>();
  const workbookId = +(_workbookId || '');
  const dispatch = useAppDispatch();
  const workbook = useAppSelector(selectWorkbook);
  const recordbook = useAppSelector(selectRecordbook);
  const problemMap = useAppSelector(selectProblemMap);
  const englishWordRecordbook = useAppSelector(selectEnglishWordRecordbook);
  const ts = useAppSelector(selectTs);

  console.log('EnglishWordMemorizationInit');
  useEffect(() => {
    dispatch(setEnglishWordRecordbook(recordbook));
  }, [dispatch, ts, recordbook]);

  if (englishWordRecordbook.records.length === 0) {
    return <div />;
  }
  const problemId = englishWordRecordbook.records[0].problemId;
  const problem = problemMap[problemId];
  // onsole.log('englishWordRecordbook.records', englishWordRecordbook.records);
  // onsole.log('problemMap', problemMap);
  // onsole.log('problemId', problemId);
  // onsole.log('problem', problem);

  if (!problem) {
    return <div>undefined</div>;
  }
  return (
    <Container fluid>
      <EnglishWordMemorizationBreadcrumb
        breadcrumbLinks={props.breadcrumbLinks}
        workbookUrl={props.workbookUrl}
        name={workbook.name}
        id={workbookId}
      />
    </Container>
  );
};

type EnglishWordMemorizationInitProps = {
  breadcrumbLinks: AppBreadcrumbLink[];
  workbookUrl: string;
};
