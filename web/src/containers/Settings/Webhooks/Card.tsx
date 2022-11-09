import React, { useCallback, useRef, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import styled from "styled-components";
import { Edit2 as Edit, Trash, Check, X } from 'react-feather';
import { motion } from "framer-motion";

import { Card } from "../../../components";
import { Typography, Toggle } from "../../../elements";
import { AYCD, Slack, Discord } from "../../../icons";

import { makeTheme } from '../../../stores/Main/reducers/theme';
import { Webhook, enableWebhook, testWebhook, removeWebhook } from '../../../stores/Main/reducers/webhooks';

import EditComponent from './Edit';

type Props = {
  webhook: Webhook;
};

const getWebhookIcon = (type: string) => {
  switch (type) {
    default:
    case "discord":
      return <Discord />;
    case "slack":
      return <Slack />;
    case "aycd":
      return <AYCD />;
  }
};

const CardComponent = ({ webhook }: Props) => {
  let refs: any = useRef([]);
  const [open, setOpen] = useState(false);
  let containerRef: any = useRef(null);

  const theme = useSelector(makeTheme);
  const dispatch = useDispatch();

  const { id, name, type, active, declines, sensitivity, profiles } = webhook;

  const handleActivate = () => {
    const data = { ...webhook, active: !webhook.active };
    return dispatch(enableWebhook(data));
  };

  const handleRemove = useCallback(() => dispatch(removeWebhook(id)), [id]);

  return (
    <Card key={id}>
      <Flex height={32} ai="center">
        <Icon component={getWebhookIcon(type)} />
        <Title>{name}</Title>
        <EditIcon onClick={() => setOpen(true)} />
        <ToggleContainer>
          <Toggle on={active} size="small" check={handleActivate} theme={theme} />
        </ToggleContainer>
        <RemoveIcon onClick={handleRemove} />
      </Flex>
      <Flex mt="8px" ai="center">
        {declines ? <PositiveBtn><CheckMark /> Declines</PositiveBtn> : <NegativeBtn><Close /> Declines</NegativeBtn>}
        {sensitivity ? <PositiveBtn><CheckMark /> Sensitivity</PositiveBtn> : <NegativeBtn><Close /> Sensitivity</NegativeBtn>}
      </Flex>
      <Border />
      <Flex fd="column" mt="0">
        <Subtitle>Profile Groups</Subtitle>
        <ProfilesFlex ref={containerRef}>
          {profiles.map(({ name }) => {
            return <ProfileChip ref={ref => refs.current.push(ref)}>{name}</ProfileChip>;
          })}
        </ProfilesFlex>
      </Flex>
      <EditComponent open={open} webhook={webhook} setOpen={setOpen} />
    </Card>
  );
};

const Flex = styled.div<{
  ai?: string;
  jc?: string;
  height?: number;
  fd?: string;
  mt?: string;
  mb?: string;
  mr?: string;
  ml?: string;
}>`
  display: flex;
  ${({ ai }) => (ai ? `align-items: ${ai}` : "")};
  ${({ jc }) => (jc ? `justify-content: ${jc}` : "")};
  ${({ height }) => (height ? `height: ${height}px` : "")};
  ${({ fd }) => (fd ? `flex-direction: ${fd}` : "")};
  ${({ mt }) => (mt ? `margin-top: ${mt}` : "")};
  ${({ mb }) => (mb ? `margin-bottom: ${mb}` : "")};
  ${({ mr }) => (mr ? `margin-right: ${mr}` : "")};
  ${({ ml }) => (ml ? `margin-left: ${ml}` : "")};
`;

const Icon = styled(({ component, ...props }) =>
  React.cloneElement(component, props)
)`
  height: 32px;
  margin-right: 8px;
`;

const Title = styled(Typography.H2)`
  display: flex;
  line-height: 24px;
  max-width: 7vw;
  white-space: nowrap;
  text-overflow: hidden;
  overflow: hidden;
  color: ${({ theme }) => theme.colors.h2};
  text-transform: capitalize;
  font-size: 16px;
  font-weight: 700;
  margin: 0 8px;
`;

const EditIcon = styled(Edit)`
  height: 16px;
  width: 16px;
  stroke: ${({ theme }) => theme.colors.h3};
  fill-opacity: 0.4;
  fill: ${({ theme }) => theme.colors.h3};
  margin: 0 8px 0 0;
`;

const RemoveIcon = styled(Trash)`
  height: 16px;
  width: 16px;
  stroke: ${({ theme }) => theme.colors.failed};
  fill-opacity: 0.4;
  fill: ${({ theme }) => theme.colors.failed};
  margin: 0 0 0 auto;
`;

const PositiveBtn = styled.div`
  border-radius: 2.5px;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4px 8px;
  color: ${({ theme }) => theme.colors.h2};
  background-color: ${({ theme }) => theme.colors.background};
  font-size: 12px;
  margin: 0 8px 0 0;
`;

const NegativeBtn = styled.div`
  border-radius: 2.5px;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4px 8px;
  color: ${({ theme }) => theme.colors.h2};
  background-color: ${({ theme }) => theme.colors.background};
  font-size: 12px;
  margin: 0 8px 0 0;
`;

const CheckMark = styled(Check)`
  height: 14px;
  width: 14px;
  cursor: default;
  stroke: ${({ theme }) => theme.colors.primary};
  margin-right: 4px;
  
  & > * {
    cursor: default;
  }
`;

const Close = styled(X)`
  height: 14px;
  width: 14px;
  cursor: default;
  stroke: ${({ theme }) => theme.colors.failed};
  margin-right: 4px;
  
  & > * {
    cursor: default;
  }
`;

const ToggleContainer = styled.div`
  margin: 0;
`;

const Border = styled(motion.div)`
  margin: 16px 0;
  border-bottom: 1px solid ${({ theme }) => theme.colors.separator};
`;

const Subtitle = styled(Typography.H4)`
  color: ${({ theme }) => theme.colors.h3};
  font-size: 14px;
  font-weight: 500;
  margin: 0 0 8px 0;
`;

const ProfilesFlex = styled.div`
  margin: 0 -4px;
  width: auto;
  display: flex;
`;

const ProfileChip = styled.div`
  background-color: ${({ theme }) => theme.colors.primary};
  color: white;
  font-size: 12px;
  font-weight: 400;
  margin: 0 4px;
  border-radius: 4px;
  padding: 4px 8px;
`;

export default CardComponent;
