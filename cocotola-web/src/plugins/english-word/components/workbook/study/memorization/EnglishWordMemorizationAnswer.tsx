import { FC, useState } from 'react';

import { useTranslation } from 'react-i18next';
import { Button, Form, Message } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { ErrorMessage } from '@/components';
import { LinkButton } from '@/components/buttons';
import { selectProblemMap } from '@/features/problem_find';
import { addRecord } from '@/features/record_add';
import { WorkbookModel } from '@/models/workbook';
import {
  selectEnglishWordRecordbook,
  nextEnglishWordProblem,
} from '@/plugins/english-word/features/english_word_study';
import { EnglishWordProblemModel } from '@/plugins/english-word/models/english-word-problem';
import { toDsiplayText } from '@/plugins/english-word/utils/util';
import { emptyFunction } from '@/utils/util';

import { EnglishWordMemorizationCard } from './EnglishWordMemorizationCard';

type EnglishWordMemorizationAnswerProps = {
  workbook: WorkbookModel;
  studyType: string;
};

export const EnglishWordMemorizationAnswer: FC<
  EnglishWordMemorizationAnswerProps
> = (props: EnglishWordMemorizationAnswerProps) => {
  const dispatch = useAppDispatch();
  const [t] = useTranslation();
  const problemMap = useAppSelector(selectProblemMap);
  const englishWordRecordbook = useAppSelector(selectEnglishWordRecordbook);
  const [mastered, setMemorized] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const problemId = englishWordRecordbook.records[0].problemId;
  const problem = EnglishWordProblemModel.of(problemMap[problemId]);
  const baseUrl = `/app/private/workbook/${props.workbook.id}/problem/${problemId}`;
  console.log('problem', problem);
  // let sentence1 = emptyTatoebaSentence;
  // let sentence2 = emptyTatoebaSentence;
  // if (problem.senexampleSentenceNote && values.exampleSentenceNote !== '') {
  //   try {
  //     const noteObj = JSON.parse(values.exampleSentenceNote);
  //     console.log('noteObj', noteObj);
  //     sentence1 = {
  //       text: values.exampleSentenceText,
  //       author: noteObj['tatoebaAuthor1'],
  //       sentenceNumber: +noteObj['tatoebaSentenceNumber1'],
  //       lang: 'en',
  //     };
  //     sentence2 = {
  //       text: values.exampleSentenceTranslated,
  //       author: noteObj['tatoebaAuthor2'],
  //       sentenceNumber: +noteObj['tatoebaSentenceNumber2'],
  //       lang: 'ja',
  //     };
  //   } catch (e) {
  //     console.log(e);
  //   }
  // }

  if (englishWordRecordbook.records.length === 0) {
    return (
      <Message info>
        <p>You answered all problems. Please try again in a few days.</p>
      </Message>
    );
  }

  const onNextButtonClick = () => {
    if (mastered) {
      const f = async () => {
        await dispatch(
          addRecord({
            param: {
              workbookId: props.workbook.id,
              studyType: props.studyType,
              problemId: problemId,
              result: true,
              mastered: true,
            },
            postSuccessProcess: emptyFunction,
            postFailureProcess: setErrorMessage,
          })
        );
      };
      f().catch(console.error);
    }
    dispatch(nextEnglishWordProblem());
  };
  const onMemorizeButtonClick = () => setMemorized(!mastered);

  return (
    <>
      <EnglishWordMemorizationCard
        workbookId={props.workbook.id}
        problemId={problemId}
        audioId={problem.audioId}
        updatedAt={problem.updatedAt}
        headerText={problem.text}
        contentList={[
          <p>{toDsiplayText(+problem.pos)}</p>,
          <p>{problem.translated}</p>,
          <Button.Group fluid>
            <LinkButton to={`${baseUrl}/edit`} value={t('Edit')} />
          </Button.Group>,
          <Form.Checkbox
            checked={mastered}
            label="完璧に覚えた"
            onClick={onMemorizeButtonClick}
          />,

          <Button.Group fluid>
            <Button color="teal" onClick={onNextButtonClick}>
              Next
            </Button>
          </Button.Group>,
        ]}
        setErrorMessage={setErrorMessage}
      ></EnglishWordMemorizationCard>
      <ErrorMessage message={errorMessage} />
      {englishWordRecordbook.records.map((record) => {
        return (
          <div key={record.problemId}>
            {record.problemId} : {record.level} : {record.reviewLevel}
          </div>
        );
      })}
    </>
  );
};
