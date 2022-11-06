import { FC, ReactElement } from 'react';

import { Link } from 'react-router-dom';
import { Card, Grid } from 'semantic-ui-react';

import { WorkbookModel } from '@/models/workbook';

export const WorkbookCard: FC<WorkbookProps> = (
  props: WorkbookProps
): ReactElement => {
  return (
    <Card className="full-width">
      <Card.Content>
        <Card.Header>
          <Link to={`${props.basePath}${props.model.id}`}>
            {props.model.name}
          </Link>
        </Card.Header>
      </Card.Content>
      <Card.Content>
        <Grid>
          <Grid.Row>
            <Grid.Column mobile={16} tablet={8} computer={8}></Grid.Column>
            <Grid.Column mobile={16} tablet={8} computer={8}></Grid.Column>
          </Grid.Row>
        </Grid>
      </Card.Content>
    </Card>
  );
};

type WorkbookProps = {
  basePath: string;
  model: WorkbookModel;
};
