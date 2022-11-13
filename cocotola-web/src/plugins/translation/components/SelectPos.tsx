import React from 'react';

import { Select } from 'formik-semantic-ui-react';
import { useTranslation } from 'react-i18next';

export const posOptions = [
  { key: 'ADJ', text: 'ADJ', value: '1' },
  { key: 'ADV', text: 'ADV', value: '2' },
  { key: 'CONJ', text: 'CONJ', value: '3' },
  { key: 'DET', text: 'DET', value: '4' },
  { key: 'NOUN', text: 'NOUN', value: '6' },
  { key: 'PREP', text: 'PREP', value: '7' },
  { key: 'PRON', text: 'PRON', value: '8' },
  { key: 'VERB', text: 'VERB', value: '9' },
  { key: '', text: '', value: '99' },
];
export type SelectPosProps = {
  disabled?: boolean;
};
export const SelectPos = (props: SelectPosProps): React.ReactElement => {
  const [t] = useTranslation();
  return (
    <Select
      name="pos"
      label={t('Pos')}
      options={posOptions}
      errorPrompt
      disabled={props.disabled}
    />
  );
};
