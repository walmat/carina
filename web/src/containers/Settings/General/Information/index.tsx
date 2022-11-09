import React from "react";
import styled from "styled-components";

import SignOff from "./SignOff";
import Version from "./version";
import Expiration from "./expiration";
import Activations from "./activations";
import Update from "./update";

const Information = () => {
  return (
    <Container>
      <Row>
        <Col>
          <Heading>Software Information</Heading>
          <Version />
        </Col>
        <Col>
          <Activations />
          <SignOff />
          <Update />
        </Col>
      </Row>
      <RowFill>
        <Expiration />
      </RowFill>
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

const Row = styled.div`
  display: flex;
`;

const RowFill = styled.div`
  display: flex;
  flex: 1;
`;

const Col = styled.div`
  display: flex;
  flex-direction: column;
`;

const ColFill = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;
`;

const Heading = styled.div`
  font-size: 18px;
  font-weight: 700;
  color: ${({ theme }) => theme.colors.subHeading};
`;

export default Information;
