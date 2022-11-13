import { ReactElement, Dispatch, SetStateAction, ComponentType } from 'react';

import { withFormik, FormikBag, FormikProps } from 'formik';
import { Form } from 'formik-semantic-ui-react';
import { Form as NativeForm } from 'semantic-ui-react';
import * as Yup from 'yup';

import { useAppDispatch, useAppSelector } from '@/app/hooks';
import { AppDimmer, UpdateButton } from '@/components';
import {
  updateProblemProperty,
  selectProblemUpdateLoading,
} from '@/features/problem_update';

export interface FormValues {}

export interface FormikFormProps {}

export const problemPropertyEditFormikForm = <
  V extends FormValues,
  P extends FormikFormProps
>(
  workbookId: number,
  problemId: number,
  problemVersion: number,
  problemType: string,
  toField: (values: V) => ReactElement,
  validationSchema: Yup.ObjectSchema<any>,
  propsToValues: (props: P) => V,
  valuesToProperites: (values: V) => { [key: string]: string },
  setValues: (v: V) => void,
  setErrorMessage: Dispatch<SetStateAction<string>>
): ComponentType<P> => {
  const dispatch = useAppDispatch();
  const loading = useAppSelector(selectProblemUpdateLoading);

  const EditForm = (props: FormikProps<V>): ReactElement => {
    const { values, isValid, submitForm } = props;
    const onClick = () => {
      if (isValid) {
        submitForm();
      }
    };

    return (
      <Form>
        {loading ? <AppDimmer /> : <div />}
        <NativeForm.Group inline>
          <NativeForm.Field>{toField(values)}</NativeForm.Field>
          <NativeForm.Field>
            <UpdateButton type="button" disabled={loading} onClick={onClick} />
          </NativeForm.Field>
        </NativeForm.Group>
      </Form>
    );
  };

  return withFormik<P, V>({
    mapPropsToValues: (props: P) => ({
      ...propsToValues(props),
      loading,
    }),
    validationSchema: validationSchema,
    handleSubmit: (values: V, formikBag: FormikBag<P, V>) => {
      console.log('handleSubmit');
      dispatch(
        updateProblemProperty({
          param: {
            workbookId: workbookId,
            problemId: problemId,
            version: problemVersion,
            problemType: problemType,
            properties: valuesToProperites(values),
          },
          postSuccessProcess: () => {},
          postFailureProcess: (error: string) => setErrorMessage(error),
        })
      );
      setValues(values);
    },
  })(EditForm);
};
