import { ReactElement, Dispatch, ComponentType, SetStateAction } from 'react';

import { FormikProps, withFormik, FormikBag } from 'formik';
import { Form } from 'formik-semantic-ui-react';
import { useNavigate } from 'react-router-dom';
import { Card, Header } from 'semantic-ui-react';
import * as Yup from 'yup';

import { useAppDispatch, useAppSelector } from '@/app/hooks';
import { AddButton, AppDimmer } from '@/components';
import { addProblem, selectProblemAddLoading } from '@/features/problem_add';

export interface FormValues {}

export interface FormikFormProps {}

export interface problemNewFormikFormArgs<
  V extends FormValues,
  P extends FormikFormProps
> {
  workbookId: number;
  problemType: string;
  toContent: (v: V) => ReactElement;
  validationSchema: Yup.ObjectSchema<any>;
  propsToValues: (props: P) => V;
  valuesToProperties: (values: V) => { [key: string]: string };
  resetValues: (v: V) => void;
  setErrorMessage: Dispatch<SetStateAction<string>>;
}

export const problemNewFormikForm = <
  V extends FormValues,
  P extends FormikFormProps
>(
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
    const { values, isSubmitting } = props;

    return (
      <Form>
        <Card>
          <Card.Content>
            <Header component="h2">New problem</Header>
          </Card.Content>
          <Card.Content>{toContent(values)}</Card.Content>
          <Card.Content>
            {loading ? <AppDimmer /> : <div />}
            <AddButton type="submit" disabled={isSubmitting} />
          </Card.Content>
        </Card>
      </Form>
    );
  };

  return withFormik<P, V>({
    mapPropsToValues: (props: P) => propsToValues(props),
    validationSchema: validationSchema,
    handleSubmit: (values: V, formikBag: FormikBag<P, V>) => {
      dispatch(
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
      resetValues(values);
    },
  })(NewForm);
};
