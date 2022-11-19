import { FC, useEffect, useState, Dispatch, SetStateAction } from 'react';

import { FormikProps } from 'formik';
import { Input, Select } from 'formik-semantic-ui-react';
import { Container, Divider } from 'semantic-ui-react';
import * as Yup from 'yup';

import {
  ErrorMessage,
  PrivateProblemBreadcrumb,
  langOptions,
} from '@/components';
import { ProblemEditFormikForm } from '@/components/problem/ProblemEditFormikForm';
import { ProblemModel, EnglishSentenceProblemTypeId } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';
import { EnglishSentenceProblemModel } from '@/plugins/english-sentence/models/english-sentence-problem';

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

const editFormikForm = (
  workbookId: number,
  problemId: number,
  problemVersion: number,
  resetValues: (v: formValues) => void,
  setErrorMessage: Dispatch<SetStateAction<string>>
) => {
  return ProblemEditFormikForm({
    workbookId: workbookId,
    problemId: problemId,
    problemVersion: problemVersion,
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
    resetValues: (values: formValues) => resetValues(values),
    setErrorMessage: setErrorMessage,
  });
};

type EnglishSentenceProblemEditProps = {
  workbook: WorkbookModel;
  problem: ProblemModel;
};

export const EnglishSentenceProblemEdit: FC<EnglishSentenceProblemEditProps> = (
  props: EnglishSentenceProblemEditProps
) => {
  // const { _workbookId, _problemId } = useParams<ParamTypes>();
  // const workbookId = +(_workbookId || '');
  // const problemId = +(_problemId || '');
  // const problemMap = useAppSelector(selectProblemMap);
  // const problem = EnglishSentenceProblemModel.of(problemMap[problemId]);
  const problem = EnglishSentenceProblemModel.of(props.problem);
  const [values, setValues] = useState({
    // number: problem.number,
    text: problem.text,
    lang2: problem.lang2,
    translated: problem.translated,
  });
  const [errorMessage, setErrorMessage] = useState('');
  console.log('values.text', values.text);
  useEffect(() => {
    setValues({
      ...values,
      // number: problem.number,
      text: problem.text,
      lang2: problem.lang2,
      translated: problem.translated,
    });
  }, [problem.id, problem.version]); // eslint-disable-line react-hooks/exhaustive-deps

  const EnglishSentenceProblemEditFormikForm = editFormikForm(
    props.workbook.id,
    problem.id,
    problem.version,
    setValues,
    setErrorMessage
  );

  return (
    <Container fluid>
      <PrivateProblemBreadcrumb
        name={props.workbook.name}
        id={props.workbook.id}
        text={props.problem.number.toString()}
      />
      <Divider hidden />
      <EnglishSentenceProblemEditFormikForm
        text={values.text}
        lang2={values.lang2}
        translated={values.translated}
      />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};
