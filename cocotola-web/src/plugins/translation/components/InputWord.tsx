import React from 'react';

import { Input } from 'formik-semantic-ui-react';
import { useTranslation } from 'react-i18next';

export type InputWordProps = {
  disabled?: boolean;
};
export const InputWord = (props: InputWordProps): React.ReactElement => {
  const [t] = useTranslation();
  return (
    <Input
      name="text"
      label={String(t('Word'))}
      placeholder="english word"
      errorPrompt
      disabled={props.disabled}
    />
  );
};
