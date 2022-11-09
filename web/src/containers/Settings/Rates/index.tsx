import React, { Fragment, useState, useCallback } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useTranslation } from "react-i18next";

import { Download, Upload, Edit, Trash, Copy } from "react-feather";

import { Action } from "../../../types";
import { Items } from "../../../components/Table/types";
import { Content, InfoBar, Table, Toolbar } from "../../../components";

import { makeRates, removeRates } from '../../../stores/Main/reducers/rates';

import AddRates from "./Create";
import { filterRates } from './filterRates';
import { RATES_COLUMNS } from './constants';

const ShippingRates = () => {
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const groups = useSelector(makeRates);

  const [create, setCreate] = useState<boolean>(false);
  const [filter, setFilter] = useState<string>("");
  const [selected, setSelected] = useState<string[]>([]);

  const items: Items[] = [
    {
      name: "Rates",
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
          onClick: () => {},
        },
        {
          name: "Remove",
          Icon: Trash,
          shortcut: ["⌘", "D"],
          onClick: () => {},
        },
      ],
    }
  ];

  const actions: Action[] = [
    {
      type: "Action",
      title: "Edit Rates",
      Icon: Edit,
      onClick: () => {},
    },
    {
      type: "Primary",
      title: "Remove Rates",
      Icon: Trash,
      onClick: useCallback(
        () => dispatch(removeRates({ ids: selected })),
        [selected]
      ),
    },
    {
      type: "First",
      title: "Import Create",
      Icon: Download,
      onClick: () => {},
    },
    {
      type: "Second",
      title: "Export Create",
      Icon: Upload,
      onClick: () => {},
    },
  ];

  const onFilter = useCallback((e) => {
    setFilter(e.target.value);
  }, []);

  return (
    <Content>
      <Fragment key="Rates">
        <Toolbar
          title={t("Pages.Shipping Rates")}
          label="Add Rates"
          hasGroups={false}
          onCreate={() => setCreate(true)}
          filter={filter}
          onFilter={onFilter}
        />
        <InfoBar
          actions={actions}
          group={groups}
          showView={false}
          list="rates"
        />
        <AddRates open={create} setOpen={setCreate} />
        <Table
          columns={RATES_COLUMNS}
          group={null}
          groups={groups}
          data={Object.values(groups)}
          view="Basic"
          filter={filter}
          filterRow={filterRates}
          setFilter={setFilter}
          selected={selected}
          setSelected={setSelected}
          items={items}
          actions={[
            {
              name: "Edit",
              onClick: () => {},
            },
            {
              name: "Remove",
              onClick: useCallback(
                (id: string) => dispatch(removeRates({ ids: [id] })),
                []
              ),
            },
          ]}
        />
      </Fragment>
    </Content>
  );
};

export default ShippingRates;
