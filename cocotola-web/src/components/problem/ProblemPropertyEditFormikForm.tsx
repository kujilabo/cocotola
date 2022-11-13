import { ReactElement, Dispatch, SetStateAction, ComponentType } from 'react';

import { withFormik, FormikProps } from 'formik';
import { Form } from 'formik-semantic-ui-react';
import { Form as NativeForm } from 'semantic-ui-react';
import * as Yup from 'yup';

import { useAppDispatch, useAppSelector } from '@/app/hooks';
import { AppDimmer, UpdateButton } from '@/components';
import {
  updateProblemProperty,
  selectProblemUpdateLoading,
} from '@/features/problem_update';
import { emptyFunction } from '@/utils/util';

export interface problemPropertyEditFormikFormArgs<
  V extends object,
  P extends object
> {
  workbookId: number;
  problemId: number;
  problemVersion: number;
  problemType: string;
  toField: (v: V) => ReactElement;
  validationSchema: Yup.ObjectSchema<any>;
  propsToValues: (props: P) => V;
  valuesToProperties: (values: V) => { [key: string]: string };
  resetValues: (v: V) => void;
  setErrorMessage: Dispatch<SetStateAction<string>>;
}

export const ProblemPropertyEditFormikForm = <
  V extends object,
  P extends object
>(
  args: problemPropertyEditFormikFormArgs<V, P>
): ComponentType<P> => {
  const dispatch = useAppDispatch();
  const loading = useAppSelector(selectProblemUpdateLoading);

  const {
    workbookId,
    problemId,
    problemVersion,
    problemType,
    toField,
    validationSchema,
    propsToValues,
    valuesToProperties,
    resetValues,
    setErrorMessage,
  } = args;

  const EditForm = (props: FormikProps<V>): ReactElement => {
    const { values, isValid, submitForm } = props;
    const onClick = () => {
      if (isValid) {
        const f = async () => {
          await submitForm();
        };
        f().catch(console.error);
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
    handleSubmit: (values: V) => {
      console.log('handleSubmit');
      const f = async () => {
        await dispatch(
          updateProblemProperty({
            param: {
              workbookId: workbookId,
              problemId: problemId,
              version: problemVersion,
              problemType: problemType,
              properties: valuesToProperties(values),
            },
            postSuccessProcess: () => emptyFunction,
            postFailureProcess: (error: string) => setErrorMessage(error),
          })
        );
      };
      f().catch(console.error);
      resetValues(values);
    },
  })(EditForm);
};
