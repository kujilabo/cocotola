import { FC, useEffect } from 'react';

import { useParams } from 'react-router-dom';
import { Container } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppBreadcrumbLink } from '@/components';
import { selectProblemMap } from '@/features/problem_find';
import { selectRecordbook } from '@/features/recordbook_get';
import { selectWorkbook } from '@/features/workbook_get';

import {
  setEnglishSentenceRecordbook,
  selectTs,
  selectEnglishSentenceRecordbook,
} from '../../../../features/english_sentence_study';

import { EnglishSentenceMemorizationBreadcrumb } from './EnglishSentenceMemorizationBreadcrumb';

type ParamTypes = {
  _workbookId: string;
  _studyType: string;
};

type EnglishSentenceMemorizationInitProps = {
  breadcrumbLinks: AppBreadcrumbLink[];
  workbookUrl: string;
};

export const EnglishSentenceMemorizationInit: FC<
  EnglishSentenceMemorizationInitProps
> = (props: EnglishSentenceMemorizationInitProps) => {
  const { _workbookId } = useParams<ParamTypes>();
  const workbookId = +(_workbookId || '');
  const dispatch = useAppDispatch();
  const workbook = useAppSelector(selectWorkbook);
  const recordbook = useAppSelector(selectRecordbook);
  const problemMap = useAppSelector(selectProblemMap);
  const englishSentenceRecordbook = useAppSelector(
    selectEnglishSentenceRecordbook
  );
  const ts = useAppSelector(selectTs);

  console.log('EnglishSentenceMemorizationInit');
  useEffect(() => {
    dispatch(setEnglishSentenceRecordbook(recordbook));
  }, [dispatch, ts, recordbook]);

  if (englishSentenceRecordbook.records.length === 0) {
    return <div />;
  }
  const problemId = englishSentenceRecordbook.records[0].problemId;
  const problem = problemMap[problemId];
  // onsole.log('englishSentenceRecordbook.records', englishSentenceRecordbook.records);
  // onsole.log('problemMap', problemMap);
  // onsole.log('problemId', problemId);
  // onsole.log('problem', problem);

  if (!problem) {
    return <div>undefined</div>;
  }
  return (
    <Container fluid>
      <EnglishSentenceMemorizationBreadcrumb
        breadcrumbLinks={props.breadcrumbLinks}
        workbookUrl={props.workbookUrl}
        name={workbook.name}
        id={workbookId}
      />
    </Container>
  );
};
