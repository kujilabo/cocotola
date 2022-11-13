import React from 'react';
import { Button, Card, Header } from 'semantic-ui-react';
import { FormikProps } from 'formik';
import { Form, Input, Select } from 'formik-semantic-ui-react';
import { posOptions } from 'plugins/translation/components';
import { langOptions } from 'components';
import { AppDimmer } from 'components';
import 'App.css';

export interface EnglishWordProblemNewFormValues {
  text: string;
  pos: string;
  lang2: string;
  loading: boolean;
}
export const EnglishWordProblemNewForm = (
  props: FormikProps<EnglishWordProblemNewFormValues>
): React.ReactElement => {
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
