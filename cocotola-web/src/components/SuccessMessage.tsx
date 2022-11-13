import { FC } from 'react';

import { Message } from 'semantic-ui-react';

type SuccessMessageProps = {
  message: string;
};

export const SuccessMessage: FC<SuccessMessageProps> = (
  props: SuccessMessageProps
) => {
  if (props.message) {
    return (
      <Message positive>
        <Message.Header>{props.message}</Message.Header>
      </Message>
    );
  }
  return <div />;
};
