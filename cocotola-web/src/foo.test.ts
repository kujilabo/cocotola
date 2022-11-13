import { sum } from './foo';

test('basic', () => {
  expect(sum()).toBe(0);
});

test('basic again', () => {
  expect(sum(1, 2)).toBe(3);
});

interface xxx {
  number: number;
  text: string;
  lang2: string;
  translated: string;
  // note: string;
}
test('aa', () => {
  const s: xxx = {
    number: 123,
    text: 'TEXT',
    lang2: 'LANG2',
    translated: 'TRANSLATED',
  };
  const t = { ...s };
  expect(t.number).toBe(124);
});
