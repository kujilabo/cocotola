import { FC, useEffect, useState } from 'react';

import { useParams } from 'react-router-dom';
import { Container, Divider } from 'semantic-ui-react';

import { useAppDispatch, useAppSelector } from '@/app/hooks';
import { ErrorMessage, PrivateProblemBreadcrumb } from '@/components';
// import { selectProblem } from '@/features/problem_get';
import { selectProblemMap } from '@/features/problem_find';
import { ProblemModel } from '@/models/problem';
import { WorkbookModel } from '@/models/workbook';

import { englishSentenceProblemEditFormikForm } from '../../../components/workbook/problem/EnglishSentenceProblemEditFormikForm';
import { EnglishSentenceProblemModel } from '../../../models/english-sentence-problem';

type ParamTypes = {
  _workbookId: string;
  _problemId: string;
};

export const EnglishSentenceProblemEdit: FC<EnglishSentenceProblemEditProps> = (
  props: EnglishSentenceProblemEditProps
) => {
  const { _workbookId, _problemId } = useParams<ParamTypes>();
  const workbookId = +(_workbookId || '');
  const problemId = +(_problemId || '');
  const dispatch = useAppDispatch();
  // const problem = EnglishSentenceProblemModel.of(useAppSelector(selectProblem));
  const problemMap = useAppSelector(selectProblemMap);
  const problem = EnglishSentenceProblemModel.of(problemMap[problemId]);
  const [values, setValues] = useState({
    number: problem.number,
    text: problem.text,
    lang2: problem.lang2,
    translated: problem.translated,
  });
  const [errorMessage, setErrorMessage] = useState('');
  console.log('values.text', values.text);
  useEffect(() => {
    setValues({
      ...values,
      number: problem.number,
      text: problem.text,
      lang2: problem.lang2,
      translated: problem.translated,
    });
  }, [dispatch, problem.id, problem.version]);

  const EnglishSentenceProblemEditFormikForm =
    englishSentenceProblemEditFormikForm(
      workbookId,
      problem,
      setErrorMessage,
      setValues
    );

  return (
    <Container fluid>
      <PrivateProblemBreadcrumb
        name={props.workbook.name}
        id={+(_workbookId || '')}
        text={'' + props.problem.number}
      />
      <Divider hidden />
      <EnglishSentenceProblemEditFormikForm
        number={values.number}
        text={values.text}
        lang2={values.lang2}
        translated={values.translated}
      />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};

type EnglishSentenceProblemEditProps = {
  workbook: WorkbookModel;
  problem: ProblemModel;
};
