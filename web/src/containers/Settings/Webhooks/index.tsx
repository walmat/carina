import React, { Fragment, useState, useCallback } from "react";
import { useSelector } from "react-redux";
import { useTranslation } from "react-i18next";
import { Download, Upload } from "react-feather";

import { Action } from "../../../types";

import { Content, InfoBar, Toolbar } from "../../../components";

import { makeWebhooks } from '../../../stores/Main/reducers/webhooks';

import Grid from "./Grid";
import Create from "./Create";
import { filterWebhooks } from './filterWebhooks';

const Webhooks = () => {
  const { t } = useTranslation();
  const webhooks = useSelector(makeWebhooks);

  const [create, setCreate] = useState<boolean>(false);
  const [filter, setFilter] = useState<string>("");

  const actions: Action[] = [
    {
      type: "First",
      title: "Import Webhooks",
      Icon: Download,
      onClick: () => {},
    },
    {
      type: "Second",
      title: "Export Webhooks",
      Icon: Upload,
      onClick: () => {},
    },
  ];

  const onFilter = useCallback((e) => {
    setFilter(e.target.value);
  }, []);

  return (
    <Content>
      <Fragment key="Webhooks">
        <Toolbar
          title={t("Pages.Webhooks")}
          label="Add Webhook"
          hasGroups={false}
          onCreateGroup={() => {}}
          onCreate={() => setCreate(true)}
          filter={filter}
          onFilter={onFilter}
        />
        <InfoBar
          actions={actions}
          group={webhooks}
          showView={false}
          list="webhooks"
        />
        {create &&
          <Create open={create} setOpen={setCreate} />
        }
        <Grid webhooks={Object.values(webhooks).filter(wh => filterWebhooks(wh, filter))} />
      </Fragment>
    </Content>
  );
};

export default Webhooks;
