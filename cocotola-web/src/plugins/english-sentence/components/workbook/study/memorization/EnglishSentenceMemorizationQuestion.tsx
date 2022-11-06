import { FC, useState } from 'react';

import { EnglishSentenceMemorizationBreadcrumb } from './EnglishSentenceMemorizationBreadcrumb';

import { useParams } from 'react-router-dom';
import {
  Accordion,
  Button,
  Card,
  Container,
  Divider,
  Header,
  Icon,
  Message,
} from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppBreadcrumbLink, AudioButton, ErrorMessage } from '@/components';
import { getAudio, selectAudioViewLoading } from '@/features/audio';
import { selectProblemMap } from '@/features/problem_find';
import { addRecord } from '@/features/record_add';
import { selectWorkbook } from '@/features/workbook_get';
import {
  selectEnglishSentenceRecordbook,
  setEnglishSentenceStatus,
  setEnglishSentenceRecord,
  ENGLISH_SENTENCE_STATUS_ANSWER,
} from '@/plugins/english-sentence/features/english_sentence_study';
import { emptyFunction } from '@/utils/util';

type ParamTypes = {
  _workbookId: string;
  _studyType: string;
};
export const EnglishSentenceMemorizationQuestion: FC<
  EnglishSentenceMemorizationQuestionProps
> = (props: EnglishSentenceMemorizationQuestionProps) => {
  //onsole.log('EnglishSentenceMemorizationQuestion');
  const { _workbookId, _studyType } = useParams<ParamTypes>();
  const dispatch = useAppDispatch();
  const workbook = useAppSelector(selectWorkbook);
  const problemMap = useAppSelector(selectProblemMap);
  const audioViewLoading = useAppSelector(selectAudioViewLoading);
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
      id={+(_workbookId || '')}
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
  // onsole.log('englishSentenceRecordbook.records', englishSentenceRecordbook.records);
  // onsole.log('problemMap', problemMap);
  // onsole.log('problemId', problemId);
  // onsole.log('problem', problem);

  const loadAndPlay = (postFunc: (value: string) => void) => {
    dispatch(
      getAudio({
        param: {
          workbookId: +(_workbookId || ''),
          problemId: problemId,
          audioId: problem.properties['audioId'],
          updatedAt: problem.updatedAt,
        },
        postFunc: postFunc,
        postSuccessProcess: emptyFunction,
        postFailureProcess: setErrorMessage,
      })
    );
  };

  const setRecord = (result: boolean) => {
    dispatch(
      addRecord({
        param: {
          workbookId: +(_workbookId || ''),
          studyType: _studyType || '',
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
  if (!problem) {
    return <div>undefined</div>;
  }

  return (
    <Container fluid>
      {breadcrumb}
      <Divider hidden />
      <Card>
        <Card.Content>
          <Header component="h2">
            {problem.properties['text']}
            <AudioButton
              id={problem.properties['audioId']}
              loadAndPlay={loadAndPlay}
              disabled={audioViewLoading}
            />
          </Header>
        </Card.Content>
        {/* <Card.Content>
      <Form.Checkbox
        checked={memorized}
        label="完璧に覚えた"
        onClick={onMemorizeButtonClick}
      />
    </Card.Content> */}
        <Card.Content>
          <Button.Group fluid>
            <Button onClick={onNoButtonClick}>わからない</Button>
            <Button.Or />
            <Button positive onClick={onYesButtonClick}>
              わかる
            </Button>
          </Button.Group>
        </Card.Content>

        <Card.Content>
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
          </Accordion>
        </Card.Content>
      </Card>
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

type EnglishSentenceMemorizationQuestionProps = {
  breadcrumbLinks: AppBreadcrumbLink[];
  workbookUrl: string;
};
