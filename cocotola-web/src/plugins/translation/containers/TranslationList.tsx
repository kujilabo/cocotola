import React, { useEffect, useState } from 'react';

import { Link } from 'react-router-dom';
import {
  Button,
  Container,
  Card,
  Divider,
  Grid,
  Label,
  Header,
} from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppBreadcrumb, AppDimmer, ErrorMessage } from '@/components';
import { emptyFunction } from '@/utils/util';

import { toDsiplayText } from '../../english-word/utils/util';
import { AlphabetLinks } from '../components/AlphabetLinks';
import { TranslationListMenu } from '../components/TranslationListMenu';
import {
  findTranslations,
  selectTranslationFindLoading,
  selectTranslations,
} from '../features/translation_find';
import { TranslationModel } from '../models/translation';

// type ParamTypes = {
//   _letter: string;
// };

export const TranslationList = (): React.ReactElement => {
  console.log('DRAW TranslationList');
  // const { _letter } = useParams<ParamTypes>();
  const dispatch = useAppDispatch();
  const loading = useAppSelector(selectTranslationFindLoading);
  const translations = useAppSelector(selectTranslations);
  const [errorMessage, setErrorMessage] = useState('');
  const [letter, setLetter] = useState('a');
  const baseUrl = `/plugin/translation`;
  const onAlphabetClick = (letter: string) => setLetter(letter);

  useEffect(() => {
    dispatch(
      findTranslations({
        param: { letter: letter },
        postSuccessProcess: emptyFunction,
        postFailureProcess: setErrorMessage,
      })
    );
  }, [dispatch, letter]);

  if (loading) {
    return <AppDimmer />;
  }

  return (
    <Container fluid>
      <AppBreadcrumb links={[]} text={'Translations'} />
      <Divider hidden />
      <Grid>
        <Grid.Row>
          {loading ? <AppDimmer /> : <div />}
          <Grid.Column mobile={16} tablet={16} computer={3}>
            <TranslationListMenu />
          </Grid.Column>
          <Grid.Column mobile={16} tablet={16} computer={13}>
            <Grid doubling columns={1}>
              <Grid.Column>
                <AlphabetLinks onClick={onAlphabetClick} />
              </Grid.Column>
            </Grid>
            <Grid doubling columns={3}>
              {translations.map((m: TranslationModel, i: number) => {
                return (
                  <Grid.Column key={i}>
                    <Card>
                      <Card.Content>
                        <Header component="h2">{m.text}</Header>
                      </Card.Content>
                      <Card.Content>{toDsiplayText(m.pos)}</Card.Content>
                      <Card.Content>{m.translated}</Card.Content>
                      <Card.Content>
                        <Button.Group>
                          <Button color="teal">
                            <Link
                              style={{ textDecoration: 'none', color: 'white' }}
                              to={`${baseUrl}/${m.text}/${m.pos}/edit`}
                            >
                              Translation
                            </Link>
                          </Button>
                        </Button.Group>
                        <Button.Group floated="right">
                          {m.provider === 'azure' ? (
                            <Label as="a" color="orange" size="tiny" tag>
                              Azure
                            </Label>
                          ) : (
                            <div />
                          )}
                        </Button.Group>
                      </Card.Content>
                    </Card>
                  </Grid.Column>
                );
              })}
            </Grid>
            <ErrorMessage message={errorMessage} />
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Container>
  );
};
