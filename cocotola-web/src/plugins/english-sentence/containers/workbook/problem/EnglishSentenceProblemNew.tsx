import { FC, useState, Dispatch, SetStateAction } from 'react';

import { FormikProps } from 'formik';
import { Input, Select } from 'formik-semantic-ui-react';
import { Container, Divider } from 'semantic-ui-react';
import * as Yup from 'yup';

import { ErrorMessage, langOptions } from '@/components';
import { PrivateProblemBreadcrumb } from '@/components/PrivateProblemBreadcrumb';
import { ProblemNewFormikForm } from '@/components/problem/ProblemNewFormikForm';
import { EnglishSentenceProblemTypeId } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';

interface formikFormProps {
  text: string;
  translated: string;
  lang2: string;
}

interface formValues {
  text: string;
  translated: string;
  lang2: string;
}

const newFormikForm = (
  workbookId: number,
  setValues: (v: formValues) => void,
  setErrorMessage: Dispatch<SetStateAction<string>>
) => {
  return ProblemNewFormikForm({
    workbookId: workbookId,
    problemType: EnglishSentenceProblemTypeId,
    toContent: (props: FormikProps<formValues>) => {
      const { values } = props;
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

type EnglishSentenceProblemNewProps = {
  workbook: WorkbookModel;
};

export const EnglishSentenceProblemNew: FC<EnglishSentenceProblemNewProps> = (
  props: EnglishSentenceProblemNewProps
) => {
  const [values, setValues] = useState({
    text: 'pen',
    lang2: 'ja',
    translated: '',
  });
  const [errorMessage, setErrorMessage] = useState('');
  const EnglishSentenceProblemNewFormikForm = newFormikForm(
    props.workbook.id,
    (values: formValues) => setValues(values),
    setErrorMessage
  );

  return (
    <Container fluid>
      <PrivateProblemBreadcrumb
        name={props.workbook.name}
        id={props.workbook.id}
        text={'New problem'}
      />
      <Divider hidden />
      <EnglishSentenceProblemNewFormikForm
        text={values.text}
        translated={values.translated}
        lang2={values.lang2}
      />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};
