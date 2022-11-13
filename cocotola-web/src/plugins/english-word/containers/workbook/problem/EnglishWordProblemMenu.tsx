import React, { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import { Label, Menu, Input } from 'semantic-ui-react';

import { useAppSelector, useAppDispatch } from 'app/hooks';
import { WorkbookModel } from 'models/workbook';
import { emptyFunction } from 'utils/util';
import {
  getCompletionRate,
  selectRecordbookCompletionRateMap,
} from 'features/recordbook_get';
import 'App.css';

export const EnglishWordProblemMenu: React.FC<EnglishWordProblemMenuProps> = (
  props: EnglishWordProblemMenuProps
) => {
  const dispatch = useAppDispatch();
  const recordbookCompletionRateMap = useAppSelector(
    selectRecordbookCompletionRateMap
  );
  const [errorMessage, setErrorMessage] = useState('');
  const history = useHistory();
  const studyButtonClicked = (studyType: string) => {
    props.initStudy('');
    history.push(`/app/workbook/${props.workbook.id}/study/${studyType}`);
  };
  const onImportButtonClick = () => {
    history.push(`/app/private/workbook/${props.workbook.id}/import`);
  };
  // when workbookId is changed
  useEffect(() => {
    // get the completion rate of the workbook
    dispatch(
      getCompletionRate({
        param: { workbookId: props.workbook.id },
        postSuccessProcess: emptyFunction,
        postFailureProcess: setErrorMessage,
      })
    );
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
                history.push(
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
                history.push(`/app/private/workbook/${props.workbook.id}/edit`)
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
type EnglishWordProblemMenuProps = {
  initStudy: (s: string) => void;
  workbook: WorkbookModel;
};
