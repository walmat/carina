import { Dispatch, SetStateAction } from "react";
import { useSelector } from "react-redux";
import styled from "styled-components";
import { motion } from "framer-motion";

import { makeTheme } from '../../stores/Main/reducers/theme';

const circle = (size: string) => ({
  off: {
    x: size ==='small' ? '3px' : "4px",
    backgroundColor: "#616161",
  },
  on: {
    x: size ==='small' ? '18px' : "26px",
    backgroundColor: "#786EF2",
  },
});

const toggle = {
  off: (theme: number) => ({
    backgroundColor: theme === 0 ? "#F0F0F0" : "#202126",
  }),
  on: (theme: number) => ({
    backgroundColor: theme === 0 ? "#F2F1FE" : "rgba(120, 110, 242, 0.2)",
  }),
};

type ToggleProps = {
  on: boolean;
  size?: "small" | "default";
  check: Dispatch<SetStateAction<boolean>>;
};

const ToggleElement = ({ on, size = "default", check }: ToggleProps) => {
  const theme = useSelector(makeTheme);

  return (
    <Toggle
      size={size}
      onClick={() => check((prev) => !prev)}
      custom={theme}
      variants={toggle}
      initial={on ? "on" : "off"}
      animate={on ? "on" : "off"}
      transition={{
        backgroundColor: {
          duration: 0.01,
        },
      }}
    >
      <Circle
        size={size}
        variants={circle(size)}
        custom={theme}
        initial={on ? "on" : "off"}
        animate={on ? "on" : "off"}
        transition={{
          backgroundColor: {
            duration: 0.01,
          },
          x: { type: "spring", stiffness: 200, damping: 17.5 },
        }}
      />
    </Toggle>
  );
};

ToggleElement.defaultProps = {
  theme: 0,
};

const Toggle = styled(motion.div)<{ size: string; }>`
  width: ${({ size }) => size === 'small' ? '36px' : '48px'};
  height: ${({ size }) => size === 'small' ? '20px' : '24px'};
  border-radius: 16px;
  position: relative;
  cursor: pointer;
`;

const Circle = styled(motion.div)<{ size: string; }>`
  width: ${({ size }) => size === 'small' ? '14px' : '18px'};
  height: ${({ size }) => size === 'small' ? '14px' : '18px'};
  border-radius: 50%;
  position: absolute;
  cursor: pointer;
  top: ${({ size }) => size === 'small' ? '3px' : '3px'};
`;

export default ToggleElement;
