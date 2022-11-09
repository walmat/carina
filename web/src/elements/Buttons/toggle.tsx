import React from "react";
import styled from "styled-components";
import { motion } from "framer-motion";

interface ToggleButtonProps {
	height: string | number;
	margin?: string;
	active: boolean;
	onClick: any;
	children: any;
}

const ToggleButton = ({
  height,
  margin,
  active,
  onClick,
  children,
 }: ToggleButtonProps) => {
	return (
		<Button
			whileHover={{ scale: 1.05 }}
			whileTap={{ scale: 0.95 }}
			{...{ active, margin, height, onClick }}
		>
			{children}
		</Button>
	);
};

const Button = styled(motion.button)<{ height: number | string; active: boolean; margin?: string; }>`
  height: ${({ height }) => typeof height === "string" ? height : `${height}px`};
  border-radius: 2px;
  ${({ margin }) => margin ? `margin: ${margin};` : ''}
  border: none;
  cursor: pointer;
  padding: 8px 16px;
  background-color: ${({ theme, active }) => active ? theme.colors.primary : theme.colors.background};
  color: ${({ theme, active }) => active ? '#fff' : theme.colors.paragraph};
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 400;
`;

export default ToggleButton;
