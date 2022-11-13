import { FC, useState, Dispatch, SetStateAction } from 'react';

import { Input, Select } from 'formik-semantic-ui-react';
import { useParams } from 'react-router-dom';
import { Container, Divider } from 'semantic-ui-react';
import * as Yup from 'yup';

import { useAppSelector } from '@/app/hooks';
import { ErrorMessage, langOptions, AppDimmer } from '@/components';
import { PrivateProblemBreadcrumb } from '@/components/PrivateProblemBreadcrumb';
import {
  FormValues,
  FormikFormProps,
  problemNewFormikForm,
} from '@/components/problem/ProblemNewFormikForm';
import {
  selectWorkbook,
  selectWorkbookGetLoading,
} from '@/features/workbook_get';
import { EnglishSentenceProblemTypeId } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';

interface formikFormProps extends FormikFormProps {
  text: string;
  translated: string;
  lang2: string;
}

interface formValues extends FormValues {
  text: string;
  translated: string;
  lang2: string;
}

const newFormikForm = (
  workbookId: number,
  setValues: (v: formValues) => void,
  setErrorMessage: Dispatch<SetStateAction<string>>
) => {
  return problemNewFormikForm({
    workbookId: workbookId,
    problemType: EnglishSentenceProblemTypeId,
    toContent: (values: formValues) => {
      return (
        <>
          <Input
            name="text"
            label="Sentence"
            placeholder="english sentence"
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
        </>
      );
    },
    validationSchema: Yup.object().shape({
      text: Yup.string().required('Sentence is required'),
    }),
    propsToValues: (props: formikFormProps) => ({ ...props }),
    valuesToProperties: (values: formValues) => ({
      text: values.text,
      translated: values.translated,
      lang2: values.lang2,
    }),
    resetValues: (values: formValues) => setValues(values),
    setErrorMessage: setErrorMessage,
  });
};

type ParamTypes = {
  _workbookId: string;
};

type EnglishSentenceProblemNewProps = {
  workbook: WorkbookModel;
};

export const EnglishSentenceProblemNew: FC<EnglishSentenceProblemNewProps> = (
  props: EnglishSentenceProblemNewProps
) => {
  const { _workbookId } = useParams<ParamTypes>();
  const workbookId = +(_workbookId || '');
  const workbook = useAppSelector(selectWorkbook);
  const workbookGetLoading = useAppSelector(selectWorkbookGetLoading);
  const [values, setValues] = useState({
    text: 'pen',
    lang2: 'ja',
    translated: '',
  });
  const [errorMessage, setErrorMessage] = useState('');
  const loading = workbookGetLoading;
  const EnglishSentenceProblemNewFormikForm = newFormikForm(
    workbookId,
    (values: formValues) => setValues(values),
    setErrorMessage
  );

  return (
    <Container fluid>
      <PrivateProblemBreadcrumb
        name={workbook.name}
        id={workbookId}
        text={'New problem'}
      />
      <Divider hidden />
      {loading ? <AppDimmer /> : <div />}
      <EnglishSentenceProblemNewFormikForm
        text={values.text}
        lang2={values.lang2}
        translated={values.translated}
      />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};
