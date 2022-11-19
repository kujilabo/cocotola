import { FC, ReactElement, useState } from 'react';

import { Accordion, Button, Icon, Message } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { ErrorMessage } from '@/components';
import { selectProblemMap } from '@/features/problem_find';
import { addRecord } from '@/features/record_add';
import { WorkbookModel } from '@/models/workbook';
import {
  selectEnglishSentenceRecordbook,
  setEnglishSentenceStatus,
  setEnglishSentenceRecord,
  ENGLISH_SENTENCE_STATUS_ANSWER,
} from '@/plugins/english-sentence/features/english_sentence_study';
import { EnglishSentenceProblemModel } from '@/plugins/english-sentence/models/english-sentence-problem';

import { EnglishSentenceMemorizationCard } from './EnglishSentenceMemorizationCard';

type EnglishSentenceMemorizationQuestionProps = {
  workbook: WorkbookModel;
  studyType: string;
};

export const EnglishSentenceMemorizationQuestion: FC<
  EnglishSentenceMemorizationQuestionProps
> = (props: EnglishSentenceMemorizationQuestionProps): ReactElement => {
  console.log('EnglishSentenceMemorizationQuestion');
  const dispatch = useAppDispatch();
  const problemMap = useAppSelector(selectProblemMap);
  const englishSentenceRecordbook = useAppSelector(
    selectEnglishSentenceRecordbook
  );
  const [answerOpen, setAnswerOpen] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  if (englishSentenceRecordbook.records.length === 0) {
    return (
      <Message info>
        <p>You answered all problems. Please try again in a few days.</p>
      </Message>
    );
  }

  const problemId = englishSentenceRecordbook.records[0].problemId;
  const problem = EnglishSentenceProblemModel.of(problemMap[problemId]);
  if (!problem) {
    return <div>undefined</div>;
  }

  const setRecord = (result: boolean) => {
    const f = async () => {
      await dispatch(
        addRecord({
          param: {
            workbookId: props.workbook.id,
            studyType: props.studyType,
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
    f().catch(console.error);
  };
  const onYesButtonClick = () => setRecord(true);
  const onNoButtonClick = () => setRecord(false);

  return (
    <>
      <EnglishSentenceMemorizationCard
        workbookId={props.workbook.id}
        problemId={problemId}
        audioId={+problem.audioId}
        updatedAt={problem.updatedAt}
        headerText={String(problem.text)}
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
              <p>{problem.translated}</p>
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
    </>
  );
};
