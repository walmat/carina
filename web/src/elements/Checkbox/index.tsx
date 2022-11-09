import React, { Dispatch, SetStateAction, useCallback } from "react";
import styled from "styled-components";
import { motion } from "framer-motion";

const tickVariants = {
  pressed: (isChecked: boolean) => ({
    borderWidth: isChecked ? `1px` : "1.5px",
  }),
  checked: { pathLength: 1, opacity: 1 },
  unchecked: { pathLength: 0, opacity: 0 },
};

interface CheckboxProps {
  isChecked: boolean;
  setIsChecked: Dispatch<SetStateAction<boolean>>;
}

const Checkbox = ({ isChecked, setIsChecked }: CheckboxProps) => {
  const onClickHandler = useCallback(() => {
    return setIsChecked(!isChecked);
  }, [setIsChecked, isChecked]);

  return (
    <Container
      isChecked={isChecked}
      initial={isChecked ? "checked" : "unchecked"}
      animate={isChecked ? "checked" : "unchecked"}
      onClick={onClickHandler}
    >
      <motion.path
        d="M 4 7.61 L 6.86 10.4 L 11 4"
        fill="transparent"
        strokeWidth="2"
        stroke="#fff"
        strokeLinecap="round"
        strokeLinejoin="round"
        variants={tickVariants}
        custom={isChecked}
      />
    </Container>
  );
};

const Container = styled(motion.svg)<{ isChecked: boolean }>`
  margin: auto 0;
  height: 16px;
  width: 16px;
  border-width: 1px;
  border-style: solid;
  box-sizing: border-box;
  opacity: ${({ isChecked }) => (isChecked ? 1 : 0.6)};
  border-color: ${({ theme }) => theme.colors.primary};
  border-radius: 4px;
  background-color: ${({ theme, isChecked }) =>
    isChecked ? theme.colors.primary : theme.colors.secondary};

  ${({ isChecked }) =>
    !isChecked
      ? `
    border-style: solid;
  `
      : ""}

  &:hover {
    border-width: ${({ isChecked }) => (isChecked ? 1 : 1.5)}px;
  }
`;

export default Checkbox;
