import React from 'react';

import { FormikProps } from 'formik';
import { Form, Input, Select } from 'formik-semantic-ui-react';
import { Button, Card, Header } from 'semantic-ui-react';

import { langOptions } from '@/components';
import {
  TatoebaSentencePairModel,
  TatoebaSentenceModel,
} from '@/plugins/tatoeba/models/tatoeba';
import { posOptions } from '@/plugins/translation/components';

import { ExampleTatoebaSentence } from './ExampleTatoebaSentence';
import { ExampleTatoebaSentenceList } from './ExampleTatoebaSentenceList';

export interface EnglishWordProblemEditFormValues {
  number: number;
  text: string;
  pos: string;
  lang2: string;
  translated: string;
  exampleSentenceText: string;
  exampleSentenceTranslated: string;
  exampleSentenceNote: string;
  sentenceProvider: string;
  tatoebaSentenceNumber1: string;
  tatoebaSentenceNumber2: string;
  tatoebaSentences: TatoebaSentencePairModel[];
  selectSentence: (index: number, checked: boolean) => void;
}
const emptyTatoebaSentence: TatoebaSentenceModel = {
  text: '',
  author: '',
  sentenceNumber: 0,
  lang2: '',
};
export const EnglishWordProblemEditForm = (
  props: FormikProps<EnglishWordProblemEditFormValues>
): React.ReactElement => {
  const { values, isSubmitting, setFieldValue } = props;
  const onCheckboxChange = (
    sentenceNumber1: number,
    sentenceNumber2: number,
    checked: boolean
  ): void => {
    if (checked) {
      setFieldValue('sentenceProvider', 'tatoeba');
      setFieldValue('tatoebaSentenceNumber1', sentenceNumber1.toString());
      setFieldValue('tatoebaSentenceNumber2', sentenceNumber2.toString());
    } else {
      setFieldValue('sentenceProvider', '');
      setFieldValue('tatoebaSentenceNumber1', '');
      setFieldValue('tatoebaSentenceNumber2', '');
    }
  };

  let sentence1 = emptyTatoebaSentence;
  let sentence2 = emptyTatoebaSentence;
  if (values.exampleSentenceNote && values.exampleSentenceNote !== '') {
    try {
      /* eslint-disable */
      const noteObj: { [key: string]: string } = JSON.parse(
        values.exampleSentenceNote
      );
      /* eslint-enable */
      console.log('noteObj', noteObj);
      sentence1 = {
        text: values.exampleSentenceText,
        author: noteObj['tatoebaAuthor1'],
        sentenceNumber: +noteObj['tatoebaSentenceNumber1'],
        lang2: 'en',
      };
      sentence2 = {
        text: values.exampleSentenceTranslated,
        author: noteObj['tatoebaAuthor2'],
        sentenceNumber: +noteObj['tatoebaSentenceNumber2'],
        lang2: 'ja',
      };
    } catch (e) {
      console.log(e);
    }
  }

  return (
    <Form>
      <Card fluid>
        <Card.Content>
          <Header component="h2">Edit problem</Header>
        </Card.Content>
        <Card.Content>
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
          {sentence1.text !== '' ? (
            <ExampleTatoebaSentence
              sentence1={sentence1}
              sentence2={sentence2}
            />
          ) : (
            <div />
          )}
          <ExampleTatoebaSentenceList
            sentences={values.tatoebaSentences}
            onCheckboxChange={onCheckboxChange}
          />
        </Card.Content>
        <Card.Content>
          <Button
            type="submit"
            // variant="true"
            color="teal"
            disabled={isSubmitting}
          >
            Update
          </Button>
        </Card.Content>
      </Card>
    </Form>
  );
};
