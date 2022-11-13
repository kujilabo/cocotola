import { FC } from 'react';

import { useTranslation } from 'react-i18next';

import { StandardButton } from './StandardButton';


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
