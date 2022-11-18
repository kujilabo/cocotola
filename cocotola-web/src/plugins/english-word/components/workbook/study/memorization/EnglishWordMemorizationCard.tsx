import { ReactElement, FC, Dispatch, SetStateAction } from 'react';

import { Card, Header } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AudioButton } from '@/components';
import { getAudio, selectAudioViewLoading } from '@/features/audio';
import { emptyFunction } from '@/utils/util';

type EnglishWordMemorizationCardProps = {
  workbookId: number;
  problemId: number;
  audioId: number;
  updatedAt: string;
  headerText: string;
  contentList: ReactElement[];
  setErrorMessage: Dispatch<SetStateAction<string>>;
};

export const EnglishWordMemorizationCard: FC<
  EnglishWordMemorizationCardProps
> = (props: EnglishWordMemorizationCardProps): ReactElement => {
  const dispatch = useAppDispatch();
  const audioViewLoading = useAppSelector(selectAudioViewLoading);

  const loadAndPlay = (postFunc: (value: string) => void) => {
    const f = async () => {
      await dispatch(
        getAudio({
          param: {
            workbookId: props.workbookId,
            problemId: props.problemId,
            audioId: props.audioId,
            updatedAt: props.updatedAt,
          },
          postFunc: postFunc,
          postSuccessProcess: emptyFunction,
          postFailureProcess: props.setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  };

  return (
    <div>
      <Card fluid>
        <Card.Content>
          <Header component="h2">
            {props.headerText}
            <AudioButton
              id={props.audioId}
              loadAndPlay={loadAndPlay}
              disabled={audioViewLoading}
            />
          </Header>
        </Card.Content>
        {props.contentList.map((content: ReactElement, i: number) => {
          return <Card.Content key={i}>{content}</Card.Content>;
        })}
      </Card>
    </div>
  );
};
