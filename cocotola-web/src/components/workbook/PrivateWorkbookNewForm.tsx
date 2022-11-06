import { ReactElement } from 'react';

import { FormikProps } from 'formik';
import { Form, Input, Select } from 'formik-semantic-ui-react';
import { Button, Card, Header } from 'semantic-ui-react';

import { AppDimmer } from '@/components';

const problemTypeOptions = [
  { key: 'word', text: '単語', value: 'english_word' },
  { key: 'phrase', text: '文章', value: 'english_sentence' },
  { key: 'phrase', text: 'フレーズ', value: 'english_phrase' },
  { key: 'template', text: 'テンプレート', value: 'template' },
];

export interface PrivateWorkbookNewFormValues {
  name: string;
  lang2: string;
  questionText: string;
  problemType: string;
  loading: boolean;
}
export const PrivateWorkbookNewForm = (
  props: FormikProps<PrivateWorkbookNewFormValues>
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
            name="name"
            label="Name"
            placeholder="Workbook name"
            className="full-width"
            errorPrompt
          />
          <Input
            name="questionText"
            label="Question text"
            placeholder=""
            className="full-width"
            errorPrompt
          />
          <Select
            name="problemType"
            label="Problem type"
            options={problemTypeOptions}
            value={values.problemType}
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
