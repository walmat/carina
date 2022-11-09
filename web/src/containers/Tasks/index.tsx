import React, {
  Fragment,
  useState,
  useCallback,
  useEffect,
  useMemo,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import { useTranslation } from "react-i18next";

import {
  Link2 as Link,
  Settings,
  ExternalLink,
  Download,
  Upload,
  Play,
  Pause,
  Edit as Pencil1,
  Edit2 as Pencil,
  Copy,
  Trash,
  Move,
} from "react-feather";

import { Action } from "../../types";

import { Content, Toolbar, Groups, InfoBar, Table } from "../../components";

import {
  TaskGroup,
  makeTasks,
  startTask,
  stopTask,
  copyTask,
  launchGroup,
  removeTask,
  copyToGroup,
  moveToGroup,
  addTaskGroup as createGroupAction,
  removeTaskGroup as removeGroupAction,
} from "../../stores/Main/reducers/tasks";

import Create from "./Create";
import Edit from './Edit';
import ActionBar from "./ActionBar";
import { Items } from "../../components/Table/types";
import { makeSettings, toggle, SETTINGS } from "../../stores/Main/reducers/settings";
import { filterTask } from './filterTask';
import { headers, getTaskColumns } from './constants';

type SortBy = {
  id: string;
  desc: boolean;
};

const Tasks = () => {
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const groups = useSelector(makeTasks);
  const { isOpen } = useSelector(makeSettings);

  const [collapsed, setCollapsed] = useState(false);
  const [view, setView] = useState<string>("Basic");
  const [edit, setEdit] = useState(false);
  const [createGroup, setCreateGroup] = useState<boolean>(false);
  const [group, setGroup] = useState<TaskGroup>(groups.default);
  const [columns, setColumns] = useState(getTaskColumns(view));
  const [filter, setFilter] = useState<string>("");
  const [selected, setSelected] = useState<string[]>([]);
  const [sortBy, setSortBy] = useState<SortBy>({ id: "id", desc: false });
  const [hiddenColumns, setHiddenColumns] = useState<string[]>([]);
  const [groupBy, setGroupBy] = useState<string>("");

  const tableData = useMemo(() => {
    if (group) {
      return Object.values(group.byId);
    }

    return Object.values(groups.default);
  }, [group, hiddenColumns]);

  useEffect(() => {
    setColumns(getTaskColumns(view));
  }, [view]);

  useEffect(() => {
    setGroup((prev) => {
      if (!groups[prev.id]) {
        return groups['default'];
      }
      return groups[prev.id];
    });
  }, [groups]);

  const items: Items[] = [
    {
      name: "Task",
      actions: [
        {
          name: "Start",
          Icon: Play,
          shortcut: ["⌘", "R"],
          onClick: useCallback((id: string) => {
            // if no rows are selected, just use the row we right clicked on
            if (!selected.length) {
              dispatch(startTask({ ids: [id], group: group.id }))

            // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              dispatch(startTask({ ids: [...selected, id], group: group.id }))

            // use the selected rows
            } else {
              dispatch(startTask({ ids: selected, group: group.id }))
            }
          }, [group.id, selected])
        },
        {
          name: "Stop",
          Icon: Pause,
          shortcut: ["⌘", "S"],
          onClick: useCallback((id: string) => {
            // if no rows are selected, just use the row we right clicked on
            if (!selected.length) {
              dispatch(stopTask({ ids: [id], group: group.id }))

              // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              dispatch(stopTask({ ids: [...selected, id], group: group.id }))

              // use the selected rows
            } else {
              dispatch(stopTask({ ids: selected, group: group.id }))
            }
          }, [group.id, selected])
        },
        {
          name: "Edit",
          Icon: Pencil,
          shortcut: ["⌘", "E"],
          onClick: useCallback((id: string) => {
            // if no rows are selected, just use the row we right clicked on
            if (!selected.length) {
              setSelected([id]);
              // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              setSelected(prev => [...prev, id]);
              // use the selected rows
            }

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
              dispatch(copyTask({ ids: [id], group: group.id }))

              // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              dispatch(copyTask({ ids: [...selected, id], group: group.id }))

              // use the selected rows
            } else {
              dispatch(copyTask({ ids: selected, group: group.id }))
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
              dispatch(removeTask({ ids: [id], group: group.id }))

              // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              dispatch(removeTask({ ids: [...selected, id], group: group.id }))

              // use the selected rows
            } else {
              dispatch(removeTask({ ids: selected, group: group.id }))
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
      title: "Paste Link",
      Icon: Link,
      onClick: () => {},
    },
    {
      type: "Primary",
      title: "Edit Tasks",
      Icon: Pencil1,
      onClick: useCallback(() => {
        if (!selected.length) {
          return;
        }
        setEdit(true);
      }, [group.id, selected])
    },
    {
      type: "Primary",
      title: "Copy Tasks",
      Icon: Copy,
      onClick: useCallback(
        () => dispatch(copyTask({ ids: selected, group: group.id })),
        [group.id, selected]
      ),
    },
    {
      type: "Primary",
      title: "Pop-Out Group",
      Icon: ExternalLink,
      onClick: useCallback(
        () => {
          if (group.id === 'default') {
            return;
          }

          dispatch(launchGroup({ group: group.id }))
        },
        [group.id]
      ),
      // @ts-ignore
      disabled: group.id === "default",
    },
    {
      type: "First",
      title: "Import Tasks",
      Icon: Download,
      onClick: () => {},
    },
    {
      type: "Second",
      title: "Export Tasks",
      Icon: Upload,
      onClick: () => {},
    },
    {
      type: "Primary",
      popover: true,
      title: "Adjust Settings",
      Icon: Settings,
      // @ts-ignore
      headers,
    },
  ];

  const onFilter = useCallback((e) => {
    setFilter(e.target.value);
  }, []);

  return (
    <Content>
      <Fragment key="Tasks">
        <Toolbar
          title={t("Pages.Tasks")}
          label="Add Tasks"
          onCreateGroup={() => setCreateGroup(true)}
          onCreate={() => dispatch(toggle({ field: SETTINGS.IS_OPEN }))}
          filter={filter}
          onFilter={onFilter}
          groups={groups}
          group={group}
          setGroup={setGroup}
        />
        <InfoBar
          hasGroups
          showView
          actions={actions}
          group={group}
          list="tasks"
          view={view}
          setView={setView}
          setSortBy={setSortBy}
          sortBy={sortBy}
          hiddenColumns={hiddenColumns}
          setHiddenColumns={setHiddenColumns}
          groupBy={groupBy}
          setGroupBy={setGroupBy}
        />
        <Groups
          tabs={{
            titles: ["Add", "Remove"],
            contents: [
              <Groups.AddGroup groups={groups} create={createGroupAction} />,
              <Groups.RemoveGroup remove={removeGroupAction} groups={groups} />,
            ],
          }}
          title="Task Groups"
          open={createGroup}
          setOpen={setCreateGroup}
        />
        {edit && (
          <Edit groups={groups} group={group} ids={selected} open={edit} setOpen={() => setEdit(prev => !prev)} />
        )}
        {isOpen && (
          <Create groups={groups} group={group} open={isOpen} setOpen={() => dispatch(toggle({ field: SETTINGS.IS_OPEN }))} />
        )}
        <Table
          columns={columns}
          addGroup={setCreateGroup}
          group={group}
          groups={groups}
          data={tableData}
          view={view}
          filter={filter}
          setFilter={setFilter}
          filterRow={filterTask}
          selected={selected}
          setSelected={setSelected}
          sortBy={sortBy}
          setSortBy={setSortBy}
          hiddenColumns={hiddenColumns}
          hotkeysEnabled={!isOpen && !edit}
          setHiddenColumns={setHiddenColumns}
          groupBy={groupBy}
          setGroupBy={setGroupBy}
          items={items}
          actions={[
            {
              name: "Run",
              onClick: useCallback(
                (type: string, id: string) => {
                  switch (type) {
                    default:
                    case "start":
                      return dispatch(
                        startTask({ ids: [id], group: group.id })
                      );
                    case "stop":
                      return dispatch(stopTask({ ids: [id], group: group.id }));
                  }
                },
                [group.id]
              ),
            },
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
                  dispatch(removeTask({ ids: [id], group: group.id })),
                [group.id]
              )
            },
          ]}
        />
        {tableData.length ? (
          <ActionBar
            selected={selected}
            group={group.id}
            open={collapsed}
            setOpen={setCollapsed}
          />
        ) : null}
      </Fragment>
    </Content>
  );
};

export default Tasks;
