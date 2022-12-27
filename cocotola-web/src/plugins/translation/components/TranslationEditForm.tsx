import React from 'react';

import { FormikProps } from 'formik';
import { Form } from 'formik-semantic-ui-react';
import { Button, Card, Header } from 'semantic-ui-react';

import { DeleteButton, UpdateButton } from '@/components/buttons';

import { InputWord, InputTrasnslatedWord, SelectPos } from '../components';

export interface TranslationEditFormValues {
  // index: number;
  // selectedLang: string;
  lang2: string;
  text: string;
  pos: string;
  translated: string;
  provider: string;
  onRemoveClick: () => void;
}
export const TranslationEditForm = (
  props: FormikProps<TranslationEditFormValues>
) => {
  const { isSubmitting } = props;

  return (
    <Form>
      <Card fluid>
        <Card.Content>
          <Header component="h2">Edit Translation</Header>
        </Card.Content>
        <Card.Content>
          <InputWord disabled />
          <SelectPos disabled />
          <InputTrasnslatedWord />
        </Card.Content>
        <Card.Content>
          <div className="ui fluid buttons">
            <UpdateButton type="submit" disabled={isSubmitting} />
          </div>

          {props.values.provider === 'custom' ? (
            <div className="ui fluid buttons">
              <DeleteButton
                type="button"
                disabled={isSubmitting}
                onClick={props.values.onRemoveClick}
              />
            </div>
          ) : (
            <div />
          )}
        </Card.Content>
      </Card>
    </Form>
  );
};
