import React from "react";
import { useSelector } from "react-redux";
import moment from "moment";
import styled from "styled-components";
import { ChevronRight } from 'react-feather';

import { Typography } from '../../../../elements';
import { makeUser } from "../../../../stores/Main/reducers/user";

const getExpirationDate = (createdAt: string) => {
  const day = moment().date();
  const renews = moment(createdAt).date();

  return moment()
    .add(renews - day, "days")
    .format("MM/DD/YYYY");
};

const Expiration = () => {
  const { type, createdAt } = useSelector(makeUser);

  if (type === "Renewal") {
    const expiration = createdAt ? getExpirationDate(createdAt) : 'Unknown';

    return (
      <Container>
        <ButtonContainer>
          {type}
          <Separator />
          <ExpirationDate>Expires {expiration}</ExpirationDate>
        </ButtonContainer>
        <LearnMore>
          Learn more
          <Right />
        </LearnMore>

      </Container>
    );
  }

  return (
    <Container>
      <ButtonContainer>
        {type}
        <Separator />
        <ExpirationDate>No expiration</ExpirationDate>
      </ButtonContainer>
      <LearnMore>
        Learn more
        <Right />
      </LearnMore>
    </Container>
  );
};

const Container = styled.div`
  display: inline-flex;
`;

const ButtonContainer = styled.div`
	border-radius: 4px;
  margin: auto 0 0 0;
  display: flex;
  padding: 8px 12px;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.h2};
  background-color: ${({ theme }) => theme.colors.background};
  font-size: 12px;
  font-weight: 500;
`;

const LearnMore = styled(Typography.Paragraph)`
  margin: auto 8px 8px 8px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: 400;
  font-size: 12px;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.h2};
`;

const Right = styled(ChevronRight)`
  height: 14px;
  cursor: pointer;
`;

const Separator = styled.div`
  width: 2px;
  height: 8px;
  cursor: pointer;
  display: inline-flex;
  border-radius: 8px;
  margin: 3px 8px;
  background-color: ${({ theme }) => theme.colors.paragraph};
`;

const ExpirationDate = styled.span`
  display: inline-flex;
  cursor: pointer;
  font-weight: 500;
  font-size: 12px;
`;

export default Expiration;
