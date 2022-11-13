import { ReactElement, useState, useEffect } from 'react';

import { useParams } from 'react-router-dom';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { problemFactory } from '@/app/store';
import { getWorkbook, selectWorkbook } from '@/features/workbook_get';
import { emptyFunction } from '@/utils/util';

type ParamTypes = {
  _workbookId: string;
};
export const PrivateProblemNew = (): ReactElement => {
  const { _workbookId } = useParams<ParamTypes>();
  const workbookId = +(_workbookId || '');
  const dispatch = useAppDispatch();
  const workbook = useAppSelector(selectWorkbook);
  const [errorMessage, setErrorMessage] = useState('');

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

  if (errorMessage !== '') {
    return <div>{errorMessage}</div>;
  }
  return problemFactory.createProblemNew(workbook.problemType, workbook);
};
