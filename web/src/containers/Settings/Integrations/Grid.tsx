import React from "react";
import styled from "styled-components";

import { Integration } from '../../../stores/Main/reducers/integrations';
import Card from './Card';

interface Props {
  integrations: Integration[];
}

const Grid = ({ integrations }: Props) => {
  return (
    <Container>
      {integrations.map((integration) => (
        <Card integration={integration} />
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
    grid-template-columns: repeat(4, 1fr);
  }

  @media (min-height: 900px) {
    grid-template-rows: repeat(4, 1fr) !important;
  }

  @media (min-height: 1200px) {
    grid-template-rows: repeat(5, 1fr) !important;
  }
`;

export default Grid;
