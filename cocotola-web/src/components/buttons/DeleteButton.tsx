import { FC } from 'react';

import { useTranslation } from 'react-i18next';

import { DangerButton } from '@/components/buttons';

type DeleteButtonProps = {
  type: 'submit' | 'reset' | 'button';
  disabled: boolean;
  onClick?: () => void;
};

export const DeleteButton: FC<DeleteButtonProps> = (
  props: DeleteButtonProps
) => {
  const [t] = useTranslation();
  return (
    <DangerButton
      type={props.type}
      disabled={props.disabled}
      onClick={props.onClick}
      value={t('Delete')}
    />
  );
};
