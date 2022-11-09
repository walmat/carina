import React from "react";
import styled from "styled-components";

import { Webhook } from '../../../stores/Main/reducers/webhooks';
import Card from './Card';

interface Props {
  webhooks: Webhook[];
}

const Grid = ({ webhooks }: Props) => {
  return (
    <Container>
      {webhooks.map((webhook) => (
        <Card webhook={webhook} />
      ))}
    </Container>
  );
};

const Container = styled.div`
  display: grid;
  margin-top: 16px;
  gap: 16px;
  height: 100%;
  overflow: auto;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(3, 1fr);

  &::-webkit-scrollbar {
    display: none;
  }

  @media (min-width: 1400px) {
    grid-template-columns: repeat(3, 1fr);
  }

  @media (min-height: 900px) {
    grid-template-rows: repeat(3, 1fr) !important;
  }

  @media (min-height: 1200px) {
    grid-template-rows: repeat(3, 1fr) !important;
  }
`;

export default Grid;
