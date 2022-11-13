import React from 'react';
import { Form, Checkbox, CheckboxProps, Table } from 'semantic-ui-react';

import { TatoebaSentencePairModel } from 'plugins/tatoeba/models/tatoeba';
import {
  makeAuthorLink,
  makeLicenseLink,
  makeSentenceLink,
} from 'plugins/tatoeba/components/util';

import 'App.css';

type ExampleTatoebaSentenceListProps = {
  sentences: TatoebaSentencePairModel[];
  onCheckboxChange: (
    sentenceNumber1: number,
    sentenceNumber2: number,
    checked: boolean
  ) => void;
};

export const ExampleTatoebaSentenceList: React.FC<
  ExampleTatoebaSentenceListProps
> = (props: ExampleTatoebaSentenceListProps) => {
  return (
    <Form.Group grouped>
      <label>Example Sentences</label>

      <Table compact celled definition>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell />
            <Table.HeaderCell>Sentence1</Table.HeaderCell>
            <Table.HeaderCell>Sentence2</Table.HeaderCell>
            <Table.HeaderCell>Tatoeba1</Table.HeaderCell>
            <Table.HeaderCell>Tatoeba2</Table.HeaderCell>
            <Table.HeaderCell>Author1</Table.HeaderCell>
            <Table.HeaderCell>Author2</Table.HeaderCell>
            <Table.HeaderCell>License</Table.HeaderCell>
          </Table.Row>
        </Table.Header>

        <Table.Body>
          {props.sentences.map((m: TatoebaSentencePairModel, i: number) => {
            return (
              <Table.Row key={i}>
                <Table.Cell collapsing>
                  <Checkbox
                    slider
                    onChange={(
                      event: React.FormEvent<HTMLInputElement>,
                      data: CheckboxProps
                    ) =>
                      props.onCheckboxChange(
                        m.src.sentenceNumber,
                        m.dst.sentenceNumber,
                        data.checked || false
                      )
                    }
                  />
                </Table.Cell>
                <Table.Cell>{m.src.text}</Table.Cell>
                <Table.Cell>{m.dst.text}</Table.Cell>
                <Table.Cell>
                  {makeSentenceLink(m.src.sentenceNumber)}
                </Table.Cell>
                <Table.Cell>
                  {makeSentenceLink(m.dst.sentenceNumber)}
                </Table.Cell>
                <Table.Cell>{makeAuthorLink(m.src.author)}</Table.Cell>
                <Table.Cell>{makeAuthorLink(m.dst.author)}</Table.Cell>
                <Table.Cell>{makeLicenseLink()}</Table.Cell>
              </Table.Row>
            );
          })}
        </Table.Body>
      </Table>
    </Form.Group>
  );
};
