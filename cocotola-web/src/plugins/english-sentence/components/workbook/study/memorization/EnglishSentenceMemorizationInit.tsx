import { ReactElement, useEffect } from 'react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppDimmer } from '@/components';
import { selectProblemMap } from '@/features/problem_find';
import { selectRecordbook } from '@/features/recordbook_get';
import {
  setEnglishSentenceRecordbook,
  selectTs,
  selectEnglishSentenceRecordbook,
} from '@/plugins/english-sentence/features/english_sentence_study';

export const EnglishSentenceMemorizationInit = (): ReactElement => {
  const dispatch = useAppDispatch();
  const recordbook = useAppSelector(selectRecordbook);
  const problemMap = useAppSelector(selectProblemMap);
  const englishSentenceRecordbook = useAppSelector(
    selectEnglishSentenceRecordbook
  );
  const ts = useAppSelector(selectTs);

  console.log('EnglishSentenceMemorizationInit');
  useEffect(() => {
    dispatch(setEnglishSentenceRecordbook(recordbook));
  }, [dispatch, ts, recordbook]);

  if (englishSentenceRecordbook.records.length === 0) {
    return <div />;
  }

  const problemId = englishSentenceRecordbook.records[0].problemId;
  const problem = problemMap[problemId];
  // onsole.log('englishSentenceRecordbook.records', englishSentenceRecordbook.records);
  // onsole.log('problemMap', problemMap);
  // onsole.log('problemId', problemId);
  // onsole.log('problem', problem);

  if (!problem) {
    return <div>undefined</div>;
  }

  return <AppDimmer />;
};
