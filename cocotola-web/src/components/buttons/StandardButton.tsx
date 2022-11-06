import { FC } from 'react';

import { Button } from 'semantic-ui-react';

type StandardButtonProps = {
  value: string;
  type: 'submit' | 'reset' | 'button';
  disabled?: boolean;
  onClick?: () => void;
};

export const StandardButton: FC<StandardButtonProps> = (
  props: StandardButtonProps
) => {
  return (
    <Button
      // variant="true"
      color="teal"
      type={props.type}
      disabled={props.disabled}
      onClick={props.onClick}
    >
      {props.value}
    </Button>
  );
};
