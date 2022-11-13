import { ReactElement, useState, useEffect } from 'react';

import { useParams } from 'react-router-dom';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { problemFactory } from '@/app/store';
import { AppDimmer } from '@/components';
// import {
//   getProblem,
//   selectProblem,
//   selectProblemGetLoading,
// } from '@/features/problem_get';import {
import {
  getProblem,
  selectProblemMap,
  // selectProblemLoadingMap,
} from '@/features/problem_find';
import { getWorkbook, selectWorkbook } from '@/features/workbook_get';
import { emptyFunction } from '@/utils/util';

type ParamTypes = {
  _workbookId: string;
  _problemId: string;
};
export const PrivateProblemEdit = (): ReactElement => {
  const { _workbookId, _problemId } = useParams<ParamTypes>();
  const workbookId = +(_workbookId || '');
  const problemId = +(_problemId || '');
  const dispatch = useAppDispatch();
  const workbook = useAppSelector(selectWorkbook);
  const problemMap = useAppSelector(selectProblemMap);
  const problem = problemMap[problemId];
  // const problemLoading = useAppSelector(selectProblemGetLoading);
  // TODO
  const problemLoading = false;
  const [errorMessage, setErrorMessage] = useState('');

  console.log('problem.version', problem.version);
  useEffect(() => {
    dispatch(
      getWorkbook({
        param: { id: workbookId },
        postSuccessProcess: emptyFunction,
        postFailureProcess: setErrorMessage,
      })
    );
  }, [dispatch, workbookId]);

  useEffect(() => {
    console.log(
      'getProblem1',
      workbookId + ',' + problemId + ',' + ',' + problem.version
    );
    if (problemLoading) {
      return;
    }

    console.log(
      'getProblem2',
      workbookId + ',' + problemId + ',' + ',' + problem.version
    );
    dispatch(
      getProblem({
        param: { workbookId: workbookId, problemId: problemId },
        postSuccessProcess: emptyFunction,
        postFailureProcess: setErrorMessage,
      })
    );
    // }, [dispatch, workbookId, problemId]);
  }, [dispatch, workbookId, problemId, problem.version]);

  if (errorMessage !== '') {
    return <div>{errorMessage}</div>;
  }
  if (+(_workbookId || '') !== workbook.id) {
    return <AppDimmer />;
  }
  if (+(_problemId || '') !== problem.id || problemLoading) {
    return <AppDimmer />;
  }

  return problemFactory.createProblemEdit(
    workbook.problemType,
    workbook,
    problem
  );
};
