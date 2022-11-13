import { useState, useEffect } from 'react';

import { useParams } from 'react-router-dom';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { problemFactory } from '@/app/store';
import { AppDimmer } from '@/components';
import {
  findAllProblems,
  selectProblemFindLoading,
} from '@/features/problem_find';
import {
  getRecordbook,
  selectRecordbookGetLoading,
} from '@/features/recordbook_get';
import {
  getWorkbook,
  selectWorkbook,
  selectWorkbookGetLoading,
} from '@/features/workbook_get';
import { emptyFunction } from '@/utils/util';

type ParamTypes = {
  _workbookId: string;
  _studyType: string;
};
export const WorkbookStudy = () => {
  const { _workbookId, _studyType } = useParams<ParamTypes>();
  const dispatch = useAppDispatch();
  const workbook = useAppSelector(selectWorkbook);
  const workbookGetLoading = useAppSelector(selectWorkbookGetLoading);
  const problemListLoading = useAppSelector(selectProblemFindLoading);
  const recordbookViewLoading = useAppSelector(selectRecordbookGetLoading);
  const [errorMessage, setErrorMessage] = useState('');

  // find workbook and all problems
  useEffect(() => {
    dispatch(
      getWorkbook({
        param: { id: +(_workbookId || '') },
        postSuccessProcess: emptyFunction,
        postFailureProcess: setErrorMessage,
      })
    );
    dispatch(
      findAllProblems({
        param: {
          workbookId: +(_workbookId || ''),
        },
        postSuccessProcess: emptyFunction,
        postFailureProcess: setErrorMessage,
      })
    );
  }, [dispatch, _workbookId]);

  // find recordbook
  useEffect(() => {
    dispatch(
      getRecordbook({
        param: {
          workbookId: +(_workbookId || ''),
          studyType: _studyType || '',
        },
        postSuccessProcess: () => {
          const now = new Date();
          const ts = now.toISOString();
          dispatch(problemFactory.initProblemStudy(workbook.problemType)(ts));
        },
        postFailureProcess: setErrorMessage,
      })
    );
  }, [dispatch, _workbookId, workbook.problemType]);

  if (workbookGetLoading || problemListLoading || recordbookViewLoading) {
    return <AppDimmer />;
  } else if (errorMessage !== '') {
    return <div>{errorMessage}</div>;
  }

  return problemFactory.createProblemStudy(
    workbook.problemType,
    _studyType || ''
  );
};
