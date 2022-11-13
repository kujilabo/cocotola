import { ReactElement, useState, useEffect } from 'react';

import { useParams } from 'react-router-dom';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { problemFactory } from '@/app/store';
import { AppDimmer } from '@/components';
import { getProblem, selectProblemMap } from '@/features/problem_find';
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
  const problemVersion = problem ? problem.version : 0;
  useEffect(() => {
    const f = async () => {
      await dispatch(
        getWorkbook({
          param: { id: workbookId },
          postSuccessProcess: emptyFunction,
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  }, [dispatch, workbookId]);

  useEffect(() => {
    if (problemLoading) {
      return;
    }

    if (!problem || (problem && problemId !== problem.id)) {
      const f = async () => {
        await dispatch(
          getProblem({
            param: { workbookId: workbookId, problemId: problemId },
            postSuccessProcess: emptyFunction,
            postFailureProcess: setErrorMessage,
          })
        );
      };
      f().catch(console.error);
    }
  }, [
    dispatch,
    workbookId,
    problemId,
    problem,
    problemVersion,
    problemLoading,
  ]);

  if (errorMessage !== '') {
    return <div>{errorMessage}</div>;
  }
  if (!problem) {
    return <AppDimmer />;
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
