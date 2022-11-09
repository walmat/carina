import React from "react";
import { useSelector, useDispatch } from "react-redux";
import styled from "styled-components";

import { Select, IndicatorSeparator, Typography } from "../../../../elements";

import {
  editLanguage,
  editCurrency,
  makeSettings
} from '../../../../stores/Main/reducers/settings';

const languages = [
  {
    label: "English (US)",
    value: "en-US",
  },
];

const currencies = [
  {
    label: "USD ($)",
    value: "usd",
  },
  {
    label: "CAD ($)",
    value: "cad",
  },
  {
    label: "EUR (€)",
    value: "eur",
  },
  {
    label: "GBP (£)",
    value: "gbp",
  },
  {
    label: "JPY (¥)",
    value: "jpy",
  },
  {
    label: "AUD ($)",
    value: "jpy",
  },
];

const Preferences = () => {
  const dispatch = useDispatch();

  const { language, currency } = useSelector(makeSettings);

  return (
    <Container>
      <Heading>Preferences</Heading>
      <Row>
        <Col>
          <Title>Language</Title>
          <SubText>Choose your language.</SubText>
        </Col>
        <ColFill>
          <SelectContainer>
            <Select
              name="language"
              isClearable={false}
              placeholder="English (US)"
              components={{ IndicatorSeparator }}
              value={language}
              options={languages}
              onChange={(e: any) => dispatch(editLanguage(e))}
            />
          </SelectContainer>
        </ColFill>
      </Row>
      <Border />
      <Row>
        <Col>
          <Title>Currency</Title>
          <SubText>Choose your currency.</SubText>
        </Col>
        <ColFill>
          <SelectContainer>
            <Select
              name="currency"
              isClearable={false}
              placeholder="USD ($)"
              components={{ IndicatorSeparator }}
              value={currency}
              options={currencies}
              onChange={(e: any) => dispatch(editCurrency(e))}
            />
          </SelectContainer>
        </ColFill>
      </Row>
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  background-color: ${({ theme }) => theme.colors.sidebar};
  margin: 0 0 0 8px;
  flex-basis: 55%;
  padding: 16px;
  border-radius: 8px;
`;

const Heading = styled.div`
  font-size: 18px;
  margin-bottom: 8px;
  font-weight: 700;
  color: ${({ theme }) => theme.colors.subHeading};
`;

const Row = styled.div`
  display: flex;
  align-items: center;
  flex: 1;
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

const Border = styled.div`
  display: flex;
  margin: 8px 0;
  padding: 0.5px 0;
  background-color: ${({ theme }) => theme.colors.separator};
`;

const Title = styled(Typography.H2)`
  font-size: 14px;
  margin: 0 0 4px 0;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.h2};
`;

const SubText = styled(Typography.Paragraph)`
  font-size: 12px;
  margin: 0;
  font-weight: 400;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const SelectContainer = styled.div`
  margin-left: auto;
  width: 15vw;
  max-width: 15vw;
`;

export default Preferences;
