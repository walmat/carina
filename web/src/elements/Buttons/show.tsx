import React, { Dispatch, SetStateAction } from "react";
import { motion, useMotionValue, useTransform } from "framer-motion";
import styled from "styled-components";

interface ShowIconProps {
  show: boolean;
  setShow: Dispatch<SetStateAction<boolean>>;
}

export const ShowIcon = ({ show, setShow }: ShowIconProps) => {
  const duration = 0.4;
  const lineVariants = {
    pressed: (show: boolean) => ({ pathLength: show ? 0.85 : 0.05 }),
    checked: { pathLength: 1 },
    unchecked: { pathLength: 0 },
  };
  const pathLength = useMotionValue(0);
  const opacity = useTransform(pathLength, [0.05, 0.15], [0, 1]);

  return (
    <Container
      show={show}
      onClick={() => {
        setShow((prev) => !prev);
      }}
    >
      <Button
        initial={false}
        animate={show ? "checked" : "unchecked"}
        whileHover="hover"
        whileTap="pressed"
        transition={{ duration }}
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <motion.path
          d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"
          stroke="#786EF2"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        <motion.circle
          cx="12"
          cy="12"
          r="3"
          stroke="#786EF2"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        <motion.line
          x1="1"
          y1="1"
          x2="24"
          y2="24"
          stroke="#786EF2"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          variants={lineVariants}
          style={{ pathLength, opacity }}
          custom={show}
          transition={{ duration }}
        />
      </Button>
    </Container>
  );
};

const Container = styled.button<{ show: boolean }>`
  display: flex;
  justify-content: center;
  align-items: center;
  background: transparent;
  border: none;
  height: 24px;
  cursor: ${({ show }) => (show ? "default" : "pointer")};
  outline: none;
`;

const Button = styled(motion.svg)`
  color: ${({ theme }) => theme.colors.primary};
`;

export default ShowIcon;
