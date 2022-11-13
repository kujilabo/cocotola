import { FC } from 'react';

import { useTranslation } from 'react-i18next';

import { StandardButton } from './StandardButton';

type UpdateButtonProps = {
  type: 'submit' | 'reset' | 'button';
  disabled?: boolean;
  onClick?: () => void;
};

export const UpdateButton: FC<UpdateButtonProps> = (
  props: UpdateButtonProps
) => {
  const [t] = useTranslation();
  return (
    <StandardButton
      type={props.type}
      disabled={props.disabled}
      onClick={props.onClick}
      value={t('Update')}
    />
  );
};
