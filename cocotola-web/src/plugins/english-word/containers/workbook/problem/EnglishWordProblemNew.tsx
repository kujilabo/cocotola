import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { Container, Divider } from 'semantic-ui-react';

import { useAppSelector } from 'app/hooks';
import {
  selectWorkbook,
  selectWorkbookGetLoading,
} from 'features/workbook_get';
import { selectProblemAddLoading } from 'features/problem_add';

import { ErrorMessage } from 'components';
import { PrivateProblemBreadcrumb } from 'components/PrivateProblemBreadcrumb';
import { WorkbookModel } from 'models/workbook';

import { englishWordProblemNewFormikForm } from '../../../components/workbook/problem/EnglishWordProblemNewFormikForm';

import 'App.css';

type ParamTypes = {
  _workbookId: string;
};

export const EnglishWordProblemNew: React.FC<EnglishWordProblemNewProps> = (
  props: EnglishWordProblemNewProps
) => {
  const { _workbookId } = useParams<ParamTypes>();
  const workbookId = +_workbookId;
  const workbook = useAppSelector(selectWorkbook);
  const workbookGetLoading = useAppSelector(selectWorkbookGetLoading);
  const problemAddLoading = useAppSelector(selectProblemAddLoading);
  const [values, setValues] = useState({
    text: 'pen',
    pos: '99',
    lang2: 'ja',
  });
  const [errorMessage, setErrorMessage] = useState('');
  const loading = workbookGetLoading || problemAddLoading;

  const EnglishWordProblemNewFormikForm = englishWordProblemNewFormikForm(
    workbookId,
    setErrorMessage,
    setValues
  );

  return (
    <Container fluid>
      <PrivateProblemBreadcrumb
        name={workbook.name}
        id={workbookId}
        text={'New problem'}
      />
      <Divider hidden />
      <EnglishWordProblemNewFormikForm
        text={values.text}
        pos={values.pos}
        lang2={values.lang2}
        loading={loading}
      />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};

type EnglishWordProblemNewProps = {
  workbook: WorkbookModel;
};
