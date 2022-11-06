import { ChangeEvent, useEffect } from 'react';

export const useDidMount = (func: () => void) =>
  useEffect(() => {
    func();
  }, []);

export const inputChangeString = (
  f: (v: string) => void
): ((e: ChangeEvent<HTMLInputElement>) => void) => {
  return (e: ChangeEvent<HTMLInputElement>): void => f(e.target.value);
};
