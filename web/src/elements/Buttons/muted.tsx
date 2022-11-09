import React from "react";
import styled from "styled-components";
import { motion } from "framer-motion";

type ButtonProps = {
  width: number | string;
  height: number | string;
  text: string;
  onClick: any;
};

const MutedButton = ({ width, height, text, onClick }: ButtonProps) => {
  return (
    <Button
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
      width={width}
      height={height}
      onClick={onClick}
    >
      {text}
    </Button>
  );
};

const Button = styled(motion.button)<{
  width: number | string;
  height: number | string;
}>`
  width: ${({ width }) => (typeof width === "string" ? width : `${width}px`)};
  height: ${({ height }) =>
    typeof height === "string" ? height : `${height}px`};
  border-radius: 4px;
  border: none;
  background-color: #fff;
  opacity: 0.6;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.primary};
  font-size: 14px;
  font-weight: 400;
`;

export default MutedButton;
