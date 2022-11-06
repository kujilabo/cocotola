import { FC } from 'react';

import { Button, Icon } from 'semantic-ui-react';

type AudioButtonProps = {
  id: number;
  disabled?: boolean;
  loadAndPlay: (postFunc: (value: string) => void) => void;
};

export const AudioButton: FC<AudioButtonProps> = (props: AudioButtonProps) => {
  const playAudio = (value: string) => {
    const audio = new Audio('data:audio/wav;base64,' + value);
    audio.play();
  };

  if (props.id === 0) {
    return <div />;
  }

  return (
    <Button.Group floated="right" icon>
      <Button
        basic
        color="teal"
        disabled={props.disabled}
        onClick={() => props.loadAndPlay(playAudio)}
      >
        <Icon name="play" />
      </Button>
    </Button.Group>
  );
};
