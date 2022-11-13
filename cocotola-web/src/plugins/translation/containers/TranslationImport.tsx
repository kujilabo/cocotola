import { ReactElement, ChangeEvent, useState } from 'react';

import { Card, Container, Divider, Header, Form } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppDimmer, AppBreadcrumb, ErrorMessage } from '@/components';
import { UploadButton } from '@/components/buttons';

import {
  importTranslation,
  selectTranslationImportLoading,
} from '../features/translation_import';

export const TranslationImport = (): ReactElement => {
  const translationImportLoading = useAppSelector(
    selectTranslationImportLoading
  );
  const loading = translationImportLoading;
  const dispatch = useAppDispatch();
  const [file, setFile] = useState({});
  // const [fileName, setFileName] = useState({});
  const [errorMessage, setErrorMessage] = useState('');
  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e == null || e.target == null || e.target.files == null) {
      return;
    }
    const fileList: FileList = e.target.files;
    if (fileList.item(0) == null) {
      return;
    }
    const file: File | null = fileList.item(0);
    if (file == null) {
      return;
    }
    setFile(file);
  };
  const onUploadButtonClick = () => {
    const formData = new FormData();
    formData.append('file', file as Blob);

    const f = async () => {
      await dispatch(
        importTranslation({
          param: formData,
          postSuccessProcess: () => setErrorMessage(''),
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  };

  return (
    <Container fluid>
      <AppBreadcrumb
        links={[{ text: 'Translations', url: '/plugin/translation/list' }]}
        text={'Import'}
      />
      <Divider hidden />
      <Card fluid>
        <Card.Content>
          <Header component="h2">Import translations</Header>
        </Card.Content>
        <Card.Content>
          <p>
            TEXT,POS,TRANSLATED
            <br />
            POS: adj,adv,conj,det,modal,noun,prep,pron,verb
          </p>
          <Form>
            <input type="file" name="text" onChange={handleFileChange} />
          </Form>
        </Card.Content>
        <Card.Content extra>
          {loading ? <AppDimmer /> : <div />}
          <UploadButton type="button" onClick={onUploadButtonClick} />
        </Card.Content>
      </Card>
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};
