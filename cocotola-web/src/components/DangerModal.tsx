import { ReactNode, ReactElement, FC, useState } from 'react';

import { Button, Modal } from 'semantic-ui-react';

type DangerModalProps = {
  triggerValue: string;
  content: string;
  standardValue: string;
  dangerValue: string;
  triggerLayout: (children: ReactNode) => ReactElement;
  standardFunc: () => any; // eslint-disable-line @typescript-eslint/no-explicit-any
  dangerFunc: () => any; // eslint-disable-line @typescript-eslint/no-explicit-any
};

export const DangerModal: FC<DangerModalProps> = (props: DangerModalProps) => {
  const [open, setOpen] = useState(false);
  return (
    <div>
      <Modal
        size="mini"
        trigger={props.triggerLayout(<>{props.triggerValue}</>)}
        // trigger={<div>{props.triggerValue}</div>}
        open={open}
        onClose={() => setOpen(false)}
        onOpen={() => setOpen(true)}
      >
        <Modal.Content>{props.content}</Modal.Content>
        <Modal.Actions>
          <Button
            negative
            onClick={() => {
              props.dangerFunc();
              setOpen(false);
            }}
          >
            {props.dangerValue}
          </Button>
          <Button
            positive
            onClick={() => {
              props.standardFunc();
              setOpen(false);
            }}
          >
            {props.standardValue}
          </Button>
        </Modal.Actions>
      </Modal>
    </div>
  );
};
