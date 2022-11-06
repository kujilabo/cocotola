import { Dispatch, ComponentType, SetStateAction } from 'react';

import {
  EnglishSentenceProblemNewForm,
  EnglishSentenceProblemNewFormValues,
} from './EnglishSentenceProblemNewForm';

import { withFormik, FormikBag } from 'formik';
import { useNavigate } from 'react-router-dom';
import * as Yup from 'yup';

import { useAppDispatch } from '@/app/hooks';
import { addProblem } from '@/features/problem_add';
import { EnglishSentenceProblemTypeId } from '@/models/problem';

export interface EnglishSentenceProblemNewFormikFormProps {
  text: string;
  lang2: string;
  translated: string;
  loading: boolean;
}
export const englishSentenceProblemNewFormikForm = (
  workbookId: number,
  setErrorMessage: Dispatch<SetStateAction<string>>,
  setProblem: (t: EnglishSentenceProblemNewFormValues) => void
): ComponentType<EnglishSentenceProblemNewFormikFormProps> => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  return withFormik<
    EnglishSentenceProblemNewFormikFormProps,
    EnglishSentenceProblemNewFormValues
  >({
    mapPropsToValues: (props: EnglishSentenceProblemNewFormikFormProps) => ({
      text: props.text,
      lang2: props.lang2,
      translated: props.translated,
      loading: props.loading,
    }),
    validationSchema: Yup.object().shape({
      text: Yup.string().required('Sentence is required'),
    }),
    handleSubmit: (
      values: EnglishSentenceProblemNewFormValues,
      formikBag: FormikBag<
        EnglishSentenceProblemNewFormikFormProps,
        EnglishSentenceProblemNewFormValues
      >
    ) => {
      // onsole.log('handleSubmit');
      dispatch(
        addProblem({
          workbookId: workbookId,
          param: {
            number: 1,
            problemType: EnglishSentenceProblemTypeId,
            properties: {
              text: values.text,
              translated: values.translated,
              lang2: values.lang2,
            },
          },
          postSuccessProcess: () =>
            navigate(`/app/private/workbook/${workbookId}`),
          postFailureProcess: setErrorMessage,
        })
      );
      setProblem(values);
    },
  })(EnglishSentenceProblemNewForm);
};
