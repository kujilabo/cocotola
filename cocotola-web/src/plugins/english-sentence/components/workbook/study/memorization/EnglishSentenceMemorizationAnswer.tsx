import { FC, ReactElement, useState } from 'react';

import { Input } from 'formik-semantic-ui-react';
import { Button, Form, Message } from 'semantic-ui-react';
import * as Yup from 'yup';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { ErrorMessage } from '@/components';
import { ProblemPropertyEditFormikForm } from '@/components/problem/ProblemPropertyEditFormikForm';
import { selectProblemMap } from '@/features/problem_find';
import { addRecord } from '@/features/record_add';
import { EnglishSentenceProblemTypeId } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';
import {
  selectEnglishSentenceRecordbook,
  nextEnglishSentenceProblem,
} from '@/plugins/english-sentence/features/english_sentence_study';
import { EnglishSentenceProblemModel } from '@/plugins/english-sentence/models/english-sentence-problem';
import { emptyFunction } from '@/utils/util';

import { EnglishSentenceMemorizationCard } from './EnglishSentenceMemorizationCard';

type EnglishSentenceMemorizationAnswerProps = {
  workbook: WorkbookModel;
  studyType: string;
};

interface formikFormPropsTranslated extends Object {
  translated: string;
}

interface formValuesTranslated extends Object {
  translated: string;
}

export const EnglishSentenceMemorizationAnswer: FC<
  EnglishSentenceMemorizationAnswerProps
> = (props: EnglishSentenceMemorizationAnswerProps): ReactElement => {
  const dispatch = useAppDispatch();
  const problemMap = useAppSelector(selectProblemMap);
  const englishSentenceRecordbook = useAppSelector(
    selectEnglishSentenceRecordbook
  );
  const [mastered, setMemorized] = useState(false);
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
  console.log('problem', problem);
  console.log('problemId', problemId);
  console.log('problem.translated', problem.translated);
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
  // useEffect(() => {
  //   console.log('get problem');
  //   dispatch(
  //     getProblem({
  //       param: { workbookId: workbookId, problemId: problemId },
  //       postSuccessProcess: (p: ProblemModel) => {
  //         const e = EnglishSentenceProblemModel.of(p);
  //         console.log(e);
  //       },
  //       postFailureProcess: setErrorMessage,
  //     })
  //   );
  //   }, [dispatch, workbookId, problemId, problem.version]);

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
    dispatch(nextEnglishSentenceProblem());
  };
  const onMemorizeButtonClick = () => setMemorized(!mastered);

  const EnglishSentenceProblemEditFormikForm = ProblemPropertyEditFormikForm<
    formValuesTranslated,
    formikFormPropsTranslated
  >({
    workbookId: props.workbook.id,
    problemId: problem.id,
    problemVersion: problem.version,
    problemType: EnglishSentenceProblemTypeId,
    toField: () => (
      <Input name="translated" placeholder="translated sentence" errorPrompt />
    ),
    validationSchema: Yup.object().shape({
      translated: Yup.string().required('Translated is required'),
    }),
    propsToValues: (props: formikFormPropsTranslated) => ({ ...props }),
    valuesToProperties: (values: formValuesTranslated) => ({
      translated: values.translated,
    }),
    resetValues: () => emptyFunction,
    setErrorMessage: setErrorMessage,
  });

  return (
    <>
      <EnglishSentenceMemorizationCard
        workbookId={props.workbook.id}
        problemId={problemId}
        audioId={+problem.audioId}
        updatedAt={problem.updatedAt}
        headerText={problem.text}
        contentList={[
          <div>
            <EnglishSentenceProblemEditFormikForm
              translated={problem.translated}
            />
          </div>,
          <Form.Checkbox
            checked={mastered}
            label="完璧に覚えた"
            onClick={onMemorizeButtonClick}
          />,

          <div className="ui fluid buttons">
            <Button color="teal" onClick={onNextButtonClick}>
              Next
            </Button>
          </div>,
        ]}
        setErrorMessage={setErrorMessage}
      />
      <ErrorMessage message={errorMessage} />
      {englishSentenceRecordbook.records.map((record) => {
        return (
          <div key={record.problemId}>
            {record.problemId} : {record.level} : {record.reviewLevel}
          </div>
        );
      })}
    </>
  );
};
