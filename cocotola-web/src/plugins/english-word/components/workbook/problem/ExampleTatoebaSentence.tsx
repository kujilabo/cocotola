import React from 'react';
import { Form, Table } from 'semantic-ui-react';

import { TatoebaSentenceModel } from 'plugins/tatoeba/models/tatoeba';
import {
  makeAuthorLink,
  makeLicenseLink,
  makeSentenceLink,
} from 'plugins/tatoeba/components/util';

import 'App.css';

type ExampleTatoebaSentenceProps = {
  sentence1: TatoebaSentenceModel;
  sentence2: TatoebaSentenceModel;
};

export const ExampleTatoebaSentence: React.FC<ExampleTatoebaSentenceProps> = (
  props: ExampleTatoebaSentenceProps
) => {
  console.log('ExampleTatoebaSentence');
  return (
    <Form.Group grouped>
      <label>Example Sentence</label>
      <Table compact celled>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>Sentence</Table.HeaderCell>
            <Table.HeaderCell>Tatoeba</Table.HeaderCell>
            <Table.HeaderCell>Author</Table.HeaderCell>
            <Table.HeaderCell>License</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          <Table.Row>
            <Table.Cell>{props.sentence1.text}</Table.Cell>
            <Table.Cell>
              {makeSentenceLink(props.sentence1.sentenceNumber)}
            </Table.Cell>
            <Table.Cell>{makeAuthorLink(props.sentence1.author)}</Table.Cell>
            <Table.Cell>{makeLicenseLink()}</Table.Cell>
          </Table.Row>
          <Table.Row>
            <Table.Cell>{props.sentence2.text}</Table.Cell>
            <Table.Cell>
              {makeSentenceLink(props.sentence2.sentenceNumber)}
            </Table.Cell>
            <Table.Cell>{makeAuthorLink(props.sentence2.author)}</Table.Cell>
            <Table.Cell>{makeLicenseLink()}</Table.Cell>
          </Table.Row>
        </Table.Body>
      </Table>
    </Form.Group>
  );
};
