import React, { Dispatch, SetStateAction } from "react";
import { motion } from "framer-motion";
import styled from "styled-components";
import { Plus, Minus } from "react-feather";

interface CounterProps {
  count: number;
  setCount: Dispatch<SetStateAction<number>>;
}

const Counter = ({ count, setCount }: CounterProps) => (
  <Container>
    <Count
      id="count"
      name="count"
      type="tel"
      placeholder="1"
      onBlur={() => {
        if (!count) {
          setCount(1);
        }
      }}
      onChange={(e) => {
        const { value } = e.target;

        const number = Number(value);
        if (!Number.isNaN(number)) {
          return setCount(number);
        }

        return setCount(1);
      }}
      value={count}
    />
    <Icon
      onClick={() => {
        if (count === 1) {
          return;
        }

        return setCount((prev) => prev - 1);
      }}
      component={<Minus />}
    />
    <Icon component={<Plus />} onClick={() => setCount((prev) => prev + 1)} />
  </Container>
);

const Container = styled(motion.div)`
  display: flex;
  justify-content: center;
  align-items: center;
`;

const Count = styled.input`
  padding: 4px;
  border-radius: 4px;
  text-align: center;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: ${({ theme }) => theme.colors.secondary};
  color: ${({ theme }) => theme.colors.primary};
  height: 24px;
  min-width: 40px;
  max-width: 50px;
  font-size: 14px;
  font-weight: 500;
  margin-right: 8px;
`;

const Icon = styled(({ component, ...props }) =>
  React.cloneElement(component, props)
)`
  height: 16px;
  width: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  margin: 0 4px;
  color: ${({ theme }) => theme.colors.paragraph};
`;

export default Counter;
