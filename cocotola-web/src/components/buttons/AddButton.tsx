import { FC } from 'react';

import { StandardButton } from './StandardButton';

import { useTranslation } from 'react-i18next';

type AddButtonProps = {
  type: 'submit' | 'reset' | 'button';
  disabled: boolean;
  onClick?: () => void;
};

export const AddButton: FC<AddButtonProps> = (props: AddButtonProps) => {
  const [t] = useTranslation();
  return (
    <StandardButton
      type={props.type}
      disabled={props.disabled}
      onClick={props.onClick}
      value={t('Add')}
    />
  );
};
