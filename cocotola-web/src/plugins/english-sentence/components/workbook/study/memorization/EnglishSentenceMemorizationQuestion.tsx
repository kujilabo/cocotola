import { ReactElement, FC, useState } from 'react';

import { useParams } from 'react-router-dom';
import {
  Accordion,
  Button,
  Container,
  Divider,
  Icon,
  Message,
} from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppBreadcrumbLink, ErrorMessage } from '@/components';
import { selectProblemMap } from '@/features/problem_find';
import { addRecord } from '@/features/record_add';
import { selectWorkbook } from '@/features/workbook_get';

import {
  selectEnglishSentenceRecordbook,
  setEnglishSentenceStatus,
  setEnglishSentenceRecord,
  ENGLISH_SENTENCE_STATUS_ANSWER,
} from '../../../../features/english_sentence_study';

import { EnglishSentenceMemorizationBreadcrumb } from './EnglishSentenceMemorizationBreadcrumb';
import { EnglishSentenceMemorizationCard } from './EnglishSentenceMemorizationCard';

type ParamTypes = {
  _workbookId: string;
  _studyType: string;
};

type EnglishSentenceMemorizationQuestionProps = {
  breadcrumbLinks: AppBreadcrumbLink[];
  workbookUrl: string;
};

export const EnglishSentenceMemorizationQuestion: FC<
  EnglishSentenceMemorizationQuestionProps
> = (props: EnglishSentenceMemorizationQuestionProps): ReactElement => {
  console.log('EnglishSentenceMemorizationQuestion');
  const { _workbookId, _studyType } = useParams<ParamTypes>();
  const workbookId = +(_workbookId || '');
  const studyType = _studyType || '';
  const dispatch = useAppDispatch();
  const workbook = useAppSelector(selectWorkbook);
  const problemMap = useAppSelector(selectProblemMap);
  const englishSentenceRecordbook = useAppSelector(
    selectEnglishSentenceRecordbook
  );
  const [answerOpen, setAnswerOpen] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const breadcrumb = (
    <EnglishSentenceMemorizationBreadcrumb
      breadcrumbLinks={props.breadcrumbLinks}
      workbookUrl={props.workbookUrl}
      name={workbook.name}
      id={workbookId}
    />
  );
  if (englishSentenceRecordbook.records.length === 0) {
    return (
      <Container fluid>
        {breadcrumb}
        <Message info>
          <p>You answered all problems. Please try again in a few days.</p>
        </Message>
      </Container>
    );
  }

  const problemId = englishSentenceRecordbook.records[0].problemId;
  const problem = problemMap[problemId];
  if (!problem) {
    return <div>undefined</div>;
  }

  // onsole.log('englishSentenceRecordbook.records', englishSentenceRecordbook.records);
  // onsole.log('problemMap', problemMap);
  // onsole.log('problemId', problemId);
  // onsole.log('problem', problem);

  const setRecord = (result: boolean) => {
    dispatch(
      addRecord({
        param: {
          workbookId: workbookId,
          studyType: studyType,
          problemId: problemId,
          result: result,
          memorized: false,
        },
        postSuccessProcess: () => {
          dispatch(setEnglishSentenceRecord(result));
          dispatch(setEnglishSentenceStatus(ENGLISH_SENTENCE_STATUS_ANSWER));
        },
        postFailureProcess: setErrorMessage,
      })
    );
  };
  const onYesButtonClick = () => setRecord(true);
  const onNoButtonClick = () => setRecord(false);
  // const onMemorizeButtonClick = () => setMemorized(!memorized);
  // onsole.log('englishSentenceRecordbook', englishSentenceRecordbook);
  // onsole.log('problemId', problemId);
  // onsole.log('problem', problem);
  // onsole.log('problemMap', problemMap);
  return (
    <Container fluid>
      {breadcrumb}
      <Divider hidden />
      <EnglishSentenceMemorizationCard
        workbookId={workbookId}
        problemId={problemId}
        audioId={problem.properties['audioId']}
        updatedAt={problem.updatedAt}
        headerText={problem.properties['text']}
        contentList={[
          <div className="ui fluid buttons">
            <Button onClick={onNoButtonClick}>わからない</Button>
            <Button.Or />
            <Button positive onClick={onYesButtonClick}>
              わかる
            </Button>
          </div>,
          <Accordion>
            <Accordion.Title
              active={answerOpen}
              index={0}
              onClick={() => setAnswerOpen(!answerOpen)}
            >
              <Icon name="dropdown" />
              Answer
            </Accordion.Title>
            <Accordion.Content active={answerOpen}>
              <p>{problem.properties['translated']}</p>
            </Accordion.Content>
          </Accordion>,
        ]}
        setErrorMessage={setErrorMessage}
      />
      <ErrorMessage message={errorMessage} />
      {englishSentenceRecordbook.records.length}
      {englishSentenceRecordbook.records.map((record) => {
        const isReview = record.isReview ? 'true' : 'false';
        return (
          <div key={record.problemId}>
            {record.problemId} : {record.level} : {isReview} :{' '}
            {record.reviewLevel}
          </div>
        );
      })}
    </Container>
  );
};
