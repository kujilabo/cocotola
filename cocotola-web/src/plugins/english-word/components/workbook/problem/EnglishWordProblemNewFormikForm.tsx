import React from 'react';

import { withFormik, FormikBag } from 'formik';
import { useNavigate } from 'react-router-dom';
import * as Yup from 'yup';

import { useAppDispatch } from '@/app/hooks';
import { addProblem } from '@/features/problem_add';
import { EnglishWordProblemTypeId } from '@/models/problem';

import {
  EnglishWordProblemNewForm,
  EnglishWordProblemNewFormValues,
} from './EnglishWordProblemNewForm';

export interface EnglishWordProblemNewFormikFormProps {
  text: string;
  pos: string;
  lang2: string;
  loading: boolean;
}
export const EnglishWordProblemNewFormikForm = (
  workbookId: number,
  setErrorMessage: React.Dispatch<React.SetStateAction<string>>,
  setProblem: (t: EnglishWordProblemNewFormValues) => void
): React.ComponentType<EnglishWordProblemNewFormikFormProps> => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  return withFormik<
    EnglishWordProblemNewFormikFormProps,
    EnglishWordProblemNewFormValues
  >({
    mapPropsToValues: (props: EnglishWordProblemNewFormikFormProps) => ({
      text: props.text,
      pos: props.pos,
      lang2: props.lang2,
      loading: props.loading,
    }),
    validationSchema: Yup.object().shape({
      text: Yup.string().required('Word is required'),
    }),
    handleSubmit: (
      values: EnglishWordProblemNewFormValues,
      formikBag: FormikBag<
        EnglishWordProblemNewFormikFormProps,
        EnglishWordProblemNewFormValues
      >
    ) => {
      // onsole.log('handleSubmit');
      const f = async () => {
        await dispatch(
          addProblem({
            workbookId: workbookId,
            param: {
              problemType: EnglishWordProblemTypeId,
              properties: {
                text: values.text,
                pos: values.pos,
                lang2: values.lang2,
              },
            },
            postSuccessProcess: () =>
              navigate(`/app/private/workbook/${workbookId}`),
            postFailureProcess: setErrorMessage,
          })
        );
      };
      f().catch(console.error);
      setProblem(values);
    },
  })(EnglishWordProblemNewForm);
};
