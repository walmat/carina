import React from "react";
import { useSelector } from "react-redux";
import styled from "styled-components";
import { Check } from 'react-feather';

import { makeUser } from "../../../../stores/Main/reducers/user";
import { Typography } from "../../../../elements";

const Activations = () => {
	const { instances, maxInstances } = useSelector(makeUser);
	return (
		<Container>
			<CheckIcon />
			<InstancesText>{instances} / {maxInstances}</InstancesText>
			 {maxInstances > 1 ? "Activations" : "Activation"}
		</Container>
	);
};

const Container = styled.div`
	border-radius: 4px;
  margin: 0 0 auto auto;
  display: flex;
  white-space: pre;
  padding: 8px 12px;
  cursor: pointer;
  justify-content: center;
  align-items: center;
  color: #fff;
  background-color: ${({ theme }) => theme.colors.failed};
  font-size: 12px;
  font-weight: 500;
`;

const InstancesText = styled(Typography.Paragraph)`
	font-size: 12px;
	margin: 0;
	font-weight: 500;
	margin-right: 8px;
  cursor: pointer;
`;

const CheckIcon = styled(Check)`
	width: 14px;
	height: auto;
	margin-right: 4px;
  cursor: pointer;
`;

export default Activations;
