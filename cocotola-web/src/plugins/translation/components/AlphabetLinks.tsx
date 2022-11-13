import React from 'react';

import { Link } from 'react-router-dom';

type AlphabetLinksProps = {
  onClick: (event: string) => void;
};

export const AlphabetLinks: React.FC<AlphabetLinksProps> = (
  props: AlphabetLinksProps
) => {
  console.log('DRAW AlphabetLinks');
  const characters = [
    'a',
    'b',
    'c',
    'd',
    'e',
    'f',
    'g',
    'h',
    'i',
    'j',
    'k',
    'l',
    'm',
    'n',
    'o',
    'p',
    'q',
    'r',
    's',
    't',
    'u',
    'v',
    'w',
    'x',
    'y',
    'z',
  ];
  return (
    <div>
      {characters.map((c: string, i: number) => {
        if (i !== characters.length - 1) {
          return (
            <span key={i}>
              <Link to="#" onClick={() => props.onClick(c)}>
                {c}
              </Link>
              {' / '}
            </span>
          );
        } else {
          return (
            <span key={i}>
              <Link key={i} to="#" onClick={() => props.onClick(c)}>
                {c}
              </Link>
            </span>
          );
        }
      })}
    </div>
  );
};
