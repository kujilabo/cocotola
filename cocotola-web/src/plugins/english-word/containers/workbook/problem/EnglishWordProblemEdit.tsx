import {
  Dispatch,
  FC,
  ReactElement,
  SetStateAction,
  useEffect,
  useState,
} from 'react';

import { FormikProps } from 'formik';
import { Input, Select } from 'formik-semantic-ui-react';
import { Container, Divider } from 'semantic-ui-react';
import * as Yup from 'yup';

import { useAppDispatch, useAppSelector } from '@/app/hooks';
import {
  ErrorMessage,
  PrivateProblemBreadcrumb,
  langOptions,
} from '@/components';
import { ProblemEditFormikForm } from '@/components/problem/ProblemEditFormikForm';
import { ProblemModel, EnglishWordProblemTypeId } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';
import { ExampleTatoebaSentence } from '@/plugins/english-word/components/workbook/problem/ExampleTatoebaSentence';
import { ExampleTatoebaSentenceList } from '@/plugins/english-word/components/workbook/problem/ExampleTatoebaSentenceList';
import { EnglishWordProblemModel } from '@/plugins/english-word/models/english-word-problem';
import {
  findTatoebaSentences,
  selectTatoebaSentences,
  selectTatoebaFindLoading,
} from '@/plugins/tatoeba/features/tatoeba_find';
import {
  TatoebaSentencePairModel,
  TatoebaSentenceModel,
} from '@/plugins/tatoeba/models/tatoeba';
import { posOptions } from '@/plugins/translation/components';
import { emptyFunction } from '@/utils/util';

interface formikFormProps {
  text: string;
  pos: string;
  translated: string;
  lang2: string;
  exampleSentenceText: string;
  exampleSentenceTranslated: string;
  exampleSentenceNote: string;
  sentenceProvider: string;
  tatoebaSentenceNumber1: string;
  tatoebaSentenceNumber2: string;
  tatoebaSentences: TatoebaSentencePairModel[];
}

interface formValues {
  text: string;
  pos: string;
  translated: string;
  lang2: string;
  exampleSentenceText: string;
  exampleSentenceTranslated: string;
  exampleSentenceNote: string;
  sentenceProvider: string;
  tatoebaSentenceNumber1: string;
  tatoebaSentenceNumber2: string;
  tatoebaSentences: TatoebaSentencePairModel[];
}

const editFormikForm = (
  workbookId: number,
  problemId: number,
  problemVersion: number,
  resetValues: (v: formValues) => void,
  setErrorMessage: Dispatch<SetStateAction<string>>
) => {
  const emptyTatoebaSentence: TatoebaSentenceModel = {
    text: '',
    author: '',
    sentenceNumber: 0,
    lang2: '',
  };
  return ProblemEditFormikForm({
    workbookId: workbookId,
    problemId: problemId,
    problemVersion: problemVersion,
    problemType: EnglishWordProblemTypeId,
    toContent: (props: FormikProps<formValues>) => {
      const { values, setFieldValue } = props;

      const onCheckboxChange = (
        sentenceNumber1: number,
        sentenceNumber2: number,
        checked: boolean
      ): void => {
        if (checked) {
          setFieldValue('sentenceProvider', 'tatoeba');
          setFieldValue('tatoebaSentenceNumber1', sentenceNumber1.toString());
          setFieldValue('tatoebaSentenceNumber2', sentenceNumber2.toString());
        } else {
          setFieldValue('sentenceProvider', '');
          setFieldValue('tatoebaSentenceNumber1', '');
          setFieldValue('tatoebaSentenceNumber2', '');
        }
      };

      let sentence1 = emptyTatoebaSentence;
      let sentence2 = emptyTatoebaSentence;
      if (values.exampleSentenceNote && values.exampleSentenceNote !== '') {
        try {
          /* eslint-disable */
          const noteObj: { [key: string]: string } = JSON.parse(
            values.exampleSentenceNote
          );
          /* eslint-enable */
          console.log('noteObj', noteObj);
          sentence1 = {
            text: values.exampleSentenceText,
            author: String(noteObj['tatoebaAuthor1']),
            sentenceNumber: +noteObj['tatoebaSentenceNumber1'],
            lang2: 'en',
          };
          sentence2 = {
            text: values.exampleSentenceTranslated,
            author: noteObj['tatoebaAuthor2'],
            sentenceNumber: +noteObj['tatoebaSentenceNumber2'],
            lang2: 'ja',
          };
        } catch (e) {
          console.log(e);
        }
      }

      return (
        <>
          <Input
            name="text"
            label="Sentence"
            placeholder="english sentence"
            errorPrompt
          />
          <Select
            name="pos"
            label="Pos"
            options={posOptions}
            value={values.pos}
            errorPrompt
          />
          <Select
            name="lang2"
            label="Lang"
            options={langOptions}
            value={values.lang2}
            errorPrompt
          />
          <Input
            name="translated"
            label="Translated sentence"
            placeholder="translated sentence"
            errorPrompt
          />

          {sentence1.text !== '' ? (
            <ExampleTatoebaSentence
              sentence1={sentence1}
              sentence2={sentence2}
            />
          ) : (
            <div />
          )}
          <ExampleTatoebaSentenceList
            sentences={values.tatoebaSentences}
            onCheckboxChange={onCheckboxChange}
          />
        </>
      );
    },
    validationSchema: Yup.object().shape({
      text: Yup.string().required('Sentence is required'),
    }),
    propsToValues: (props: formikFormProps) => ({ ...props }),
    valuesToProperties: (values: formValues) => ({
      text: values.text,
      pos: values.pos,
      translated: values.translated,
      lang2: values.lang2,
      sentenceProvider: values.sentenceProvider,
      tatoebaSentenceNumber1: values.tatoebaSentenceNumber1,
      tatoebaSentenceNumber2: values.tatoebaSentenceNumber2,
    }),
    resetValues: (values: formValues) => resetValues(values),
    setErrorMessage: setErrorMessage,
  });
};

