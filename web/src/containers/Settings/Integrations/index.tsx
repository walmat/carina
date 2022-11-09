import React, { Fragment } from "react";
import { useSelector } from "react-redux";
import { useTranslation } from "react-i18next";

import { Content, InfoBar, Toolbar } from "../../../components";
import { Action } from "../../../types";
import { Download, Upload } from "react-feather";
import { makeIntegrations } from "../../../stores/Main/reducers/integrations";
import Grid from "./Grid";

const Integrations = () => {
  const { t } = useTranslation();
  const integrations = useSelector(makeIntegrations);

  const actions: Action[] = [
    {
      type: "First",
      title: "Import Integrations",
      Icon: Download,
      onClick: () => {},
    },
    {
      type: "Second",
      title: "Export Integrations",
      Icon: Upload,
      onClick: () => {},
    },
  ];
  return (
    <Content>
      <Fragment key="Integrations">
        <Toolbar simple title={t("Pages.Integrations")} />
        <InfoBar
          actions={actions}
          group={integrations}
          showView={false}
          list="integrations"
        />
        <Grid integrations={Object.values(integrations)} />
      </Fragment>
    </Content>
  );
};

export default Integrations;
