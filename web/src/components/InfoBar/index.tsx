import React, { Dispatch, SetStateAction, useEffect, useState } from "react";
import styled from "styled-components";

import View from "./view";
import Info from "./info";
import Actions from "./actions";

import { TaskGroup } from '../../stores/Main/reducers/tasks';
import { ProfileGroup } from '../../stores/Main/reducers/profiles';
import { ProxyGroup } from '../../stores/Main/reducers/proxies';
import { AccountGroup } from '../../stores/Main/reducers/accounts';
import { Rates } from '../../stores/Main/reducers/rates';
import { Webhooks } from '../../stores/Main/reducers/webhooks';
import { Integrations } from '../../stores/Main/reducers/integrations';

type Action = {
  type: "Primary" | "First" | "Second" | "Action" | "Button";
  title: string;
  Icon: any;
  onClick?: any;
  headers?: any;
};

type SortBy = {
  id: string;
  desc: boolean;
};

type Props = {
  group:
    | TaskGroup
    | ProfileGroup
    | ProxyGroup
    | AccountGroup
    | Rates
    | Webhooks
    | Integrations;
  list: string;
  actions: Action[];
  view?: string;
  setView?: Dispatch<SetStateAction<string>>;
  sortBy?: SortBy;
  setSortBy?: Dispatch<SetStateAction<SortBy>>;
  hiddenColumns?: string[];
  setHiddenColumns?: Dispatch<SetStateAction<string[]>>;
  groupBy?: string;
  setGroupBy?: Dispatch<SetStateAction<string>>;
  showView?: boolean;
  hasGroups?: boolean;
};

const InfoBar = ({
  group,
  list,
  actions,
  view = "Basic",
  setView = () => {},
  sortBy,
  setSortBy = () => {},
  hiddenColumns,
  setHiddenColumns = () => {},
  groupBy,
  setGroupBy = () => {},
  showView = false,
  hasGroups = false,
}: Props) => {
  const [info, setInfo] = useState(
    !group ? [] : hasGroups ? Object.values(group.byId) : Object.values(group)
  );

  useEffect(() => {
    if (!group) {
      return;
    }

    setInfo(hasGroups ? Object.values(group.byId) : Object.values(group));
  }, [group, hasGroups]);

  return (
    <Container>
      {showView ? <View view={view} setView={setView} /> : null}
      <Info info={info} list={list} />
      <Actions
        actions={actions}
        setSortBy={setSortBy}
        sortBy={sortBy}
        hiddenColumns={hiddenColumns}
        setHiddenColumns={setHiddenColumns}
        groupBy={groupBy}
        setGroupBy={setGroupBy}
      />
    </Container>
  );
};

const Container = styled.div`
  height: 32px;
  margin-top: 16px;
  padding: 16px;
  display: flex;
  border-radius: 4px;
  align-items: center;
  background-color: ${({ theme }) => theme.colors.sidebar};
`;

InfoBar.defaultProps = {
  view: "Basic",
  setView: () => {},
  showView: false,
  hasGroups: false,
};

export default InfoBar;
