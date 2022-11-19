import { ReactElement, ReactNode, useState, useEffect } from 'react';

import { withFormik, FormikProps } from 'formik';
import { Form, Input } from 'formik-semantic-ui-react';
import { useParams, useNavigate } from 'react-router-dom';
import { Button, Card, Container, Divider, Header } from 'semantic-ui-react';
import * as Yup from 'yup';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import {
  AppDimmer,
  AppBreadcrumb,
  DangerModal,
  ErrorMessage,
} from '@/components';
import {
  getWorkbook,
  selectWorkbook,
  selectWorkbookGetLoading,
} from '@/features/workbook_get';
import {
  removeWorkbook,
  selectWorkbookRemoveLoading,
} from '@/features/workbook_remove';
import {
  updateWorkbook,
  selectWorkbookUpdateLoading,
} from '@/features/workbook_update';
import { WorkbookModel } from '@/models/workbook';

interface OtherProps {
  loading: boolean;
  onRemoveButtonClick: () => void;
}
interface FormValues {
  id: number;
  version: number;
  name: string;
  questionText: string;
}
const InnerForm = (props: OtherProps & FormikProps<FormValues>) => {
  const { isSubmitting, loading, onRemoveButtonClick } = props;
  return (
    <Form>
      <Card>
        <Card.Content>
          <Header component="h2">Edit workbook</Header>
        </Card.Content>
        <Card.Content>
          <Input
            name="name"
            label="Name"
            placeholder="Workbook name"
            errorPrompt
          />
          <Input
            name="questionText"
            label="Question text"
            placeholder=""
            errorPrompt
          />
        </Card.Content>
        <Card.Content>
          {loading ? <AppDimmer /> : <div />}
          <Button
            type="submit"
            // variant="true"
            color="teal"
            disabled={isSubmitting}
          >
            Update
          </Button>

          <DangerModal
            triggerValue="Delete"
            content="Are you sure you want to delete this problem?"
            standardValue="Cancel"
            dangerValue="Delete"
            triggerLayout={(children: ReactNode) => <>{children}</>}
            standardFunc={() => {
              return;
            }}
            dangerFunc={onRemoveButtonClick}
          />
        </Card.Content>
      </Card>
    </Form>
  );
};

type ParamTypes = {
  _workbookId: string;
};
export const PrivateWorkbookEdit = (): ReactElement => {
  const { _workbookId } = useParams<ParamTypes>();
  const workbookId = +(_workbookId || '');
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const workbook = useAppSelector(selectWorkbook);
  const workbookGetLoading = useAppSelector(selectWorkbookGetLoading);
  const workbookUpdateLoading = useAppSelector(selectWorkbookUpdateLoading);
  const workbookRemoveLoading = useAppSelector(selectWorkbookRemoveLoading);
  const [values, setValues] = useState({
    version: 0,
    name: '',
    questionText: '',
  });
  const [errorMessage, setErrorMessage] = useState('');
  const onRemoveButtonClick = () => {
    const f = async () => {
      await dispatch(
        removeWorkbook({
          param: {
            id: workbookId,
            version: values.version,
          },
          postSuccessProcess: () => navigate('/app/private/workbook'),
          postFailureProcess: (error: string) => setErrorMessage(error),
        })
      );
    };
    f().catch(console.error);
  };

  useEffect(() => {
    const f = async () => {
      await dispatch(
        getWorkbook({
          param: { id: workbookId },
          postSuccessProcess: (workbook: WorkbookModel) => {
            setValues({
              version: workbook.version,
              name: workbook.name,
              questionText: workbook.questionText,
            });
          },
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  }, [dispatch, workbookId]);

  const loading =
    workbookGetLoading ||
    workbookUpdateLoading ||
    workbookRemoveLoading ||
    +(_workbookId || '') !== workbook.id;

  interface FormProps {
    id: number;
    version: number;
    name: string;
    questionText: string;
    loading: boolean;
    onRemoveButtonClick: () => void;
  }
  const InnerFormikForm = withFormik<FormProps, FormValues>({
    mapPropsToValues: (props: FormProps) => ({
      id: props.id,
      version: props.version,
      name: props.name,
      questionText: props.questionText,
    }),
    validationSchema: Yup.object().shape({
      name: Yup.string().required('Name is required'),
    }),
    handleSubmit: (formValues: FormValues) => {
      const f = async () => {
        await dispatch(
          updateWorkbook({
            param: { ...formValues },
            postSuccessProcess: () =>
              navigate(`/app/private/workbook/${workbookId}`),
            postFailureProcess: setErrorMessage,
          })
        );
      };
      f().catch(console.error);
      setValues(formValues);
    },
  })(InnerForm);

  return (
    <Container fluid>
      <AppBreadcrumb
        links={[
          { text: 'My Workbooks', url: '/app/private/workbook' },
          { text: workbook.name, url: `/app/private/workbook/${workbookId}` },
        ]}
        text={'Edit'}
      />
      <Divider hidden />
      <InnerFormikForm
        id={workbookId}
        version={values.version}
        name={values.name}
        questionText={values.questionText}
        loading={loading}
        onRemoveButtonClick={onRemoveButtonClick}
      />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};
