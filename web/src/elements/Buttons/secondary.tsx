import React from "react";
import styled from "styled-components";
import { motion } from "framer-motion";
import { Typography } from "../Typography";

import { ButtonProps } from "./types";

const SecondaryButton = ({
  variant = "Button",
  text = "",
  command = "",
  width,
  height,
  onClick,
  disabled,
  children,
}: ButtonProps) => {
  if (typeof disabled != "boolean") {
    disabled = false;
  }

  if (variant === "Icon") {
    return (
      <Icon
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        width={width}
        height={height}
        onClick={onClick}
      >
        {children}
      </Icon>
    );
  }

  if (variant === "IconButton") {
    return (
      <IconButton
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        width={width}
        height={height}
        onClick={onClick}
      >
        {command ? (
          <CommandText>
            {text}
            <span>{command}</span>
          </CommandText>
        ) : (
          text || children
        )}
      </IconButton>
    );
  }

  return (
    <Button
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
      width={width}
      height={height}
      onClick={onClick}
      disabled={disabled}
    >
      {text}
    </Button>
  );
};

const Button = styled(motion.button)<{
  width: number | string;
  height: number | string;
}>`
  & {
    width: ${({ width }) => (typeof width === "string" ? width : `${width}px`)};
    height: ${({ height }) =>
      typeof height === "string" ? height : `${height}px`};
    border-radius: 4px;
    border: none;
    cursor: pointer;
    background-color: ${({ theme }) => theme.colors.failed};
    color: #ffffff;
    font-size: 14px;
    font-weight: 700;
    margin: 0;
  }
  &:disabled {
    width: ${({ width }) => (typeof width === "string" ? width : `${width}px`)};
    height: ${({ height }) =>
      typeof height === "string" ? height : `${height}px`};
    border-radius: 4px;
    border: none;
    cursor: pointer;
    background-color: #565656;
    color: #ffffff;
    font-size: 14px;
    font-weight: 700;
    margin: 0;
  }
`;

const IconButton = styled(motion.button)<{
  width?: number | string;
  height: number | string;
}>`
  height: ${({ height }) =>
    typeof height === "string" ? height : `${height}px`};
  border-radius: 4px;
  border: none;
  cursor: pointer;
  padding: 8px 16px;
  background-color: ${({ theme }) => theme.colors.failed};
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;

  ${({ width }) =>
    typeof width != "undefined"
      ? `
    width: ${width}px;
    max-width: ${width}px;
    min-width: ${width}px;    
  `
      : ""}
`;

const Icon = styled(motion.button)<{
  width: number | string;
  height: number | string;
}>`
  width: ${({ width }) => (typeof width === "string" ? width : `${width}px`)};
  height: ${({ height }) =>
    typeof height === "string" ? height : `${height}px`};
  border-radius: 4px;
  border: none;
  cursor: pointer;
  background-color: ${({ theme }) => theme.colors.failed};
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
`;

const CommandText = styled(Typography.Paragraph)`
  color: #fff;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;

  & > span {
    cursor: pointer;
    display: inline-flex;
    margin-left: 8px;
    padding: 4px 8px;
    border-radius: 2px;
    color: #fff;
    background-color: ${({ theme }) => `${theme.colors.lightHue}, 0.15)`};
    font-size: 10px;
  }
`;

export default SecondaryButton;
