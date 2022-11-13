import { FC } from 'react';

import { Button } from 'semantic-ui-react';

type DangerButtonProps = {
  value: string;
  type: 'submit' | 'reset' | 'button';
  disabled?: boolean;
  onClick?: () => void;
};

export const DangerButton: FC<DangerButtonProps> = (
  props: DangerButtonProps
) => {
  return (
    <Button
      // variant="true"
      color="red"
      type={props.type}
      disabled={props.disabled}
      onClick={props.onClick}
    >
      {props.value}
    </Button>
  );
};
