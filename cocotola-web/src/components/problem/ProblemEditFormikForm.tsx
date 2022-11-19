import { ReactElement, Dispatch, ComponentType, SetStateAction } from 'react';

import { FormikProps, withFormik } from 'formik';
import { Form } from 'formik-semantic-ui-react';
import { useNavigate } from 'react-router-dom';
import { Card, Header } from 'semantic-ui-react';
import * as Yup from 'yup';

import { useAppDispatch, useAppSelector } from '@/app/hooks';
import { AppDimmer, UpdateButton } from '@/components';
import {
  updateProblem,
  selectProblemUpdateLoading,
} from '@/features/problem_update';

export interface problemEditFormikFormArgs<V extends object, P extends object> {
  workbookId: number;
  problemId: number;
  problemVersion: number;
  problemType: string;
  toContent: (props: FormikProps<V>) => ReactElement;
  validationSchema: Yup.ObjectSchema<any>; // eslint-disable-line @typescript-eslint/no-explicit-any
  propsToValues: (props: P) => V;
  valuesToProperties: (values: V) => { [key: string]: string };
  resetValues: (v: V) => void;
  setErrorMessage: Dispatch<SetStateAction<string>>;
}

export const ProblemEditFormikForm = <V extends object, P extends object>(
  args: problemEditFormikFormArgs<V, P>
): ComponentType<P> => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const loading = useAppSelector(selectProblemUpdateLoading);

  const {
    workbookId,
    problemId,
    problemVersion,
    problemType,
    toContent,
    validationSchema,
    propsToValues,
    valuesToProperties,
    resetValues,
    setErrorMessage,
  } = args;

  const EditForm = (props: FormikProps<V>): ReactElement => {
    // const { values, isSubmitting } = props;
    return (
      <Form>
        <Card fluid>
          <Card.Content>
            <Header component="h2">Edit problem</Header>
          </Card.Content>
          <Card.Content>{toContent(props)}</Card.Content>
          <Card.Content>
            {loading ? <AppDimmer /> : <div />}
            <UpdateButton type="submit" disabled={props.isSubmitting} />
          </Card.Content>
        </Card>
      </Form>
    );
  };

  return withFormik<P, V>({
    mapPropsToValues: (props: P) => propsToValues(props),
    validationSchema: validationSchema,
    handleSubmit: (values: V) => {
      const f = async () => {
        await dispatch(
          updateProblem({
            param: {
              workbookId: workbookId,
              problemId: problemId,
              version: problemVersion,
              problemType: problemType,
              properties: valuesToProperties(values),
            },
            postSuccessProcess: () =>
              navigate(`/app/private/workbook/${workbookId}`),
            postFailureProcess: (error: string) => setErrorMessage(error),
          })
        );
      };
      f().catch(console.error);
      resetValues(values);
    },
  })(EditForm);
};
