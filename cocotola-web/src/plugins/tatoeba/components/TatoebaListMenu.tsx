import { ReactElement, FC, useState } from 'react';

import { useNavigate } from 'react-router-dom';
import { Icon, Menu, Input } from 'semantic-ui-react';

import { inputChangeString } from '@/components/util';

type TatoebaListMenuProps = {
  keyword: string;
  onSearch: (keyword: string) => void;
};

export const TatoebaListMenu: FC<TatoebaListMenuProps> = (
  props: TatoebaListMenuProps
): ReactElement => {
  const navigate = useNavigate();
  const [keyword, setKeyword] = useState(props.keyword);
  const onImportButtonClick = () => {
    navigate(`/plugin/tatoeba/import`);
  };
  return (
    <Menu vertical fluid>
      <Menu.Item>
        Tatoeba
        <Menu.Menu>
          <Menu.Item>
            <Input
              icon={
                <Icon
                  name="search"
                  inverted
                  circular
                  link
                  onClick={() => props.onSearch(keyword)}
                />
              }
              placeholder="Search..."
              onChange={inputChangeString(setKeyword)}
              value={keyword}
            />
          </Menu.Item>
          <Menu.Item onClick={onImportButtonClick}>
            Import Tatoeba data
          </Menu.Item>
        </Menu.Menu>
      </Menu.Item>
    </Menu>
  );
};
