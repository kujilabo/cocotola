import { FC, ReactElement } from 'react';

import { Message } from 'semantic-ui-react';

type ErrorMessageProps = {
  message: string;
};

export const ErrorMessage: FC<ErrorMessageProps> = (
  props: ErrorMessageProps
): ReactElement => {
  if (props.message) {
    return (
      <Message negative>
        <Message.Header>{props.message}</Message.Header>
      </Message>
    );
  }
  return <div />;
};
