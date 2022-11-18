import { ReactElement, useEffect } from 'react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppDimmer } from '@/components';
import { selectProblemMap } from '@/features/problem_find';
import { selectRecordbook } from '@/features/recordbook_get';
import {
  setEnglishWordRecordbook,
  selectTs,
  selectEnglishWordRecordbook,
} from '@/plugins/english-word/features/english_word_study';

export const EnglishWordMemorizationInit = (): ReactElement => {
  const dispatch = useAppDispatch();
  const recordbook = useAppSelector(selectRecordbook);
  const problemMap = useAppSelector(selectProblemMap);
  const englishWordRecordbook = useAppSelector(selectEnglishWordRecordbook);
  const ts = useAppSelector(selectTs);

  console.log('EnglishWordMemorizationInit');
  useEffect(() => {
    dispatch(setEnglishWordRecordbook(recordbook));
  }, [dispatch, ts, recordbook]);

  if (englishWordRecordbook.records.length === 0) {
    return <div />;
  }
  const problemId = englishWordRecordbook.records[0].problemId;
  const problem = problemMap[problemId];
  // onsole.log('englishWordRecordbook.records', englishWordRecordbook.records);
  // onsole.log('problemMap', problemMap);
  // onsole.log('problemId', problemId);
  // onsole.log('problem', problem);

  if (!problem) {
    return <div>undefined</div>;
  }

  return <AppDimmer />;
};
