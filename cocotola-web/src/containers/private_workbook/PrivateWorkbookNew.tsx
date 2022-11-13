import { ReactElement, useState } from 'react';

import { Container, Divider } from 'semantic-ui-react';

import { useAppSelector } from '@/app/hooks';
import { AppBreadcrumb, ErrorMessage } from '@/components';
import { PrivateWorkbookNewFormikForm } from '@/components/workbook/PrivateWorkbookNewFormikForm';
import { selectWorkbookAddLoading } from '@/features/workbook_add';

export const PrivateWorkbookNew = (): ReactElement => {
  const workbookAddLoading = useAppSelector(selectWorkbookAddLoading);
  const [values, setValues] = useState({
    name: '',
    lang2: 'ja',
    questionText: '',
    problemType: '',
  });
  const [errorMessage, setErrorMessage] = useState('');

  const NewFormikForm = PrivateWorkbookNewFormikForm(
    setErrorMessage,
    setValues
  );

  return (
    <Container fluid>
      <AppBreadcrumb
        links={[{ text: 'workbook', url: '/app/private/workbook' }]}
        text={'New workbook'}
      />
      <Divider hidden />
      <NewFormikForm
        loading={workbookAddLoading}
        name={values.name}
        lang2={values.lang2}
        questionText={values.questionText}
        problemType={values.problemType}
      />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};
