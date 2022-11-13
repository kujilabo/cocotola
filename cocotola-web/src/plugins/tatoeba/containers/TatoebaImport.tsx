import React, { ChangeEvent, useState } from 'react';

import { Card, Container, Divider, Header, Form } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppDimmer, AppBreadcrumb, ErrorMessage } from '@/components';
import { UploadButton } from '@/components/buttons';

import {
  importTatoebaSentence,
  importTatoebaLink,
  selectTatoebaImportLoading,
} from '../features/tatoeba_import';

export const TatoebaImport = (): React.ReactElement => {
  const tatoebaImportLoading = useAppSelector(selectTatoebaImportLoading);
  const loading = tatoebaImportLoading;
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
  const onUploadSentenceButtonClick = () => {
    const formData = new FormData();
    formData.append('file', file as Blob);

    dispatch(
      importTatoebaSentence({
        param: formData,
        postSuccessProcess: () => setErrorMessage(''),
        postFailureProcess: setErrorMessage,
      })
    );
  };
  const onUploadLinkButtonClick = () => {
    const formData = new FormData();
    formData.append('file', file as Blob);

    dispatch(
      importTatoebaLink({
        param: formData,
        postSuccessProcess: () => setErrorMessage(''),
        postFailureProcess: setErrorMessage,
      })
    );
  };

  return (
    <Container fluid>
      <AppBreadcrumb links={[]} text={'Import'} />
      <Divider hidden />
      <Card fluid>
        <Card.Content>
          <Header component="h2">Import sentences</Header>
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
          <UploadButton type="button" onClick={onUploadSentenceButtonClick} />
        </Card.Content>
      </Card>
      <Card fluid>
        <Card.Content>
          <Header component="h2">Import links</Header>
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
          <UploadButton type="button" onClick={onUploadLinkButtonClick} />
        </Card.Content>
      </Card>
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};
