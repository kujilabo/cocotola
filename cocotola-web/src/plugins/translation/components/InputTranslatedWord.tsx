import React from 'react';

import { Input } from 'formik-semantic-ui-react';
import { useTranslation } from 'react-i18next';

export type InputTrasnslatedWordProps = {
  disabled?: boolean;
};
export const InputTrasnslatedWord = (
  props: InputTrasnslatedWordProps
): React.ReactElement => {
  const [t] = useTranslation();
  return (
    <Input
      name="translated"
      label={t('Translated')}
      placeholder="translated"
      errorPrompt
      disabled={props.disabled}
    />
  );
};
