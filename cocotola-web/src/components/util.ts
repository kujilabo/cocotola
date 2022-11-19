import { ChangeEvent, useEffect } from 'react';

export const useDidMount = (func: () => void) =>
  useEffect(() => {
    func();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

export const inputChangeString = (
  f: (v: string) => void
): ((e: ChangeEvent<HTMLInputElement>) => void) => {
  return (e: ChangeEvent<HTMLInputElement>): void => f(e.target.value);
};
