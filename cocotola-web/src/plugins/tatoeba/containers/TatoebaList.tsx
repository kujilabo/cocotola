import React, { useEffect, useState } from 'react';

import { Container, Card, Divider, Grid, Header } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppBreadcrumb, AppDimmer, ErrorMessage } from '@/components';
import { emptyFunction } from '@/utils/util';

import { TatoebaListMenu } from '../components/TatoebaListMenu';
import {
  findTatoebaSentences,
  selectTatoebaFindLoading,
  selectTatoebaSentences,
} from '../features/tatoeba_find';
import { TatoebaSentencePairModel } from '../models/tatoeba';

const makeAuthorUrl = (author: string): string => {
  return `https://tatoeba.org/ja/user/profile/${author}`;
};

const makeStentenceUrl = (sentenceNumber: number): string => {
  return `https://tatoeba.org/ja/sentences/show/${sentenceNumber}`;
};

export const TatoebaList = (): React.ReactElement => {
  const dispatch = useAppDispatch();
  const loading = useAppSelector(selectTatoebaFindLoading);
  const tatoebaSentences = useAppSelector(selectTatoebaSentences);
  const [errorMessage, setErrorMessage] = useState('');
  const [keyword, setKeyword] = useState('');
  console.log('keyword', keyword);
  const onSearch = (value: string) => {
    dispatch(
      findTatoebaSentences({
        param: {
          pageNo: 1,
          pageSize: 10,
          keyword: value,
          random: true,
        },
        postSuccessProcess: emptyFunction,
        postFailureProcess: setErrorMessage,
      })
    );
    setKeyword(value);
  };
  useEffect(() => {
    dispatch(
      findTatoebaSentences({
        param: {
          pageNo: 1,
          pageSize: 10,
          keyword: '',
          random: true,
        },
        postSuccessProcess: emptyFunction,
        postFailureProcess: setErrorMessage,
      })
    );
  }, [dispatch]);

  if (loading) {
    return <AppDimmer />;
  }

  return (
    <Container fluid>
      <AppBreadcrumb links={[]} text={'Tatoeba'} />
      <Divider hidden />
      <Grid>
        <Grid.Row>
          {loading ? <AppDimmer /> : <div />}
          <Grid.Column mobile={16} tablet={16} computer={3}>
            <TatoebaListMenu keyword={keyword} onSearch={onSearch} />
          </Grid.Column>
          <Grid.Column mobile={16} tablet={16} computer={13}>
            <Grid doubling columns={1}>
              {tatoebaSentences.map(
                (m: TatoebaSentencePairModel, i: number) => {
                  return (
                    <Grid.Column key={i}>
                      <Card fluid>
                        <Card.Content>
                          <Header component="h2">{m.src.text}</Header>
                        </Card.Content>
                        <Card.Content>
                          <p>
                            <a
                              href={makeStentenceUrl(m.src.sentenceNumber)}
                              target="_blank"
                            >
                              {makeStentenceUrl(m.src.sentenceNumber)}
                            </a>
                          </p>
                          <p>
                            Author:{' '}
                            <a
                              href={makeAuthorUrl(m.src.author)}
                              target="_blank"
                            >
                              {m.src.author}
                            </a>
                          </p>
                          <p>
                            License:{' '}
                            <a
                              href="https://creativecommons.org/licenses/by/2.0/fr/"
                              target="_blank"
                            >
                              CC BY 2.0 FR
                            </a>
                          </p>
                        </Card.Content>
                      </Card>
                    </Grid.Column>
                  );
                }
              )}
            </Grid>
            <ErrorMessage message={errorMessage} />
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Container>
  );
};
