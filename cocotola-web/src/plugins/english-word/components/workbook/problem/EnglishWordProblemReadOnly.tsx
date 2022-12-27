import React from 'react';

import { Button, Card, Grid, Header } from 'semantic-ui-react';

import { ProblemModel } from '@/models/problem';
import { EnglishWordProblemModel } from '@/plugins/english-word/models/english-word-problem';
import { toDsiplayText } from '@/plugins/english-word/utils/util';

type EnglishWordProblemReadOnlyProps = {
  workbookId: number;
  problem: ProblemModel;
};

export const EnglishWordProblemReadOnly: React.FC<
  EnglishWordProblemReadOnlyProps
> = (props: EnglishWordProblemReadOnlyProps) => {
  const problem = EnglishWordProblemModel.of(props.problem);

  return (
    <Card fluid>
      <Card.Content>
        <Card.Header>{problem.text}</Card.Header>
      </Card.Content>
      <Card.Content>
        <Grid columns={2}>
          <Grid.Row>
            <Grid.Column>
              <Header component="h2" className="border-bottom g-mb-15">
                {toDsiplayText(problem.pos)}
              </Header>
            </Grid.Column>
            <Grid.Column>
              <Header component="h2" className="border-bottom g-mb-15">
                {problem.translated}
              </Header>
            </Grid.Column>
          </Grid.Row>
        </Grid>
        <Grid>
          <Grid.Row>
            <Grid.Column>
              <Header component="h2" className="border-bottom g-mb-15">
                {/* {problem.phonetic} */}
              </Header>
            </Grid.Column>
            <Grid.Column></Grid.Column>
          </Grid.Row>
        </Grid>
      </Card.Content>
      <Card.Content extra>
        <div className="ui fluid buttons">
          <Button
            basic
            color="teal"
            // onClick={() =>
            //   props.getAudio(
            //     props.problem.properties['audioId'],
            //     props.problem.updatedAt,
            //     playAudio
            //   )
            // }
          >
            Play
          </Button>
        </div>
      </Card.Content>
    </Card>
  );
};
