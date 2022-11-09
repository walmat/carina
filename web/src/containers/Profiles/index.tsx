import {
  Fragment,
  useState,
  useCallback,
  useMemo,
  useEffect,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import { useTranslation } from "react-i18next";

import {
  Edit as Pencil,
  Trash,
  Copy,
  Download,
  Upload,
  Move
} from "react-feather";

import { Action } from "../../types";

import { Content, Toolbar, Groups, InfoBar, Table } from "../../components";
import { Items } from "../../components/Table/types";

import {
  ProfileGroup,
  makeProfiles,
  removeProfiles,
  copyProfiles,
  copyToGroup,
  moveToGroup,
  addProfileGroup as addGroupAction,
  removeProfileGroup as removeGroupAction
} from '../../stores/Main/reducers/profiles';

import Create from "./Create";
import Edit from './Edit';
import { PROFILE_COLUMNS } from './constants';
import { filterProfiles } from "./filterProfiles";

const Profiles = () => {
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const groups = useSelector(makeProfiles);

  const [view, setView] = useState<string>("Basic");
  const [edit, setEdit] = useState<boolean>(false);
  const [create, setCreate] = useState<boolean>(false);
  const [createGroup, setCreateGroup] = useState<boolean>(false);
  const [group, setGroup] = useState<ProfileGroup>(groups.default);
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
      name: "Profile",
      actions: [
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
              dispatch(copyProfiles({ ids: [id], group: group.id }))

              // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              dispatch(copyProfiles({ ids: [...selected, id], group: group.id }))

              // use the selected rows
            } else {
              dispatch(copyProfiles({ ids: selected, group: group.id }))
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
              dispatch(removeProfiles({ ids: [id], group: group.id }))

              // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              dispatch(removeProfiles({ ids: [...selected, id], group: group.id }))

              // use the selected rows
            } else {
              dispatch(removeProfiles({ ids: selected, group: group.id }))
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
      title: "Edit Profiles",
      Icon: Pencil,
      onClick: useCallback(() => {
        if (!selected.length) {
          return;
        }
        setEdit(true);
      }, [group.id, selected])
    },
    {
      type: "Primary",
      title: "Copy Profiles",
      Icon: Copy,
      onClick: useCallback(
        () => dispatch(copyProfiles({ ids: selected, group: group.id })),
        [group.id, selected]
      ),
    },
    {
      type: "Primary",
      title: "Remove Profiles",
      Icon: Trash,
      onClick: useCallback(
        () => dispatch(removeProfiles({ ids: selected, group: group.id })),
        [group.id, selected]
      ),
    },
    {
      type: "First",
      title: "Import Profiles",
      Icon: Download,
      onClick: useCallback(() => {}, []),
    },
    {
      type: "Second",
      title: "Export Profiles",
      Icon: Upload,
      onClick: useCallback(() => {}, []),
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
      <Fragment key="Profiles">
        <Toolbar
          title={t("Pages.Profiles")}
          label="Add Profile"
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
          list="profiles"
          view={view}
          setView={setView}
        />
        <Groups
          tabs={{
            titles: ["Add", "Remove"],
            contents: [
              <Groups.AddGroup groups={groups} create={addGroupAction} />,
              <Groups.RemoveGroup remove={removeGroupAction} groups={groups} />,
            ],
          }}
          title="Profile Groups"
          open={createGroup}
          setOpen={setCreateGroup}
        />
        {edit && (
          <Edit group={group} ids={selected} open={edit} setOpen={setEdit} />
        )}
        {create && (
          <Create group={group} open={create} setOpen={setCreate} />
        )}
        <Table
          columns={PROFILE_COLUMNS}
          addGroup={setCreateGroup}
          group={group}
          groups={groups}
          data={tableData}
          view={view}
          filter={filter}
          filterRow={filterProfiles}
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
              name: "Copy",
              onClick: useCallback(
                (id: string) =>
                  dispatch(copyProfiles({ ids: [id], group: group.id })),
                [group.id]
              ),
            },
            {
              name: "Remove",
              onClick: useCallback(
                (id: string) =>
                  dispatch(removeProfiles({ ids: [id], group: group.id })),
                [group.id]
              ),
            },
          ]}
        />
      </Fragment>
    </Content>
  );
};

export default Profiles;
