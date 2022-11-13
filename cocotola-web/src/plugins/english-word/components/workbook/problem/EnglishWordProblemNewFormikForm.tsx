import React from 'react';
import { useHistory } from 'react-router-dom';
import { withFormik, FormikBag } from 'formik';
import * as Yup from 'yup';

import { useAppDispatch } from 'app/hooks';
import { addProblem } from 'features/problem_add';
import { EnglishWordProblemTypeId } from 'models/problem';
import {
  EnglishWordProblemNewForm,
  EnglishWordProblemNewFormValues,
} from './EnglishWordProblemNewForm';
import 'App.css';

export interface EnglishWordProblemNewFormikFormProps {
  text: string;
  pos: string;
  lang2: string;
  loading: boolean;
}
export const englishWordProblemNewFormikForm = (
  workbookId: number,
  setErrorMessage: React.Dispatch<React.SetStateAction<string>>,
  setProblem: (t: EnglishWordProblemNewFormValues) => void
): React.ComponentType<EnglishWordProblemNewFormikFormProps> => {
  const history = useHistory();
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
      dispatch(
        addProblem({
          workbookId: workbookId,
          param: {
            number: 1,
            problemType: EnglishWordProblemTypeId,
            properties: {
              text: values.text,
              pos: values.pos,
              lang2: values.lang2,
            },
          },
          postSuccessProcess: () =>
            history.push(`/app/private/workbook/${workbookId}`),
          postFailureProcess: setErrorMessage,
        })
      );
      setProblem(values);
    },
  })(EnglishWordProblemNewForm);
};
