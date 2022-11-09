import React from "react";
import { useSelector, useDispatch } from "react-redux";
import styled from "styled-components";

import { Typography, Select, IndicatorSeparator } from "../../../../elements";

import { makeTheme, setTheme } from '../../../../stores/Main/reducers/theme';

import { themes } from '../../../../constants';

const Theme = () => {
  const dispatch = useDispatch();
  const theme = useSelector(makeTheme);

  const handleTheme = ({ value }: any) => {
    if (theme === value) {
      return;
    }

    return dispatch(setTheme({ theme }));
  };

  return (
    <Container>
      <Heading>Theme</Heading>
      <RowBottom>
        <ColFill>
          <Title>Choose Theme</Title>
          <SelectContainer>
            <Select
              name="theme"
              isClearable={false}
              placeholder="Light (Default)"
              components={{ IndicatorSeparator }}
              value={themes[theme]}
              onChange={handleTheme}
              options={themes}
            />
          </SelectContainer>
        </ColFill>
      </RowBottom>
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  background-color: ${({ theme }) => theme.colors.sidebar};
  margin: 0 8px 0 0;
  flex-basis: 45%;
  padding: 16px;
  border-radius: 8px;
`;

const Heading = styled.div`
  font-size: 18px;
  margin-bottom: 8px;
  font-weight: 700;
  color: ${({ theme }) => theme.colors.subHeading};
`;

const RowBottom = styled.div`
  display: flex;
  flex: 1;
  align-items: flex-end;
`;

const ColFill = styled.div`
  display: flex;
  flex: 1;
  flex-direction: column;
`;

const Title = styled(Typography.H2)`
  font-size: 14px;
  margin: 0 0 8px 0;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const SelectContainer = styled.div`
  margin-right: auto;
  width: 15vw;
  max-width: 15vw;
`;

export default Theme;
