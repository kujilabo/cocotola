import { ReactElement, useEffect, useState } from 'react';

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import { Container, Divider } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppBreadcrumb, AppDimmer, ErrorMessage } from '@/components';
import {
  getStat,
  selectStatFailed,
  selectStatLoading,
  selectStatLoaded,
  selectStatHistory,
} from '@/features/stat';
import { StatResultModel } from '@/models/stat';
import { emptyFunction } from '@/utils/util';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);
export const options = {
  responsive: true,
  plugins: {
    legend: {
      position: 'top' as const,
    },
    title: {
      display: true,
      text: 'Chart.js Line Chart',
    },
  },
  scales: {
    y: {
      suggestedMin: 0,
      suggestedMax: 50,
    },
  },
};

export const Dashboard = (): ReactElement => {
  const dispatch = useAppDispatch();
  const loading = useAppSelector(selectStatLoading);
  const loaded = useAppSelector(selectStatLoaded);
  const failed = useAppSelector(selectStatFailed);
  const history = useAppSelector(selectStatHistory);
  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    if (loaded) {
      return;
    }
    const f = async () => {
      await dispatch(
        getStat({
          param: {
            pageNo: 1,
          },
          postSuccessProcess: emptyFunction,
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  }, [dispatch, loaded]);

  if (loading || failed) {
    return <AppDimmer />;
  }

  const labels = history.results.map((x) => x.date);
  const results: StatResultModel[] = history.results;
  const mastered: number[] = results.map((_) => _.mastered);
  const answered: number[] = results.map((_) => _.answered);

  const data = {
    labels: labels,
    datasets: [
      {
        label: 'Mastered',
        data: mastered,
        borderColor: 'rgb(255, 99, 132)',
        backgroundColor: 'rgba(255, 99, 132, 0.5)',
      },
      {
        label: 'Answered',
        data: answered,
        borderColor: 'rgb(53, 162, 235)',
        backgroundColor: 'rgba(53, 162, 235, 0.5)',
      },
    ],
  };
  console.log('data', data);
  return (
    <Container fluid>
      <AppBreadcrumb links={[]} text={''} home={true} />
      <Divider hidden />
      <Line options={options} data={data} />
      <ErrorMessage message={errorMessage} />
    </Container>
  );
};
