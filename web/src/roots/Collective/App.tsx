import React, { useMemo } from "react";
import styled from "styled-components";

import { Actions, Dragbar } from "../../components";
import Routes from "./Routes";

import "../../i18n/config";

const TopBar = () => {
	return (
		<Float>
			<Flex margin="0 0 0 auto">
				<Actions useHide useMargin={false} simple />
			</Flex>
		</Float>
	);
};

const App = () => {
	return useMemo(
		() => (
			<Container>
				<Dragbar />
				<TopBar />
				<Routes />
			</Container>
		),
		[]
	);
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1 0 0;
  min-width: 0;
  background-color: ${({ theme }) => theme.colors.background};
`;

const Float = styled.div`
  position: absolute;
  display: flex;
  top: 16px;
  left: 32px;
  width: calc(100% - 40px);
  z-index: 999;
`;

const Flex = styled.div<{ margin: string; }>`
  display: flex;
  flex-direction: column;
  margin: ${({ margin }) => margin};
`;

export default App;