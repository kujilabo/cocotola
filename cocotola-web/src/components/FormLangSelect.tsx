import { FC, ChangeEvent } from 'react';

import { Form, Select } from 'semantic-ui-react';

type FormLangSelectProps = {
  label?: string;
  value: string;
  onChange: (e: ChangeEvent<HTMLSelectElement>, d: any) => any; // eslint-disable-line @typescript-eslint/no-explicit-any
};

export const langOptions = [{ key: 'ja', text: 'ja', value: 'ja' }];

export const FormLangSelect: FC<FormLangSelectProps> = (
  props: FormLangSelectProps
) => {
  const label = props.label || 'Lang';
  return (
    <Form.Field
      control={Select}
      options={langOptions}
      label={label}
      onChange={props.onChange}
      value={props.value}
    />
  );
};
