import { FC, ReactElement, ReactNode, useState } from 'react';

import { useTranslation } from 'react-i18next';
import { useNavigate, Link } from 'react-router-dom';
import { Card, Label, Grid, Header, Dropdown } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AudioButton, DangerModal, ErrorMessage } from '@/components';
import { getAudio, selectAudioViewLoading } from '@/features/audio';
import { removeProblem } from '@/features/problem_remove';
import { ProblemModel } from '@/models/problem';
import { EnglishSentenceProblemModel } from '@/plugins/english-sentence/models/english-sentence-problem';
import { emptyFunction } from '@/utils/util';

export const EnglishSentenceProblemReadWrite: FC<
  EnglishSentenceProblemReadWriteProps
> = (props: EnglishSentenceProblemReadWriteProps): ReactElement => {
  console.log('EnglishSentenceProblemReadWrite a');
  const workbookId = props.workbookId;
  const problemId = props.problem.id;
  const problemVersion = props.problem.version;
  const dispatch = useAppDispatch();
  const [t] = useTranslation();
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState('');
  const baseUrl = `/app/private/workbook/${workbookId}/problem/${problemId}`;
  const audioViewLoading = useAppSelector(selectAudioViewLoading);
  const problem = EnglishSentenceProblemModel.of(props.problem);
  const loadAndPlay = (postFunc: (value: string) => void) => {
    const f = async () => {
      await dispatch(
        getAudio({
          param: {
            updatedAt: props.problem.updatedAt,
            workbookId: workbookId,
            problemId: problemId,
            audioId: problem.audioId,
          },
          postFunc: postFunc,
          postSuccessProcess: emptyFunction,
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  };
  const onRemoveButtonClick = () => {
    const f = async () => {
      await dispatch(
        removeProblem({
          param: {
            workbookId: workbookId,
            problemId: problemId,
            version: problemVersion,
          },
          postSuccessProcess: () => navigate(props.baseWorkbookPath),
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  };
  // console.log('EnglishSentenceProblemReadWrite b');
  // console.log(props.problem.properties['audioId']);
  return (
    <Card fluid>
      <Card.Content>
        <Card.Header>
          <Header floated="left">{problem.text}</Header>
          <Header floated="right">
            <Dropdown item text="" icon="bars">
              <Dropdown.Menu>
                <Dropdown.Item>
                  <Link to={`${baseUrl}/edit`}>{t('Edit')}</Link>
                </Dropdown.Item>
                <Dropdown.Item>
                  <DangerModal
                    triggerValue={t('Delete')}
                    content="Are you sure you want to delete this problem?"
                    standardValue="Cancel"
                    dangerValue={t('Delete')}
                    triggerLayout={(children: ReactNode) => (
                      <Label color="red">{children}</Label>
                    )}
                    standardFunc={() => {
                      return;
                    }}
                    dangerFunc={onRemoveButtonClick}
                  />
                </Dropdown.Item>
              </Dropdown.Menu>
            </Dropdown>
          </Header>
        </Card.Header>
        <Card.Header textAlign="right"></Card.Header>
      </Card.Content>
      <Card.Content>
        <Grid columns={1}>
          <Grid.Row>
            <Grid.Column>
              <Header component="h2" className="border-bottom g-mb-15">
                {problem.translated}
              </Header>
            </Grid.Column>
          </Grid.Row>
        </Grid>
        <Grid>
          <Grid.Row>
            <Grid.Column>
              <Header component="h2" className="border-bottom g-mb-15">
                {/* {problem.phonetic`} */}
              </Header>
            </Grid.Column>
            <Grid.Column></Grid.Column>
          </Grid.Row>
        </Grid>
      </Card.Content>
      <Card.Content extra>
        {problem.audioId !== 0 ? (
          // <utton.Group floated="left">
          <AudioButton
            id={problem.audioId}
            loadAndPlay={(postFunc: (value: string) => void) =>
              loadAndPlay(postFunc)
            }
            disabled={audioViewLoading}
          />
        ) : (
          // </utton.Group>
          <div />
        )}
      </Card.Content>
      <ErrorMessage message={errorMessage} />
    </Card>
  );
};

type EnglishSentenceProblemReadWriteProps = {
  workbookId: number;
  problem: ProblemModel;
  baseWorkbookPath: string;
};
