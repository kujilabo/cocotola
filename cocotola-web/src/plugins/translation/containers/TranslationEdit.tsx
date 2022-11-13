import React, { useEffect, useState } from 'react';

import { useParams } from 'react-router-dom';
import { Container, Divider, Grid } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import {
  AppBreadcrumb,
  AppDimmer,
  ErrorMessage,
  SuccessMessage,
} from '@/components';

import { TranslationEditFormikForm } from '../components/TranslationEditFormikForm';
import { TranslationNewFormikForm } from '../components/TranslationNewFormikForm';
import { selectTranslationAddLoading } from '../features/translation_add';
import {
  getTranslations,
  selectTranslationGetListLoading,
  selectTranslations,
} from '../features/translation_get_list';
import { selectTranslationUpdateLoading } from '../features/translation_update';
import { TranslationModel } from '../models/translation';

const findTranslationByPos = (
  translations: TranslationModel[],
  pos: number
): TranslationModel => {
  for (let i = 0; i < translations.length; i++) {
    if (translations[i].pos == pos) {
      return translations[i];
    }
  }
  throw 'not found';
};

const removeTranslationByPos = (
  translations: TranslationModel[],
  pos: number
): TranslationModel[] => {
  let index = 0;
  for (let i = 0; i < translations.length; i++) {
    if (translations[i].pos == pos) {
      index = i;
      break;
    }
  }
  translations.splice(index, 1);

  return translations;
};
type ParamTypes = {
  _text: string;
  _pos: string;
};
export const TranslationEdit = (): React.ReactElement => {
  const { _text, _pos } = useParams<ParamTypes>();
  const text = _text || '';
  const pos = +(_pos || '');
  const dispatch = useAppDispatch();
  const translationGetListLoading = useAppSelector(
    selectTranslationGetListLoading
  );
  const translationAddLoading = useAppSelector(selectTranslationAddLoading);
  const translationUpdateLoading = useAppSelector(
    selectTranslationUpdateLoading
  );
  const loading =
    translationGetListLoading ||
    translationAddLoading ||
    translationUpdateLoading;
  const orgTranslations = useAppSelector(selectTranslations);
  const [translations, setTranslations] = useState(orgTranslations);
  const [errorMessage, setErrorMessage] = useState('');
  const [successMessage, setSuccessMessage] = useState('');
  const [newValues, setNewValues] = useState({
    text: text,
    pos: '',
    translated: '',
  });
  const [editValues, setEditValues] = useState({
    lang2: '',
    text: '',
    pos: '',
    translated: '',
    provider: '',
  });
  const localGetTranslations = (text: string, pos: number) => {
    return getTranslations({
      param: {
        text: text,
      },
      postSuccessProcess: (translations: TranslationModel[]) => {
        const t = findTranslationByPos(translations, pos);
        setEditValues({
          lang2: 'ja',
          text: t.text,
          pos: t.pos.toString(),
          translated: t.translated,
          provider: t.provider,
        });
        setTranslations(removeTranslationByPos(translations, pos));
      },
      postFailureProcess: setErrorMessage,
    });
  };

  useEffect(() => {
    const f = async () => {
      await dispatch(localGetTranslations(text, pos));
    };
    f().catch(console.error);
  }, [dispatch, text, pos]);

  const EditFormikForm = TranslationEditFormikForm(
    setSuccessMessage,
    setErrorMessage,
    setEditValues
  );
  const NewFormikForm = TranslationNewFormikForm(
    setSuccessMessage,
    setErrorMessage,
    setNewValues
  );

  return (
    <Container fluid>
      <AppBreadcrumb
        links={[{ text: 'Translations', url: '/plugin/translation/list' }]}
        text={text}
      />
      <Divider hidden />
      {loading ? <AppDimmer /> : <div />}
      <Grid padded>
        <Grid.Row>
          <Grid doubling columns={3}>
            <Grid.Column>
              <EditFormikForm
                index={0}
                selectedLang2={'ja'}
                lang2={editValues.lang2}
                text={editValues.text}
                pos={editValues.pos}
                translated={editValues.translated}
                provider={editValues.provider}
                refreshTranslations={() => {
                  const f = async () => {
                    await dispatch(
                      localGetTranslations(editValues.text, +editValues.pos)
                    );
                  };
                  f().catch(console.error);
                }}
              />
            </Grid.Column>
            <Grid.Column>
              <NewFormikForm
                text={newValues.text}
                pos={newValues.pos}
                translated={newValues.translated}
                refreshTranslations={() => {
                  const f = async () => {
                    await dispatch(
                      localGetTranslations(editValues.text, +editValues.pos)
                    );
                  };
                  f().catch(console.error);
                }}
              />
            </Grid.Column>
          </Grid>
        </Grid.Row>
        <Grid.Row>
          <Grid doubling columns={3}>
            {translations.map((t: TranslationModel, i: number) => {
              return (
                <Grid.Column key={t.pos}>
                  <EditFormikForm
                    index={i}
                    selectedLang2={'ja'}
                    text={t.text}
                    pos={t.pos.toString()}
                    translated={t.translated}
                    lang2={t.lang2}
                    provider={t.provider}
                    refreshTranslations={() => {
                      const f = async () => {
                        await dispatch(
                          localGetTranslations(editValues.text, +editValues.pos)
                        );
                      };
                      f().catch(console.error);
                    }}
                  />
                </Grid.Column>
              );
            })}
          </Grid>
        </Grid.Row>
        <Grid.Row>
          <Grid.Column>
            <SuccessMessage message={successMessage} />
            <ErrorMessage message={errorMessage} />
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Container>
  );
};
