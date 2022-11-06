import { FC } from 'react';

import { StandardButton } from './StandardButton';

import { useTranslation } from 'react-i18next';

type UploadButtonProps = {
  type: 'submit' | 'reset' | 'button';
  disabled?: boolean;
  onClick?: () => void;
};

export const UploadButton: FC<UploadButtonProps> = (
  props: UploadButtonProps
) => {
  const [t] = useTranslation();
  return (
    <StandardButton
      type={props.type}
      disabled={props.disabled}
      onClick={props.onClick}
      value={t('Upload')}
    />
  );
};
