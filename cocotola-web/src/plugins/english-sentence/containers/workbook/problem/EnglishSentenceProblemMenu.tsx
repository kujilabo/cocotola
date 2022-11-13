import { FC, useEffect, useState } from 'react';

import { useNavigate } from 'react-router-dom';
import { Label, Menu, Input } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import {
  getCompletionRate,
  selectRecordbookCompletionRateMap,
} from '@/features/recordbook_get';
import { WorkbookModel } from '@/models/workbook';
import { emptyFunction } from '@/utils/util';

type EnglishSentenceProblemMenuProps = {
  initStudy: (s: string) => void;
  workbook: WorkbookModel;
};

export const EnglishSentenceProblemMenu: FC<EnglishSentenceProblemMenuProps> = (
  props: EnglishSentenceProblemMenuProps
) => {
  const dispatch = useAppDispatch();
  const recordbookCompletionRateMap = useAppSelector(
    selectRecordbookCompletionRateMap
  );
  const [errorMessage, setErrorMessage] = useState('');
  const navigate = useNavigate();
  const studyButtonClicked = (studyType: string) => {
    props.initStudy('');
    navigate(`/app/workbook/${props.workbook.id}/study/${studyType}`);
  };
  const onImportButtonClick = () => {
    navigate(`/app/private/workbook/${props.workbook.id}/import`);
  };

  // when workbookId is changed
  useEffect(() => {
    // get the completion rate of the workbook
    const f = async () => {
      await dispatch(
        getCompletionRate({
          param: { workbookId: props.workbook.id },
          postSuccessProcess: emptyFunction,
          postFailureProcess: setErrorMessage,
        })
      );
    };
    f().catch(console.error);
  }, [dispatch, props.workbook.id]);

  console.log('recordbookCompletionRateMap', recordbookCompletionRateMap);
  const memorizationCompRate = recordbookCompletionRateMap['memorization'] ?? 0;
  const dictationCompRate = recordbookCompletionRateMap['dictation'] ?? 0;

  if (props.workbook.subscribed) {
    return (
      <Menu vertical fluid>
        <Menu.Item>
          Study
          <Menu.Menu>
            <Menu.Item onClick={() => studyButtonClicked('memorization')}>
              Memorization
            </Menu.Item>
            <Menu.Item onClick={() => studyButtonClicked('dictation')}>
              Dictation
            </Menu.Item>
          </Menu.Menu>
        </Menu.Item>
      </Menu>
    );
  } else {
    return (
      <Menu vertical fluid>
        <Menu.Item>
          Study
          <Menu.Menu>
            <Menu.Item onClick={() => studyButtonClicked('memorization')}>
              Memorization
              <Label>{memorizationCompRate} %</Label>
            </Menu.Item>
            <Menu.Item onClick={() => studyButtonClicked('dictation')}>
              Dictation
              <Label>{dictationCompRate} %</Label>
            </Menu.Item>
          </Menu.Menu>
        </Menu.Item>
        <Menu.Item>
          Problems
          <Menu.Menu>
            <Menu.Item>
              <Input placeholder="Search..." />
            </Menu.Item>
            <Menu.Item
              onClick={() => {
                navigate(
                  `/app/private/workbook/${props.workbook.id}/problem/new`
                );
              }}
            >
              New problem
            </Menu.Item>
            <Menu.Item onClick={onImportButtonClick}>Import problems</Menu.Item>
          </Menu.Menu>
        </Menu.Item>
        <Menu.Item>
          Workbook settings
          <Menu.Menu>
            <Menu.Item
              onClick={() =>
                navigate(`/app/private/workbook/${props.workbook.id}/edit`)
              }
            >
              Edit workbook
            </Menu.Item>
          </Menu.Menu>
        </Menu.Item>
        <Menu.Item>{errorMessage}</Menu.Item>
      </Menu>
    );
  }
};
