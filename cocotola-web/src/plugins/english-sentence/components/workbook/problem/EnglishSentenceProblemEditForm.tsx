import { ReactElement } from 'react';

import { FormikProps } from 'formik';
import { Form, Input, Select } from 'formik-semantic-ui-react';
import { Button, Card, Header } from 'semantic-ui-react';

import { langOptions } from '@/components';

export interface EnglishSentenceProblemEditFormValues {
  number: number;
  text: string;
  lang2: string;
  translated: string;
  // note: string;
}

export const EnglishSentenceProblemEditForm = (
  props: FormikProps<EnglishSentenceProblemEditFormValues>
): ReactElement => {
  const { values, isSubmitting } = props;

  return (
    <Form>
      <Card fluid>
        <Card.Content>
          <Header component="h2">Edit problem</Header>
        </Card.Content>
        <Card.Content>
          <Input
            name="text"
            label="Sentence"
            placeholder="english word"
            errorPrompt
          />
          <Select
            name="lang2"
            label="Lang"
            options={langOptions}
            value={values.lang2}
            errorPrompt
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
