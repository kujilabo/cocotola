import { ReactElement } from 'react';

import { Input } from 'formik-semantic-ui-react';
import { useTranslation } from 'react-i18next';

export type InputTrasnslatedWordProps = {
  disabled?: boolean;
};
export const InputTrasnslatedWord = (
  props: InputTrasnslatedWordProps
): ReactElement => {
  const [t] = useTranslation();
  return (
    <Input
      name="translated"
      label={String(t('Translated'))}
      placeholder="translated"
      errorPrompt
      disabled={props.disabled}
    />
  );
};
