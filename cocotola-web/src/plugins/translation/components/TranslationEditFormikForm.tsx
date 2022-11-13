import React from 'react';

import { withFormik, FormikBag } from 'formik';
import * as Yup from 'yup';

import { useAppDispatch } from '@/app/hooks';

import { removeTranslation } from '../features/translation_remove';
import { updateTranslation } from '../features/translation_update';

import {
  TranslationEditForm,
  TranslationEditFormValues,
} from './TranslationEditForm';

export interface TranslationEditFormikFormProps {
  index: number;
  selectedLang2: string;
  text: string;
  pos: string;
  translated: string;
  lang2: string;
  provider: string;
  //   onRemoveClick: () => void;
  refreshTranslations: () => void;
}
export const translationEditFormikForm = (
  setSuccessMessage: React.Dispatch<React.SetStateAction<string>>,
  setErrorMessage: React.Dispatch<React.SetStateAction<string>>,
  setTranslation: (t: TranslationEditFormValues) => void
): React.ComponentType<TranslationEditFormikFormProps> => {
  const dispatch = useAppDispatch();

  return withFormik<TranslationEditFormikFormProps, TranslationEditFormValues>({
    mapPropsToValues: (props: TranslationEditFormikFormProps) => ({
      // ...props,
      //   index: props.index,
      //   selectedLang: props.selectedLang,
      lang2: props.lang2,
      text: props.text,
      pos: props.pos,
      translated: props.translated,
      provider: props.provider,
      //   onRemoveClick: props.onRemoveClick,
      onRemoveClick: () => {
        dispatch(
          removeTranslation({
            param: {
              text: props.text,
              pos: +props.pos,
            },
            postSuccessProcess: () => props.refreshTranslations(),
            postFailureProcess: setErrorMessage,
          })
        );
      },
    }),
    validationSchema: Yup.object().shape({
      text: Yup.string().required('Word is required'),
    }),
    handleSubmit: (
      formValues: TranslationEditFormValues,
      formikBag: FormikBag<
        TranslationEditFormikFormProps,
        TranslationEditFormValues
      >
    ) => {
      // onsole.log('handleSubmit');
      setTranslation(formValues);

      setErrorMessage('');
      setSuccessMessage('');
      dispatch(
        updateTranslation({
          param: {
            lang2: formikBag.props.lang2,
            text: formValues.text,
            pos: +formValues.pos,
            translated: formValues.translated,
          },
          postSuccessProcess: () => {
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
  })(TranslationEditForm);
};
