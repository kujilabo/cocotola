import React from 'react';

import { withFormik, FormikBag } from 'formik';
import { useNavigate } from 'react-router-dom';
import * as Yup from 'yup';

import { useAppDispatch } from '@/app/hooks';
import { updateProblem } from '@/features/problem_update';
import { EnglishWordProblemTypeId } from '@/models/problem';
import { TatoebaSentencePairModel } from '@/plugins/tatoeba/models/tatoeba';

import { EnglishWordProblemModel } from '../../../models/english-word-problem';

import {
  EnglishWordProblemEditForm,
  EnglishWordProblemEditFormValues,
} from './EnglishWordProblemEditForm';

export interface EnglishWordProblemEditFormikFormProps {
  number: number;
  text: string;
  pos: string;
  lang2: string;
  translated: string;
  exampleSentenceText: string;
  exampleSentenceTranslated: string;
  exampleSentenceNote: string;
  sentenceProvider: string;
  tatoebaSentenceNumber1: string;
  tatoebaSentenceNumber2: string;
  tatoebaSentences: TatoebaSentencePairModel[];
}
export const EnglishWordProblemEditFormikForm = (
  workbookId: number,
  problem: EnglishWordProblemModel,
  setErrorMessage: React.Dispatch<React.SetStateAction<string>>,
  setProblem: (t: EnglishWordProblemEditFormValues) => void,
  selectSentence: (index: number, checked: boolean) => void
): React.ComponentType<EnglishWordProblemEditFormikFormProps> => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  return withFormik<
    EnglishWordProblemEditFormikFormProps,
    EnglishWordProblemEditFormValues
  >({
    mapPropsToValues: (props: EnglishWordProblemEditFormikFormProps) => ({
      number: props.number,
      text: props.text,
      pos: props.pos,
      lang2: props.lang2,
      translated: props.translated,
      exampleSentenceText: props.exampleSentenceText,
      exampleSentenceTranslated: props.exampleSentenceTranslated,
      exampleSentenceNote: props.exampleSentenceNote,
      sentenceProvider: props.sentenceProvider,
      tatoebaSentenceNumber1: props.tatoebaSentenceNumber1,
      tatoebaSentenceNumber2: props.tatoebaSentenceNumber2,
      tatoebaSentences: props.tatoebaSentences,
      selectSentence: selectSentence,
    }),
    validationSchema: Yup.object().shape({
      text: Yup.string().required('Word is required'),
    }),
    handleSubmit: (
      values: EnglishWordProblemEditFormValues,
      formikBag: FormikBag<
        EnglishWordProblemEditFormikFormProps,
        EnglishWordProblemEditFormValues
      >
    ) => {
      // onsole.log('handleSubmit');
      const f = async () => {
        await dispatch(
          updateProblem({
            param: {
              workbookId: workbookId,
              problemId: problem.id,
              version: problem.version,
              problemType: EnglishWordProblemTypeId,
              properties: {
                text: values.text,
                pos: values.pos,
                lang2: values.lang2,
                sentenceProvider: values.sentenceProvider,
                tatoebaSentenceNumber1: values.tatoebaSentenceNumber1,
                tatoebaSentenceNumber2: values.tatoebaSentenceNumber2,
              },
            },
            postSuccessProcess: () =>
              navigate(`/app/private/workbook/${workbookId}`),
            postFailureProcess: (error: string) => setErrorMessage(error),
          })
        );
      };
      f().catch(console.error);
      setProblem(values);
    },
  })(EnglishWordProblemEditForm);
};
