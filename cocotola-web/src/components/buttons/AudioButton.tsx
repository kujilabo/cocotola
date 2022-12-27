import { FC, ReactElement } from 'react';

import { Button, Icon } from 'semantic-ui-react';

type AudioButtonProps = {
  id: number;
  disabled?: boolean;
  loadAndPlay: (postFunc: (value: string) => void) => void;
};

export const AudioButton: FC<AudioButtonProps> = (
  props: AudioButtonProps
): ReactElement => {
  const playAudio = (value: string) => {
    const audio = new Audio('data:audio/wav;base64,' + value);
    const f = async () => {
      await audio.play();
    };
    f().catch(console.error);
  };

  if (props.id === 0) {
    return <div />;
  }

  return (
    <div className="ui fluid buttons">
      <Button
        basic
        color="teal"
        disabled={props.disabled}
        onClick={() => props.loadAndPlay(playAudio)}
      >
        <Icon name="play" />
      </Button>
    </div>
  );
};
