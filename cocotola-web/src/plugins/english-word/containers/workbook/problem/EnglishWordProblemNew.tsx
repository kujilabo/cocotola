import { FC, ReactElement, useState, Dispatch, SetStateAction } from 'react';

import { FormikProps } from 'formik';
import { Input, Select } from 'formik-semantic-ui-react';
import { Container, Divider } from 'semantic-ui-react';
import * as Yup from 'yup';

import { ErrorMessage, langOptions } from '@/components';
import { PrivateProblemBreadcrumb } from '@/components/PrivateProblemBreadcrumb';
import { ProblemNewFormikForm } from '@/components/problem/ProblemNewFormikForm';
import { EnglishWordProblemTypeId } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';
import { posOptions } from '@/plugins/translation/components';

interface formikFormProps {
  text: string;
  pos: string;
  translated: string;
  lang2: string;
}
interface formValues {
  text: string;
  pos: string;
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
    problemType: EnglishWordProblemTypeId,
    toContent: (props: FormikProps<formValues>) => {
      const { values } = props;
      return (
        <>
          <Input
            name="text"
            label="Word"
            placeholder="english word"
            errorPrompt
          />
          <Select
            name="pos"
            label="Pos"
            options={posOptions}
            value={values.pos}
            errorPrompt
          />
          <Select
            name="lang2"
            label="Lang"
            options={langOptions}
            value={values.lang2}
            errorPrompt
          />
        </>
      );
    },
    validationSchema: Yup.object().shape({
      text: Yup.string().required('Word is required'),
    }),
    propsToValues: (props: formikFormProps) => ({ ...props }),
    valuesToProperties: (values: formValues) => ({
      text: values.text,
      pos: values.pos,
      translated: values.translated,
      lang2: values.lang2,
    }),
    resetValues: (values: formValues) => setValues(values),
    setErrorMessage: setErrorMessage,
  });
};

type EnglishWordProblemNewProps = {
  workbook: WorkbookModel;
};

export const EnglishWordProblemNew: FC<EnglishWordProblemNewProps> = (
  props: EnglishWordProblemNewProps
): ReactElement => {
  const [values, setValues] = useState({
    text: 'pen',
    pos: '99',
    translated: '',
    lang2: 'ja',
  });
  const [errorMessage, setErrorMessage] = useState('');
  const EnglishWordProblemNewFormikForm = newFormikForm(
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
      <EnglishWordProblemNewFormikForm
        text={values.text}
        pos={values.pos}
        translated={values.translated}
        lang2={values.lang2}
      />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};
