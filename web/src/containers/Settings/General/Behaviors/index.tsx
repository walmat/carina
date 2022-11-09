import React from "react";
import { useSelector, useDispatch } from "react-redux";
import styled from "styled-components";

import { Toggle, Typography } from "../../../../elements";

import {
  toggle,
  stash,
  BEHAVIORS,
  makeSettings
} from '../../../../stores/Main/reducers/settings';
import { makeTheme } from '../../../../stores/Main/reducers/theme';

import { noop } from "../../../../utils";

const Behaviors = () => {
  const dispatch = useDispatch();

  const theme = useSelector(makeTheme);
  const { behaviors } = useSelector(makeSettings);

  const handleToggle = (field: string) => {
    dispatch(toggle({ base: "behaviors", field }));
    dispatch(stash());
  };

  return (
    <Container>
      <Heading>Behaviors</Heading>
      <Row onClick={() => handleToggle(BEHAVIORS.NOTIFICATIONS)}>
        <ColFill>
          <Title>Software Notifications</Title>
          <SubText>Enables application push notifications.</SubText>
        </ColFill>
        <Col>
          <Toggle on={behaviors.notifications} theme={theme} check={noop} />
        </Col>
      </Row>
      <BottomBorder />
      <Row onClick={() => handleToggle(BEHAVIORS.RETRY_CHECKOUTS)}>
        <ColFill>
          <Title>Retry Checkouts</Title>
          <SubText>Retries checkout upon decline or other error.</SubText>
        </ColFill>
        <Col>
          <Toggle on={behaviors.retryCheckouts} theme={theme} check={noop} />
        </Col>
      </Row>
      <BottomBorder />
      <Row onClick={() => handleToggle(BEHAVIORS.MONITOR_POOLS)}>
        <ColFill>
          <Title>Monitoring Pooling</Title>
          <SubText>Enabling this helps save data when monitoring.</SubText>
        </ColFill>
        <Col>
          <Toggle on={behaviors.monitorPooling} theme={theme} check={noop} />
        </Col>
      </Row>
      <BottomBorder />
      <Row onClick={() => handleToggle(BEHAVIORS.AUTO_CLICK)}>
        <ColFill>
          <Title>Auto-Click Harvester</Title>
          <SubText>Automatically clicks checkbox (recommended).</SubText>
        </ColFill>
        <Col>
          <Toggle on={behaviors.autoClick} theme={theme} check={noop} />
        </Col>
      </Row>
      <BottomBorder />
      <Row onClick={() => handleToggle(BEHAVIORS.AUTO_LAUNCH)}>
        <ColFill>
          <Title>Auto-Launch Collective</Title>
          <SubText>Launches collective when needed (recommended).</SubText>
        </ColFill>
        <Col>
          <Toggle on={behaviors.autoLaunch} theme={theme} check={noop} />
        </Col>
      </Row>
    </Container>
  );
};

const Container = styled.div`
  padding: 16px;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  background-color: ${({ theme }) => theme.colors.sidebar};
  margin: 0 8px 0 0;
  flex-basis: 45%;
`;

const Heading = styled.div`
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 8px;
  color: ${({ theme }) => theme.colors.subHeading};
`;

const Row = styled.div`
  display: flex;
  flex: 1;
  align-items: center;
`;

const Col = styled.div`
  display: flex;
  flex-direction: column;
`;

const ColFill = styled.div`
  display: flex;
  flex: 1;
  flex-direction: column;
`;

const Title = styled(Typography.H2)`
  font-size: 14px;
  margin: 0 0 4px 0;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.subHeading};
`;

const SubText = styled(Typography.Paragraph)`
  font-size: 12px;
  margin: 0;
  font-weight: 400;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const BottomBorder = styled.div`
  display: flex;
  margin: 10px 0 8px 0;
  padding: 0.5px; 0;
  background-color: ${({ theme }) => theme.colors.separator};
`;

export default Behaviors;
