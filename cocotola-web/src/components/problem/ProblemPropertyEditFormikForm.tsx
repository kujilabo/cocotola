import { ReactElement, Dispatch, SetStateAction, ComponentType } from 'react';

import { withFormik, FormikBag, FormikProps } from 'formik';
import { Form } from 'formik-semantic-ui-react';
import { Form as NativeForm } from 'semantic-ui-react';
import * as Yup from 'yup';

import { useAppDispatch, useAppSelector } from '@/app/hooks';
import { UpdateButton } from '@/components';
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
  validationSchema: Yup.ObjectSchema<any>,
  setErrorMessage: Dispatch<SetStateAction<string>>,
  setProblem: (t: V) => void,
  propsToValues: (props: P) => V,
  valuesToProperites: (values: V) => { [key: string]: string },
  form: ReactElement
): ComponentType<P> => {
  const dispatch = useAppDispatch();
  const loading = useAppSelector(selectProblemUpdateLoading);

  const EditForm = <V extends FormValues & { loading: boolean }>(
    props: FormikProps<V>
  ): ReactElement => {
    const { values, isValid, submitForm } = props;

    return (
      <Form>
        <NativeForm.Group inline>
          <NativeForm.Field>{form}</NativeForm.Field>
          <NativeForm.Field>
            <UpdateButton
              type="button"
              disabled={values.loading}
              onClick={() => {
                if (isValid) {
                  submitForm();
                }
              }}
            />
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
      setProblem(values);
    },
  })(EditForm);
};
