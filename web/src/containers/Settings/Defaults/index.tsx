import React, {Fragment, useCallback, useMemo, useState} from "react";
import { useDispatch, useSelector } from "react-redux";
import { useTranslation } from "react-i18next";
import { Copy, Download, Edit, Trash, Upload } from "react-feather";

import { Content, InfoBar, Table, Toolbar } from "../../../components";

import {
  Default,
  copyDefault,
  makeDefaults,
  removeDefault
} from '../../../stores/Main/reducers/defaults';

import Create from "./Create";

import { Items } from "../../../components/Table/types";
import { Action } from "../../../types";
import { DEFAULTS_COLUMNS, extractAllStores } from './constants';
import { filterDefaults } from './filterDefaults';

const DefaultsComponent = () => {
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const groups = useSelector(makeDefaults);

  const [view, setView] = useState<string>("Basic");
  const [create, setCreate] = useState<boolean>(false);
  const [group, setGroup] = useState<Default | null>(null);
  const [filter, setFilter] = useState<string>("");
  const [selected, setSelected] = useState<string[]>([]);

  const tableData = useMemo(() => {
    if (group) {
      if (group.id === "all") {
        return Object.values(groups).map(({ byId }) =>
          Object.values(byId).reduce((acc, curr) => ({ ...acc, ...curr }), {})
        );
      }
      return Object.values(group.byId);
    }

    return Object.values(groups).map(({ byId }) =>
      Object.values(byId).reduce((acc, curr) => ({ ...acc, ...curr }), {})
    );
  }, [groups, group]);

  const items: Items[] = [
    {
      name: "Default",
      actions: [
        {
          name: "Edit",
          Icon: Edit,
          shortcut: ["⌘", "E"],
          onClick: () => {},
        },
        {
          name: "Copy",
          Icon: Copy,
          shortcut: ["⌘", "C"],
          onClick: useCallback((id: string) => {
            // if no rows are selected, just use the row we right clicked on
            if (!selected.length) {
              dispatch(copyDefault({ ids: [id] }))

              // if the selected rows don't contain the row we clicked on, patch it in
            } else if (!selected.find(_id => _id === id)) {
              dispatch(copyDefault({ ids: [...selected, id] }))

              // use the selected rows
            } else {
              dispatch(copyDefault({ ids: selected }))
            }
          }, [selected])
        },
        {
          name: "Remove",
          Icon: Trash,
          shortcut: ["⌘", "D"],
          onClick: () => {},
        },
      ],
    },
  ];

  const actions: Action[] = [
    {
      type: "Action",
      title: "Edit Defaults",
      Icon: Edit,
      onClick: useCallback(() => {}, []),
    },
    {
      type: "Primary",
      title: "Copy Defaults",
      Icon: Copy,
      onClick: useCallback(
        () => dispatch(copyDefault({ ids: selected })),
        [selected]
      ),
    },
    {
      type: "Primary",
      title: "Remove Defaults",
      Icon: Trash,
      onClick: useCallback(
        () => dispatch(removeDefault({ ids: selected })),
        [selected]
      ),
    },
    {
      type: "First",
      title: "Import Defaults",
      Icon: Download,
      onClick: useCallback(() => {}, []),
    },
    {
      type: "Second",
      title: "Export Defaults",
      Icon: Upload,
      onClick: useCallback(() => {}, []),
    },
  ];

  const onFilter = useCallback((e) => {
    setFilter(e.target.value);
  }, []);

  return (
    <Content>
      <Fragment key="Defaults">
        <Toolbar
          showCreateGroup={false}
          title={t("Pages.Defaults")}
          label="Add Default"
          onCreate={() => setCreate(true)}
          filter={filter}
          onFilter={onFilter}
          groups={{
            all: { id: "all", name: "All Stores", byId: {} },
            ...groups,
          }}
          group={group ? group : { name: "All Stores" }}
          setGroup={setGroup}
        />
        <InfoBar
          hasGroups
          showView={false}
          actions={actions}
          group={
            !group || group.id === "all" ? extractAllStores(groups) : group
          }
          list="defaults"
          view={view}
          setView={setView}
        />
        <Create open={create} setOpen={setCreate} />
        <Table
          columns={DEFAULTS_COLUMNS}
          group={group ? group : {}}
          groups={groups}
          data={tableData}
          view={view}
          filter={filter}
          filterRow={filterDefaults}
          setFilter={setFilter}
          selected={selected}
          setSelected={setSelected}
          items={items}
          actions={[
            {
              name: "Edit",
              onClick: useCallback(() => {}, []),
            },
            {
              name: "Remove",
              onClick: useCallback(
                (id: string) => dispatch(removeDefault({ ids: [id] })),
                []
              ),
            },
          ]}
        />
      </Fragment>
    </Content>
  );
};

export default DefaultsComponent;
