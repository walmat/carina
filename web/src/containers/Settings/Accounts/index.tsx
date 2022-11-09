import React, {
  Fragment,
  useState,
  useCallback,
  useMemo,
  useEffect,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import { useTranslation } from "react-i18next";

import { Edit, Trash, Copy, Download, Upload, Move } from "react-feather";

import { Action } from "../../../types";

import { Content, Toolbar, Groups, InfoBar, Table } from "../../../components";

import {
  AccountGroup,
  makeAccounts,
  copyAccount,
  removeAccount,
  copyToGroup,
  moveToGroup,
  createAccountsGroup as createGroupAction,
  removeAccountsGroup as removeGroupAction
} from '../../../stores/Main/reducers/accounts';

import AddAccounts from "./Create";
import { ACCOUNT_COLUMNS } from './constants';
import { TableItems } from "../../../types/table";
import { filterAccounts } from './filterAccounts';
import EditComponent from './Edit';

const Accounts = () => {
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const groups = useSelector(makeAccounts);

  const [create, setCreate] = useState<boolean>(false);
  const [edit, setEdit] = useState<boolean>(false);
  const [createGroup, setCreateGroup] = useState<boolean>(false);
  const [group, setGroup] = useState<AccountGroup>(groups.default);
  const [filter, setFilter] = useState<string>("");
  const [selected, setSelected] = useState<string[]>([]);

  const tableData = useMemo(() => {
    if (group) {
      return Object.values(group.byId);
    }

    return Object.values(groups.default);
  }, [group]);

  useEffect(() => {
    setGroup((prev) => {
      if (!groups[prev.id]) {
        return groups['default'];
      }
      return groups[prev.id];
    });
  }, [groups]);

  const items: TableItems[] = [
    {
      name: "Account",
      actions: [
        {
          name: "Edit",
          Icon: Edit,
          shortcut: ["⌘", "E"],
          onClick: useCallback((id: string) => {
            setSelected([id]);
            setEdit(true);
          }, [group.id, selected])
        },
        {
          name: "Copy",
          Icon: Copy,
          shortcut: ["⌘", "C"],
          onClick: useCallback((id: string) => {
            // if no rows are selected, just use the row we right clicked on
            if (!selected.length) {
              dispatch(copyAccount({ ids: [id], group: group.id }))

              // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              dispatch(copyAccount({ ids: [...selected, id], group: group.id }))

              // use the selected rows
            } else {
              dispatch(copyAccount({ ids: selected, group: group.id }))
            }
          }, [group.id, selected])
        },
        {
          name: "Remove",
          Icon: Trash,
          shortcut: ["⌘", "D"],
          onClick: useCallback((id: string) => {
            // if no rows are selected, just use the row we right clicked on
            if (!selected.length) {
              dispatch(removeAccount({ ids: [id], group: group.id }))

              // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              dispatch(removeAccount({ ids: [...selected, id], group: group.id }))

              // use the selected rows
            } else {
              dispatch(removeAccount({ ids: selected, group: group.id }))
            }
          }, [group.id, selected])
        },
      ],
    },
    {
      name: "Group",
      actions: [
        {
          name: "Copy",
          Icon: Copy,
          onClick: useCallback((newGroup: string) => {
            if (!selected.length) {
              return;
            }

            dispatch(copyToGroup({ oldGroup: group.id, newGroup, ids: selected }));
          }, [group.id, selected])
        },
        {
          name: "Move",
          Icon: Move,
          onClick: useCallback((newGroup: string) => {
            if (!selected.length) {
              return;
            }

            dispatch(moveToGroup({ oldGroup: group.id, newGroup, ids: selected }));
          }, [group.id, selected])
        },
      ],
    },
  ];

  const actions: Action[] = [
    {
      type: "Primary",
      title: "Remove Accounts",
      Icon: Trash,
      onClick: useCallback(
        () => dispatch(removeAccount({ ids: selected, group: group.id })),
        [group.id, selected]
      ),
    },
    {
      type: "First",
      title: "Import Accounts",
      Icon: Download,
      onClick: () => {},
    },
    {
      type: "Second",
      title: "Export Accounts",
      Icon: Upload,
      onClick: () => {},
    },
  ];

  const onFilter = useCallback((e) => {
    setFilter(e.target.value);
  }, []);

  return (
    <Content>
      <Fragment key="Accounts">
        <Toolbar
          title={t("Pages.Accounts")}
          label="Add Accounts"
          onCreateGroup={() => setCreateGroup(true)}
          onCreate={() => setCreate(true)}
          filter={filter}
          onFilter={onFilter}
          groups={groups}
          group={group}
          setGroup={setGroup}
        />
        <InfoBar hasGroups actions={actions} group={group} list="accounts" />
        <Groups
          tabs={{
            titles: ["Add", "Remove"],
            contents: [
              <Groups.AddGroup groups={groups} create={createGroupAction} />,
              <Groups.RemoveGroup remove={removeGroupAction} groups={groups} />,
            ],
          }}
          title="Account Groups"
          open={createGroup}
          setOpen={setCreateGroup}
        />
        {edit && (
          <EditComponent group={group} ids={selected} open={edit} setOpen={setEdit} />
        )}
        {create && (
          <AddAccounts group={group} open={create} setOpen={setCreate} />
        )}
        <Table
          columns={ACCOUNT_COLUMNS}
          addGroup={setCreateGroup}
          group={group}
          groups={groups}
          data={tableData}
          view="Basic"
          filter={filter}
          filterRow={filterAccounts}
          setFilter={setFilter}
          selected={selected}
          setSelected={setSelected}
          items={items}
          actions={[
            {
              name: "Edit",
              onClick: useCallback(
                (id: string) => {
                  setSelected([id]);
                  setEdit(true);
                },
                [selected]
              ),
            },
            {
              name: "Remove",
              onClick: useCallback(
                (id: string) =>
                  dispatch(removeAccount({ ids: [id], group: group.id })),
                [group.id]
              ),
            },
          ]}
        />
      </Fragment>
    </Content>
  );
};

export default Accounts;
