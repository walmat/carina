import React from "react";
import { motion } from "framer-motion";
import styled from "styled-components";

interface LoaderProps {
  width: number;
  height: number;
}

const spin = {
  loop: Infinity,
  duration: 1,
  ease: "linear",
};

export const Loader = ({ width, height }: LoaderProps) => {
  return (
    <Container width={width} height={height}>
      <Circle
        width={width}
        height={height}
        animate={{ rotate: 360 }}
        transition={spin}
      />
    </Container>
  );
};

const Container = styled(motion.div) <{ width: number; height: number; }>`
	position: relative;
	width: ${({ width }) => width}px;
	height: ${({ height }) => height}px;
	box-sizing: border-box;
`;

const Circle = styled(motion.span) <{ width: number; height: number; }>`
	display: block;
	width: ${({ width }) => width}px;
	height: ${({ height }) => height}px;
	border: 2.5px solid ${({ theme }) => theme.colors.expand};
	border-top: 2.5px solid ${({ theme }) => theme.colors.altLogin};
	border-radius: 50%;
	position: absolute;
	box-sizing: border-box;
	top: 0;
	left: 0;
`;

export default Loader;
