import React, { useCallback, useMemo } from "react";
import styled from "styled-components";
import { motion } from "framer-motion";
import { Edit2 as Edit, Copy, Square, Play, Trash, Zap } from "react-feather";

import Drag from "./drag";

import { IconComponentProps } from "./types";

export const maskInput = (number: number, useSpaces = false) => {
  const input = "•".repeat(number);

  if (useSpaces) {
    const matches = input.match(/.{1,4}/g);
    if (matches?.length) {
      return matches.join(" ");
    }

    return input;
  }

  return input;
};

export const renderArrayCell = (value: any[]) => {
  if (!value.length) {
    return "None";
  }

  if (value.length > 1) {
    return "Multiple";
  }

  return value[0].name;
};

export const IconComponent = ({
  id,
  state,
  name,
  onClick,
}: IconComponentProps) => {
  let Icon: any;
  let handler = useCallback(() => {}, []);

  switch (name) {
    default:
    case "Run":
      // TODO: Figure out why useCallback was not rendering properly
      if (state) {
        Icon = StopIcon;
        handler = () => onClick("stop", id);
      } else {
        Icon = StartIcon;
        handler = () => onClick("start", id);
      }
      break;
    case "Copy":
      Icon = CopyIcon;
      handler = () => onClick(id);
      break;
    case "Test":
      Icon = TestIcon;
      handler = () => onClick(id);
      break;
    case "Edit":
      Icon = EditIcon;
      handler = () => onClick(id);
      break;
    case "Remove":
      Icon = RemoveIcon;
      handler = () => onClick(id);
      break;
    case "Drag":
      Icon = DragIcon;
  }

  return useMemo(
    () => (
      <IconContainer
        onClick={handler}
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
      >
        <Icon />
      </IconContainer>
    ),
    [handler, state]
  );
};

const IconContainer = styled(motion.div)``;

const TestIcon = styled(Zap)`
  width: 16px;
  height: 16px;
  cursor: pointer;
  padding: 0 4px;
  stroke: ${({ theme }) => theme.colors.primary};
  fill: ${({ theme }) => theme.colors.secondary};
`;

const DragIcon = styled(Drag)`
	color: ${({ theme }) => theme.colors.primary};
	opacity: 0.5;
	height: 16px;
	width: 12px;
	0 4px 0 6px;
	cursor: grab !important;

	& > * {
		cursor: grab !important;
	}
`;

const StartIcon = styled(Play)`
  width: 16px;
  height: 16px;
  cursor: pointer;
  padding: 0 4px;
  stroke: ${({ theme }) => theme.colors.primary};
  fill: ${({ theme }) => theme.colors.secondary};
`;

const StopIcon = styled(Square)`
  width: 16px;
  height: 16px;
  cursor: pointer;
  padding: 0 4px;
  fill-opacity: 0.4;
  fill: ${({ theme }) => theme.colors.stop};
  stroke: ${({ theme }) => theme.colors.stop};
`;

const EditIcon = styled(Edit)`
  width: 16px;
  height: 16px;
  cursor: pointer;
  padding: 0 4px;
  fill: ${({ theme }) => theme.colors.edit};
  stroke: ${({ theme }) => theme.colors.placeholder};
`;

const CopyIcon = styled(Copy)`
  width: 16px;
  height: 16px;
  cursor: pointer;
  padding: 0 4px;
  stroke: ${({ theme }) => theme.colors.placeholder};
`;

const RemoveIcon = styled(Trash)`
  width: 16px;
  height: 16px;
  cursor: pointer;
  padding: 0 4px;
  fill: rgba(242, 110, 134, 0.4);
  stroke: ${({ theme }) => theme.colors.failed};
`;
