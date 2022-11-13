import { Dispatch, SetStateAction, ComponentType } from 'react';

import { withFormik, FormikBag } from 'formik';
import { useNavigate } from 'react-router-dom';
import * as Yup from 'yup';

import { useAppDispatch } from '@/app/hooks';
import { updateProblem } from '@/features/problem_update';
import { EnglishSentenceProblemTypeId } from '@/models/problem';

import { EnglishSentenceProblemModel } from '../../../models/english-sentence-problem';

import {
  EnglishSentenceProblemEditForm,
  EnglishSentenceProblemEditFormValues,
} from './EnglishSentenceProblemEditForm';

export interface EnglishSentenceProblemEditFormikFormProps {
  number: number;
  text: string;
  lang2: string;
  translated: string;
  // note: string;
}
export const englishSentenceProblemEditFormikForm = (
  workbookId: number,
  problem: EnglishSentenceProblemModel,
  setErrorMessage: Dispatch<SetStateAction<string>>,
  setProblem: (t: EnglishSentenceProblemEditFormValues) => void
): ComponentType<EnglishSentenceProblemEditFormikFormProps> => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  return withFormik<
    EnglishSentenceProblemEditFormikFormProps,
    EnglishSentenceProblemEditFormValues
  >({
    // mapPropsToValues: (props: EnglishSentenceProblemEditFormikFormProps) => ({
    //   number: props.number,
    //   text: props.text,
    //   lang2: props.lang2,
    //   translated: props.translated,
    //   // note: props.note,
    // }),
    mapPropsToValues: (props: EnglishSentenceProblemEditFormikFormProps) => ({
      ...props,
    }),
    validationSchema: Yup.object().shape({
      text: Yup.string().required('Sentence is required'),
    }),
    handleSubmit: (
      values: EnglishSentenceProblemEditFormValues,
      formikBag: FormikBag<
        EnglishSentenceProblemEditFormikFormProps,
        EnglishSentenceProblemEditFormValues
      >
    ) => {
      // onsole.log('handleSubmit');
      dispatch(
        updateProblem({
          param: {
            workbookId: workbookId,
            problemId: problem.id,
            version: problem.version,
            number: 1,
            problemType: EnglishSentenceProblemTypeId,
            properties: {
              text: values.text,
              lang2: values.lang2,
              translated: values.translated,
              // note: values.note,
            },
          },
          postSuccessProcess: () =>
            navigate(`/app/private/workbook/${workbookId}`),
          postFailureProcess: (error: string) => setErrorMessage(error),
        })
      );
      setProblem(values);
    },
  })(EnglishSentenceProblemEditForm);
};
