import { ReactElement } from 'react';

import { FormikProps } from 'formik';
import { Form, Input, Select } from 'formik-semantic-ui-react';
import { Button, Card, Header } from 'semantic-ui-react';

import { langOptions, AppDimmer } from '@/components';

export interface EnglishSentenceProblemNewFormValues {
  text: string;
  lang2: string;
  translated: string;
  loading: boolean;
}
export const EnglishSentenceProblemNewForm = (
  props: FormikProps<EnglishSentenceProblemNewFormValues>
): ReactElement => {
  const { values, isSubmitting } = props;

  return (
    <Form>
      <Card>
        <Card.Content>
          <Header component="h2">New problem</Header>
        </Card.Content>
        <Card.Content>
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
        </Card.Content>
        <Card.Content>
          {values.loading ? <AppDimmer /> : <div />}
          <Button
            type="submit"
            // variant="true"
            color="teal"
            disabled={isSubmitting}
          >
            Create
          </Button>
        </Card.Content>
      </Card>
    </Form>
  );
};
