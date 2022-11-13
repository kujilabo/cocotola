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
  const workbookId = +(_workbookId || '');
  const studyType = _studyType || '';
  const dispatch = useAppDispatch();
  const workbook = useAppSelector(selectWorkbook);
  const workbookGetLoading = useAppSelector(selectWorkbookGetLoading);
  const problemListLoading = useAppSelector(selectProblemFindLoading);
  const recordbookViewLoading = useAppSelector(selectRecordbookGetLoading);
  const [errorMessage, setErrorMessage] = useState('');

  // find workbook and all problems
  useEffect(() => {
    const f1 = async () => {
      await dispatch(
        getWorkbook({
          param: { id: workbookId },
          postSuccessProcess: emptyFunction,
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f1().catch(console.error);
    const f2 = async () => {
      await dispatch(
        findAllProblems({
          param: {
            workbookId: workbookId,
          },
          postSuccessProcess: emptyFunction,
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f2().catch(console.error);
  }, [dispatch, workbookId]);

  // find recordbook
  useEffect(() => {
    const f = async () => {
      await dispatch(
        getRecordbook({
          param: {
            workbookId: workbookId,
            studyType: studyType,
          },
          postSuccessProcess: () => {
            const now = new Date();
            const ts = now.toISOString();
            dispatch(problemFactory.initProblemStudy(workbook.problemType)(ts));
          },
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  }, [dispatch, workbookId, studyType, workbook.problemType]);

  if (workbookGetLoading || problemListLoading || recordbookViewLoading) {
    return <AppDimmer />;
  } else if (errorMessage !== '') {
    return <div>{errorMessage}</div>;
  }

  return problemFactory.createProblemStudy(workbook.problemType, studyType);
};
