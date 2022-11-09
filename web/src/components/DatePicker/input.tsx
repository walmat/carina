import React from "react";
import { Icon } from "react-feather";
import styled from "styled-components";
import NumberFormat from "react-number-format";

interface InputProps {
  id: string;
  name: string;
  value: string | number;
  focused?: boolean;
  disabled?: boolean;
  autoFocus?: boolean;
  onChange: (e: any) => void;
  onBlur?: (e: any) => void;
  onFocus?: (e: any) => void;
}

export const DateInput = ({
  id,
  name,
  focused,
  value,
  disabled = false,
  autoFocus = false,
  onChange,
  onBlur,
  onFocus,
}: InputProps) => {
  // type: date causes an error on react-number-format
  return (
    <Container>
      <Input
        format="##/##/####"
        mask={["M", "M", "D", "D", "Y", "Y", "Y", "Y"]}
        {...{
          id,
          name,
          // type: "date",
          value,
          disabled,
          focused,
          onChange,
          onFocus,
          onBlur,
          autoFocus,
        }}
      />
    </Container>
  );
};

const Container = styled.div`
  position: relative;
`;

const Input = styled(NumberFormat)<{ focused: boolean }>`
  text-align: center;
  border-radius: 4px;
  font-size: 12px;
  width: 71px;
  font-family: Inter, sans-serif;
  font-weight: 400;
  height: 24px;
  outline: none;
  // flex: 1;
  padding: 0 8px;
  margin: 0;
  color: ${({ theme }) => theme.colors.h2};
  background-color: ${({ theme }) => theme.colors.sidebar};
  border: 1px solid ${({ theme }) => theme.colors.border};
  ${({ theme, focused }) =>
    focused ? `border-color: ${theme.colors.h2};` : ""}

  &:focus {
    border-color: ${({ theme }) => theme.colors.h2};
  }

  &::placeholder {
    color: ${({ theme }) => theme.colors.placeholder};
  }
`;

export default DateInput;
