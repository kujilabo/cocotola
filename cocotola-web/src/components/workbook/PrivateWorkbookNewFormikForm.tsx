import { ComponentType, SetStateAction, Dispatch } from 'react';

import {
  PrivateWorkbookNewForm,
  PrivateWorkbookNewFormValues,
} from './PrivateWorkbookNewForm';

import { withFormik, FormikBag } from 'formik';
import { useNavigate } from 'react-router-dom';
import * as Yup from 'yup';

import { useAppDispatch } from '@/app/hooks';
import { addWorkbook } from '@/features/workbook_add';

export interface PrivateWorkbookNewFormikFormProps {
  name: string;
  lang2: string;
  questionText: string;
  problemType: string;
  loading: boolean;
}
export const privateWorkbookNewFormikForm = (
  setErrorMessage: Dispatch<SetStateAction<string>>,
  setWorkbook: (t: PrivateWorkbookNewFormValues) => void
): ComponentType<PrivateWorkbookNewFormikFormProps> => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  return withFormik<
    PrivateWorkbookNewFormikFormProps,
    PrivateWorkbookNewFormValues
  >({
    mapPropsToValues: (props: PrivateWorkbookNewFormikFormProps) => ({
      name: props.name,
      lang2: props.lang2,
      questionText: props.questionText,
      problemType: props.problemType,
      loading: props.loading,
    }),
    validationSchema: Yup.object().shape({
      name: Yup.string().required('Name is required'),
      problemType: Yup.string().required('Problem type is required'),
    }),
    handleSubmit: (
      values: PrivateWorkbookNewFormValues,
      formikBag: FormikBag<
        PrivateWorkbookNewFormikFormProps,
        PrivateWorkbookNewFormValues
      >
    ) => {
      // onsole.log('handleSubmit');
      dispatch(
        addWorkbook({
          param: { ...values, spaceKey: 'personal' },
          postSuccessProcess: (id: number) => navigate('/app/private/workbook'),
          postFailureProcess: setErrorMessage,
        })
      );
      setWorkbook(values);
    },
  })(PrivateWorkbookNewForm);
};
