import { ReactElement } from 'react';

export const makeAuthorUrl = (author: string): string => {
  return `https://tatoeba.org/ja/user/profile/${author}`;
};

export const makeAuthorLink = (author: string): ReactElement => {
  return (
    <a href={makeAuthorUrl(author)} target="_blank">
      {author}
    </a>
  );
};

export const makeStentenceUrl = (sentenceNumber: number): string => {
  return `https://tatoeba.org/ja/sentences/show/${sentenceNumber}`;
};

export const makeSentenceLink = (sentenceNumber: number): ReactElement => {
  return (
    <a href={makeStentenceUrl(sentenceNumber)} target="_blank">
      Sentence
    </a>
  );
};

export const makeLicenseLink = (): ReactElement => {
  return (
    <a href="https://creativecommons.org/licenses/by/2.0/fr/" target="_blank">
      CC BY 2.0 FR
    </a>
  );
};
