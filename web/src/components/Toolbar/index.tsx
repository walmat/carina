import React from "react";
import { Search, Folder } from "react-feather";
import styled from "styled-components";
import { useHotkeys } from "react-hotkeys-hook";

import {
  Input,
  Select,
  Buttons,
  IndicatorSeparator,
  Typography,
} from "../../elements";

import { Tasks, TaskGroup } from '../../stores/Main/reducers/tasks';
import { WorkFlows, WorkFlowGroup } from '../../stores/Main/reducers/workflows';
import { Accounts, AccountGroup } from '../../stores/Main/reducers/accounts';
import { Proxies, ProxyGroup } from '../../stores/Main/reducers/proxies';
import { Profiles, ProfileGroup } from '../../stores/Main/reducers/profiles';
import { Default, Defaults } from "../../stores/Main/reducers/defaults";

import Actions from "../Actions";

const inputPlaceholder = (title: string) => {
  switch (title) {
    default:
    case "Tasks":
      return "Filter tasks";
    case "Workflows":
      return "Filter workflows";
    case "Profiles":
      return "Filter profiles";
    case "Proxies":
      return "Filter proxies";
    case "Accounts":
      return "Filter accounts";
    case "Shipping Rates":
      return "Filter shipping rates";
    case "Quick Tasks":
      return "Filter quick tasks";
    case "Webhooks":
      return "Filter webhooks";
    case "Defaults":
      return "Filter defaults";
  }
};

interface ToolbarProps {
  title: string;
  label?: string;
  simple?: boolean;
  hasGroups?: boolean;
  onCreateGroup?: any;
  showCreateGroup?: boolean;
  onCreate?: any;
  filter?: string;
  onFilter?: any;
  groups?: Tasks | Profiles | Proxies | Accounts | WorkFlows | Defaults;
  group?: TaskGroup | ProfileGroup | ProxyGroup | AccountGroup | WorkFlowGroup | Default | any;
  setGroup?: any;
}

const Toolbar = ({
  title,
  label,
  simple,
  hasGroups,
  onCreateGroup,
  showCreateGroup = true,
  onCreate,
  filter = "",
  onFilter,
  groups,
  group,
  setGroup,
}: ToolbarProps) => {
  const onChange = (
    group: TaskGroup | ProfileGroup | ProxyGroup | AccountGroup | WorkFlowGroup
  ): any => {
    return setGroup(group);
  };

  useHotkeys("command+n,ctrl+n", onCreate, []);

  if (simple) {
    return (
      <Container>
        <Title hasGroups={false}>{title}</Title>
        <Actions simple={simple} />
      </Container>
    );
  }

  if (!hasGroups) {
    return (
      <Container>
        <Title hasGroups={false}>{title}</Title>

        <Filter>
          <Input
            id="filter"
            name="filter"
            onChange={onFilter}
            value={filter || ""}
            placeholder={inputPlaceholder(title)}
            Icon={Search}
          />
        </Filter>

        <Action
          variant="IconButton"
          text={label}
          command="⌘ N"
          width={165}
          height={36}
          onClick={onCreate}
        />

        <Actions simple={simple || false} />
      </Container>
    );
  }

  return (
    <Container>
      <Title hasGroups>{title}</Title>

      <Groups showCreateGroup={showCreateGroup}>
        <Select
          name="group"
          required
          isClearable={false}
          placeholder="Select group"
          components={{
            IndicatorSeparator,
            ClearIndicator: null,
          }}
          value={group}
          getOptionLabel={(option: TaskGroup | ProxyGroup | ProfileGroup) =>
            option.name
          }
          getOptionValue={(option: TaskGroup | ProxyGroup | ProfileGroup) =>
            option.id
          }
          onChange={onChange}
          options={Object.values(groups || {})}
        />
      </Groups>

      {showCreateGroup ? (
        <CreateGroup>
          <Buttons.Primary
            variant="Icon"
            width={36}
            height={36}
            onClick={onCreateGroup}
          >
            <FolderIcon height={14} />
          </Buttons.Primary>
        </CreateGroup>
      ) : null}

      <Filter>
        <Input
          id="filter"
          name="filter"
          onChange={onFilter}
          value={filter}
          placeholder={inputPlaceholder(title)}
          Icon={Search}
        />
      </Filter>

      <Action
        variant="IconButton"
        command="⌘ N"
        text={label}
        width={165}
        height={36}
        onClick={onCreate}
      />

      <Actions simple={simple || false} />
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  height: 36px;
  align-items: flex-start;
`;

const Title = styled(Typography.H2)<{ hasGroups: boolean }>`
  display: flex;
  justify-content: flex-start;
  align-items: center;
  font-size: 24.5px;
  height: 36px;
  flex: 0 1 auto;
  overflow: hidden;
  font-weight: 700;
  margin: 0;
  margin-right: ${({ hasGroups }) => (!hasGroups ? "0" : "32")}px;
  color: ${({ theme }) => theme.colors.h2};
`;

const Action = styled(Buttons.Primary)`
  color: #fff;
  cursor: pointer;
  background-color: ${({ theme }) => theme.colors.primary};
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  text-transform: none;
`;

const FolderIcon = styled(Folder)`
  display: flex;
  cursor: pointer;
  justify-content: center;
  align-items: center;
  margin: auto 0;
`;

const Groups = styled.div<{ showCreateGroup: boolean }>`
  margin-top: 0;
  margin-right: ${({ showCreateGroup }) => (showCreateGroup ? "16px" : "")};
  min-width: 140px;
  max-width: 170px;
  display: flex;
  flex-direction: column;
  flex-basis: 100%;
  flex: 1;
`;

const CreateGroup = styled.div`
  margin: auto 0;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const Filter = styled.div`
  margin: 0 32px;
  flex: 1 1 auto;
  display: flex;
`;

Toolbar.defaultProps = {
  simple: false,
  hasGroups: true,
  filter: "",
};

export default Toolbar;
