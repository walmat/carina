import React from "react";
import styled from "styled-components";

interface CardProps {
  children: React.ReactElement | React.ReactElement[];
}

const Card = ({ children, ...props }: CardProps) => {
  return <Container {...props}>{children}</Container>;
};

const Container = styled.div<{ height?: number; width?: number }>`
  display: flex;
  flex-direction: column;
  border-radius: 8px;
  background-color: ${({ theme }) => theme.colors.sidebar};
  padding: 16px;
  ${({ height, width }) => `
		${height ? `height: ${height}px` : ""};
		${width ? `width: ${width}px` : ""};
	`}
`;

export default Card;
