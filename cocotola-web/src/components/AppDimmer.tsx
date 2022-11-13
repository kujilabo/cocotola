import { ReactElement } from 'react';

import { Dimmer, Loader } from 'semantic-ui-react';

export const AppDimmer = (): ReactElement => {
  return (
    <Dimmer active>
      <Loader>Loading</Loader>
    </Dimmer>
  );
};
