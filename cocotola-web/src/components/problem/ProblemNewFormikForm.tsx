import { ReactElement, Dispatch, ComponentType, SetStateAction } from 'react';

import { FormikProps, withFormik } from 'formik';
import { Form } from 'formik-semantic-ui-react';
import { useNavigate } from 'react-router-dom';
import { Card, Header } from 'semantic-ui-react';
import * as Yup from 'yup';

import { useAppDispatch, useAppSelector } from '@/app/hooks';
import { AddButton, AppDimmer } from '@/components';
import { addProblem, selectProblemAddLoading } from '@/features/problem_add';

export interface problemNewFormikFormArgs<V extends object, P extends object> {
  workbookId: number;
  problemType: string;
  toContent: (props: FormikProps<V>) => ReactElement;
  validationSchema: Yup.ObjectSchema<any>; // eslint-disable-line @typescript-eslint/no-explicit-any
  propsToValues: (props: P) => V;
  valuesToProperties: (values: V) => { [key: string]: string };
  resetValues: (v: V) => void;
  setErrorMessage: Dispatch<SetStateAction<string>>;
}

export const ProblemNewFormikForm = <V extends object, P extends object>(
  args: problemNewFormikFormArgs<V, P>
): ComponentType<P> => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const loading = useAppSelector(selectProblemAddLoading);

  const {
    workbookId,
    problemType,
    toContent,
    validationSchema,
    propsToValues,
    valuesToProperties,
    resetValues,
    setErrorMessage,
  } = args;

  const NewForm = (props: FormikProps<V>): ReactElement => {
    return (
      <Form>
        <Card>
          <Card.Content>
            <Header component="h2">New problem</Header>
          </Card.Content>
          <Card.Content>{toContent(props)}</Card.Content>
          <Card.Content>
            {loading ? <AppDimmer /> : <div />}
            <AddButton type="submit" disabled={props.isSubmitting} />
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
          addProblem({
            workbookId: workbookId,
            param: {
              problemType: problemType,
              properties: valuesToProperties(values),
            },
            postSuccessProcess: () =>
              navigate(`/app/private/workbook/${workbookId}`),
            postFailureProcess: setErrorMessage,
          })
        );
      };
      f().catch(console.error);
      resetValues(values);
    },
  })(NewForm);
};
