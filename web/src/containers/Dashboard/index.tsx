import React, { Fragment } from "react";
import styled from "styled-components";
import { useTranslation } from "react-i18next";

import { Content, Toolbar } from "../../components";

import Analytics from "./Analytics";

import Shipments from "./Shipments/Shipments";
import OrderHistory from "./OrderHistory/OrderHistory";

const Dashboard = () => {
  const { t } = useTranslation();

  return (
    <Content>
      <Fragment key="Dashboard">
        <Toolbar simple title={t("Pages.Dashboard")} />
        <Container>
          <Analytics />
          <InnerContainer>
            {/*<Shipments />*/}
            <OrderHistory />
          </InnerContainer>
        </Container>
      </Fragment>
    </Content>
  );
};

const Container = styled.div`
  gap: 16px;
  height: calc(100% - 52px);
  margin: 16px 0 0 0;
  display: grid;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  grid-template-rows: repeat(2, minmax(0, 1fr));
`;

const InnerContainer = styled.div`
  gap: 16px;
  margin: 0;
  display: grid;
`;

/*
  grid-template-columns: 0.75fr 1fr;
  grid-template-rows: repeat(1, minmax(0, 1fr));
 */

export default Dashboard;
