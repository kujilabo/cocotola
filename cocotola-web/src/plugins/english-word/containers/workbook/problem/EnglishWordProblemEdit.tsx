import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Container, Divider } from 'semantic-ui-react';

import { useAppDispatch, useAppSelector } from 'app/hooks';
import { ErrorMessage, PrivateProblemBreadcrumb } from 'components';
import { selectProblem } from 'features/problem_get';
import { WorkbookModel } from 'models/workbook';
import { ProblemModel } from 'models/problem';
import {
  findTatoebaSentences,
  selectTatoebaSentences,
  selectTatoebaFindLoading,
} from 'plugins/tatoeba/features/tatoeba_find';
import { emptyFunction } from 'utils/util';

import { englishWordProblemEditFormikForm } from '../../../components/workbook/problem/EnglishWordProblemEditFormikForm';
import { EnglishWordProblemModel } from '../../../models/english-word-problem';

import 'App.css';

type ParamTypes = {
  _workbookId: string;
  _problemId: string;
};

export const EnglishWordProblemEdit: React.FC<EnglishWordProblemEditProps> = (
  props: EnglishWordProblemEditProps
) => {
  const { _workbookId } = useParams<ParamTypes>();
  const workbookId = +_workbookId;
  const dispatch = useAppDispatch();
  const problem = EnglishWordProblemModel.of(useAppSelector(selectProblem));
  const tatoebaSentences = useAppSelector(selectTatoebaSentences);
  const tatoebaSentenceFindLoading = useAppSelector(selectTatoebaFindLoading);
  const [values, setValues] = useState({
    number: problem.number,
    text: problem.text,
    pos: problem.pos,
    lang2: problem.lang2,
    translated: problem.translated,
    exampleSentenceText: problem.sentence1.text,
    exampleSentenceTranslated: problem.sentence1.translated,
    exampleSentenceNote: problem.sentence1.note,
    sentenceProvider: '',
    tatoebaSentenceNumber1: '',
    tatoebaSentenceNumber2: '',
  });
  const [errorMessage, setErrorMessage] = useState('');
  console.log('values.text', values.text);
  useEffect(() => {
    setValues({
      ...values,
      number: problem.number,
      text: problem.text,
      pos: problem.pos,
      lang2: problem.lang2,
      translated: problem.translated,
      exampleSentenceText: problem.sentence1.text,
      exampleSentenceTranslated: problem.sentence1.translated,
      exampleSentenceNote: problem.sentence1.note,
    });
  }, [dispatch, problem.id, problem.version]);
  useEffect(() => {
    if (values.text.length === 0) {
      return;
    }
    if (tatoebaSentenceFindLoading) {
      return;
    }
    dispatch(
      findTatoebaSentences({
        param: {
          pageNo: 1,
          pageSize: 10,
          keyword: values.text,
          random: true,
        },
        postSuccessProcess: emptyFunction,
        postFailureProcess: setErrorMessage,
      })
    );
  }, [dispatch, values.text]);
  const selectSentence = (index: number, checked: boolean): void => {
    // dispatch(selectTatoebaSentence);
  };

  const EnglishWordProblemEditFormikForm = englishWordProblemEditFormikForm(
    workbookId,
    problem,
    setErrorMessage,
    setValues,
    selectSentence
  );

  console.log('note', values.exampleSentenceNote);

  return (
    <Container fluid>
      <PrivateProblemBreadcrumb
        name={props.workbook.name}
        id={+_workbookId}
        text={'' + props.problem.number}
      />
      <Divider hidden />
      <EnglishWordProblemEditFormikForm
        number={values.number}
        text={values.text}
        pos={values.pos}
        lang2={values.lang2}
        translated={values.translated}
        exampleSentenceText={values.exampleSentenceText}
        exampleSentenceTranslated={values.exampleSentenceTranslated}
        exampleSentenceNote={values.exampleSentenceNote}
        sentenceProvider={values.sentenceProvider}
        tatoebaSentenceNumber1={values.tatoebaSentenceNumber1}
        tatoebaSentenceNumber2={values.tatoebaSentenceNumber2}
        tatoebaSentences={tatoebaSentences}
      />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};

type EnglishWordProblemEditProps = {
  workbook: WorkbookModel;
  problem: ProblemModel;
};