type EnglishWordProblemEditProps = {
  workbook: WorkbookModel;
  problem: ProblemModel;
};

export const EnglishWordProblemEdit: FC<EnglishWordProblemEditProps> = (
  props: EnglishWordProblemEditProps
): ReactElement => {
  const dispatch = useAppDispatch();
  const problem = EnglishWordProblemModel.of(props.problem);
  const tatoebaSentences = useAppSelector(selectTatoebaSentences);
  const tatoebaSentenceFindLoading = useAppSelector(selectTatoebaFindLoading);
  const [values, setValues] = useState({
    text: problem.text,
    pos: String(problem.pos),
    lang2: problem.lang2,
    translated: problem.translated,
    exampleSentenceText: problem.sentence1.text,
    exampleSentenceTranslated: problem.sentence1.translated,
    exampleSentenceNote: problem.sentence1.note,
    sentenceProvider: '',
    tatoebaSentenceNumber1: '',
    tatoebaSentenceNumber2: '',
  });
  const [errorMessage, setErrorMessage] = useState('');
  console.log('values.text', values.text);
  useEffect(() => {
    setValues({
      ...values,
      text: problem.text,
      pos: String(problem.pos),
      lang2: problem.lang2,
      translated: problem.translated,
      exampleSentenceText: problem.sentence1.text,
      exampleSentenceTranslated: problem.sentence1.translated,
      exampleSentenceNote: problem.sentence1.note,
    });
  }, [problem.id, problem.version]);

  useEffect(() => {
    if (values.text.length === 0) {
      return;
    }
    if (tatoebaSentenceFindLoading) {
      return;
    }
    const f = async () => {
      await dispatch(
        findTatoebaSentences({
          param: {
            pageNo: 1,
            pageSize: 10,
            keyword: values.text,
            random: true,
          },
          postSuccessProcess: emptyFunction,
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  }, [dispatch, values.text, tatoebaSentenceFindLoading]);

  const EnglishWordProblemEditFormikForm = editFormikForm(
    props.workbook.id,
    problem.id,
    problem.version,
    setValues,
    setErrorMessage
  );

  console.log('note', values.exampleSentenceNote);

  return (
    <Container fluid>
      <PrivateProblemBreadcrumb
        name={props.workbook.name}
        id={props.workbook.id}
        text={props.problem.number.toString()}
      />
      <Divider hidden />
      <EnglishWordProblemEditFormikForm
        text={values.text}
        pos={String(values.pos)}
        lang2={values.lang2}
        translated={values.translated}
        exampleSentenceText={values.exampleSentenceText}
        exampleSentenceTranslated={values.exampleSentenceTranslated}
        exampleSentenceNote={values.exampleSentenceNote}
        sentenceProvider={values.sentenceProvider}
        tatoebaSentenceNumber1={values.tatoebaSentenceNumber1}
        tatoebaSentenceNumber2={values.tatoebaSentenceNumber2}
        tatoebaSentences={tatoebaSentences}
      />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};
