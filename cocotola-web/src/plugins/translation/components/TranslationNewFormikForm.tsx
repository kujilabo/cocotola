import React from 'react';

import { withFormik, FormikBag } from 'formik';
import * as Yup from 'yup';

import { useAppDispatch } from '@/app/hooks';

import { addTranslation } from '../features/translation_add';

import {
  TranslationNewForm,
  TranslationNewFormValues,
} from './TranslationNewForm';

export interface TranslationNewFormikFormProps {
  text: string;
  pos: string;
  translated: string;
  refreshTranslations: () => void;
}
export const translationNewFormikForm = (
  setSuccessMessage: React.Dispatch<React.SetStateAction<string>>,
  setErrorMessage: React.Dispatch<React.SetStateAction<string>>,
  setValues: (v: TranslationNewFormValues) => void
): React.ComponentType<TranslationNewFormikFormProps> => {
  const dispatch = useAppDispatch();
  return withFormik<TranslationNewFormikFormProps, TranslationNewFormValues>({
    mapPropsToValues: (props: TranslationNewFormikFormProps) => ({
      text: props.text,
      pos: props.pos,
      translated: props.translated,
    }),
    validationSchema: Yup.object().shape({
      text: Yup.string().required('Word is required'),
      translated: Yup.string().required('Translated is required'),
    }),
    handleSubmit: (
      formValues: TranslationNewFormValues,
      formikBag: FormikBag<
        TranslationNewFormikFormProps,
        TranslationNewFormValues
      >
    ) => {
      // onsole.log('handleSubmit');
      setValues(formValues);
      setErrorMessage('');
      setSuccessMessage('');
      dispatch(
        addTranslation({
          param: {
            text: formValues.text,
            pos: +formValues.pos,
            translated: formValues.translated,
            lang2: 'ja',
          },
          postSuccessProcess: () => {
            formikBag.props.refreshTranslations();
            setErrorMessage('');
            setSuccessMessage('Word has been updated successfully');
          },
          postFailureProcess: (err: string) => {
            setErrorMessage(err);
            setSuccessMessage('');
          },
        })
      );
    },
  })(TranslationNewForm);
};
