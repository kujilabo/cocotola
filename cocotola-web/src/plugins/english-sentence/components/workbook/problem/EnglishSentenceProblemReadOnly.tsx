import { FC, ReactElement } from 'react';

import { Button, Card, Grid, Header } from 'semantic-ui-react';

import { ProblemModel } from '@/models/problem';
import { EnglishSentenceProblemModel } from '@/plugins/english-sentence/models/english-sentence-problem';

type EnglishSentenceProblemReadOnlyProps = {
  workbookId: number;
  problem: ProblemModel;
};

export const EnglishSentenceProblemReadOnly: FC<
  EnglishSentenceProblemReadOnlyProps
> = (props: EnglishSentenceProblemReadOnlyProps): ReactElement => {
  const problem = EnglishSentenceProblemModel.of(props.problem);

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
        <Button.Group floated="left">
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
        </Button.Group>
      </Card.Content>
    </Card>
  );
};
