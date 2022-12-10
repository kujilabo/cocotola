import { FC, ReactElement, useState } from 'react';

import { Accordion, Button, Icon, Message } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { ErrorMessage } from '@/components';
import { selectProblemMap } from '@/features/problem_find';
import { addRecord } from '@/features/record_add';
import { WorkbookModel } from '@/models/workbook';
import {
  selectEnglishWordRecordbook,
  setEnglishWordStatus,
  setEnglishWordRecord,
  ENGLISH_WORD_STATUS_ANSWER,
} from '@/plugins/english-word/features/english_word_study';
import { EnglishWordProblemModel } from '@/plugins/english-word/models/english-word-problem';

import { EnglishWordMemorizationCard } from './EnglishWordMemorizationCard';
type EnglishWordMemorizationQuestionProps = {
  workbook: WorkbookModel;
  studyType: string;
};

export const EnglishWordMemorizationQuestion: FC<
  EnglishWordMemorizationQuestionProps
> = (props: EnglishWordMemorizationQuestionProps): ReactElement => {
  const dispatch = useAppDispatch();
  const problemMap = useAppSelector(selectProblemMap);
  const englishWordRecordbook = useAppSelector(selectEnglishWordRecordbook);
  const [answerOpen, setAnswerOpen] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  if (englishWordRecordbook.records.length === 0) {
    return (
      <Message info>
        <p>You answered all problems. Please try again in a few days.</p>
      </Message>
    );
  }

  const problemId = englishWordRecordbook.records[0].problemId;
  const problem = EnglishWordProblemModel.of(problemMap[problemId]);
  // onsole.log('englishWordRecordbook.records', englishWordRecordbook.records);
  // onsole.log('problemMap', problemMap);
  // onsole.log('problemId', problemId);
  // onsole.log('problem', problem);

  const setRecord = (result: boolean) => {
    const f = async () => {
      await dispatch(
        addRecord({
          param: {
            workbookId: props.workbook.id,
            studyType: props.studyType,
            problemId: problemId,
            result: result,
            mastered: false,
          },
          postSuccessProcess: () => {
            dispatch(setEnglishWordRecord(result));
            dispatch(setEnglishWordStatus(ENGLISH_WORD_STATUS_ANSWER));
          },
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  };
  const onYesButtonClick = () => setRecord(true);
  const onNoButtonClick = () => setRecord(false);

  if (!problem) {
    return <div>undefined</div>;
  }

  return (
    <>
      <EnglishWordMemorizationCard
        workbookId={props.workbook.id}
        problemId={problemId}
        audioId={problem.audioId}
        updatedAt={problem.updatedAt}
        headerText={problem.text}
        contentList={[
          <Button.Group fluid>
            <Button onClick={onNoButtonClick}>わからない</Button>
            <Button.Or />
            <Button positive onClick={onYesButtonClick}>
              わかる
            </Button>
          </Button.Group>,
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
              <p>{problem.translated}</p>
            </Accordion.Content>
          </Accordion>,
        ]}
        setErrorMessage={setErrorMessage}
      ></EnglishWordMemorizationCard>
      <ErrorMessage message={errorMessage} />
      {englishWordRecordbook.records.length}
      {englishWordRecordbook.records.map((record) => {
        const isReview = record.isReview ? 'true' : 'false';
        return (
          <div key={record.problemId}>
            {record.problemId} : {record.level} : {isReview} :{' '}
            {record.reviewLevel}
          </div>
        );
      })}
    </>
  );
};
