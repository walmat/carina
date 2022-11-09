import React from "react";
import styled from "styled-components";

import { Typography } from "../../elements";

const Empty = ({ label }: { label: string }) => {
  return (
    <Container>
      <Text>{label}</Text>
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
`;

const Text = styled(Typography.H2)`
  color: #919191;
  font-size: 1.65vw;
`;

export default Empty;
