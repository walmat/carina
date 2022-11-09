import React, { Fragment, useState, useCallback } from "react";
import { useSelector } from "react-redux";
import { useTranslation } from "react-i18next";

import { Content, Toolbar, Groups } from "../../components";

import {
  WorkFlowGroup,
  makeWorkflows,
  createGroup as createGroupAction,
  removeWorkflow as removeGroupAction
} from '../../stores/Main/reducers/workflows';

const Workflows = () => {
  const groups = useSelector(makeWorkflows);
  const { t } = useTranslation();

  const [create, setCreate] = useState<boolean>(false);
  const [createGroup, setCreateGroup] = useState<boolean>(false);
  const [group, setGroup] = useState<WorkFlowGroup>(groups.default);
  const [filter, setFilter] = useState<string>("");

  const onFilter = useCallback((e) => {
    setFilter(e.target.value);
  }, []);

  return (
    <Content>
      <Fragment key="Workflows">
        <Toolbar
          title={t("Pages.Workflows")}
          label="Add Trigger"
          onCreateGroup={() => setCreateGroup(true)}
          onCreate={() => setCreate(true)}
          filter={filter}
          onFilter={onFilter}
          groups={groups}
          group={group}
          setGroup={setGroup}
        />
        <Groups
          tabs={{
            titles: ["Add", "Remove"],
            contents: [
              <Groups.AddGroup groups={groups} create={createGroupAction} />,
              <Groups.RemoveGroup remove={removeGroupAction} groups={groups} />,
            ],
          }}
          title="Workflows"
          open={createGroup}
          setOpen={setCreateGroup}
        />
      </Fragment>
    </Content>
  );
};

export default Workflows;
