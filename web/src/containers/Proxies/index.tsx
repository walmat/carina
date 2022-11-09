import React, {
  Fragment,
  useState,
  useCallback,
  useMemo,
  useEffect,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import { useTranslation } from "react-i18next";

import { Zap, Slash, Edit as EditIcon, Trash, Move, Download, Upload, Copy } from "react-feather";

import { Action } from "../../types";

import { Content, Toolbar, Groups, InfoBar, Table } from "../../components";

import {
  ProxyGroup,
  makeProxies,
  removeProxies,
  copyToGroup,
  moveToGroup,
  addProxiesGroup as createGroupAction,
  removeProxiesGroup as removeGroupAction
} from '../../stores/Main/reducers/proxies';

import AddProxies from "./Create";
import TestProxies from "./Test";
import { PROXY_COLUMNS } from './constants';

import { Items } from "../../components/Table/types";
import { filterProxies } from './filterProxies';
import Edit from './Edit';

const Proxies = () => {
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const groups = useSelector(makeProxies);

  const [view, setView] = useState<string>("Basic");
  const [create, setCreate] = useState<boolean>(false);
  const [edit, setEdit] = useState<boolean>(false);
  const [isTesting, setIsTesting] = useState<boolean>(false);
  const [createGroup, setCreateGroup] = useState<boolean>(false);
  const [group, setGroup] = useState<ProxyGroup>(groups.default);
  const [filter, setFilter] = useState<string>("");
  const [selected, setSelected] = useState<string[]>([]);

  const tableData = useMemo(() => {
    if (group) {
      return Object.values(group.byId);
    }

    return Object.values(groups.default);
  }, [group]);

  const items: Items[] = [
    {
      name: "Proxy",
      actions: [
        {
          name: "Edit",
          Icon: EditIcon,
          shortcut: ["⌘", "E"],
          onClick: useCallback((id: string) => {
            setSelected([id]);
            setEdit(true);
          }, [group.id, selected])
        },
        {
          name: "Remove",
          Icon: Trash,
          shortcut: ["⌘", "D"],
          onClick: useCallback((id: string) => {
            // if no rows are selected, just use the row we right clicked on
            if (!selected.length) {
              dispatch(removeProxies({ ids: [id], group: group.id }))

              // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              dispatch(removeProxies({ ids: [...selected, id], group: group.id }))

              // use the selected rows
            } else {
              dispatch(removeProxies({ ids: selected, group: group.id }))
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
      type: "Action",
      title: "Test Proxies",
      Icon: Zap,
      onClick: () => setIsTesting((prev) => !prev),
    },
    {
      type: "Primary",
      title: "Clear Proxies",
      Icon: Slash,
      onClick: () => {},
    },
    {
      type: "Primary",
      title: "Remove Proxies",
      Icon: Trash,
      onClick: useCallback(
        () => dispatch(removeProxies({ ids: selected, group: group.id })),
        [group.id, selected]
      ),
    },
    {
      type: "First",
      title: "Import Proxies",
      Icon: Download,
      onClick: () => {},
    },
    {
      type: "Second",
      title: "Export Proxies",
      Icon: Upload,
      onClick: () => {},
    },
  ];

  const onFilter = useCallback((e) => {
    setFilter(e.target.value);
  }, []);

  useEffect(() => {
    setGroup((prev) => {
      if (!groups[prev.id]) {
        return groups['default'];
      }
      return groups[prev.id];
    });
  }, [groups]);

  return (
    <Content>
      <Fragment key="Proxies">
        <Toolbar
          title={t("Pages.Proxies")}
          label="Add Proxies"
          onCreateGroup={() => setCreateGroup(true)}
          onCreate={() => setCreate(true)}
          filter={filter}
          onFilter={onFilter}
          groups={groups}
          group={group}
          setGroup={setGroup}
        />
        <InfoBar
          hasGroups
          showView={false}
          actions={actions}
          group={group}
          list="proxies"
          view={view}
          setView={setView}
        />
        {edit && (
          <Edit group={group} ids={selected} open={edit} setOpen={() => setEdit(prev => !prev)} />
        )}
        <AddProxies group={group} open={create} setOpen={setCreate} />
        <TestProxies
          open={isTesting}
          setOpen={setIsTesting}
          selected={selected}
          groupId={group.id}
        />
        <Groups
          tabs={{
            titles: ["Add", "Remove"],
            contents: [
              <Groups.AddGroup groups={groups} create={createGroupAction} />,
              <Groups.RemoveGroup remove={removeGroupAction} groups={groups} />,
            ],
          }}
          title="Proxy Groups"
          open={createGroup}
          setOpen={setCreateGroup}
        />
        <Table
          columns={PROXY_COLUMNS}
          addGroup={setCreateGroup}
          group={group}
          groups={groups}
          data={tableData}
          view={view}
          filter={filter}
          filterRow={filterProxies}
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
                  dispatch(removeProxies({ ids: [id], group: group.id })),
                [group.id]
              ),
            },
          ]}
        />
      </Fragment>
    </Content>
  );
};

export default Proxies;
