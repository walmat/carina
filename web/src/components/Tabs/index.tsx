import React, { useState } from "react";
import styled from "styled-components";
import { Tab, Tabs, TabList, TabPanel } from "react-tabs";

import { Typography } from "../../elements";

interface TabsProps {
  titles: string[];
  contents: any[];
}

const TabsPrimitive = ({ titles, contents }: TabsProps) => {
  const [tabIndex, setTabIndex] = useState(0);

  return (
    <Container tabIndex={tabIndex} onSelect={(index) => setTabIndex(index)}>
      <TabContainer>
        {titles.map((title, index) => (
          <TabTitleContainer>
            <TabTitle active={tabIndex === index}>{title}</TabTitle>
          </TabTitleContainer>
        ))}
      </TabContainer>

      {contents.map((content) => (
        <TabBody>{content}</TabBody>
      ))}
    </Container>
  );
};

const Container = styled(Tabs)`
  display: flex;
  flex-direction: column;
  flex: 1;
`;

const TabContainer = styled(TabList)`
  display: flex;
  height: 20px;
  list-style: none;
  margin: 0 0 16px 0;
  padding: 0;
  border-bottom: 1px solid ${({ theme }) => theme.colors.separator};
`;

const TabTitleContainer = styled(Tab)`
  display: flex;
  flex-direction: column;
`;

const TabTitle = styled(Typography.H4)<{ active: boolean }>`
  font-size: 14px;
  font-weight: ${({ active }) => (active ? 500 : 400)};
  margin: 0 8px;
  color: ${({ theme }) => theme.colors.paragraph};
  opacity: ${({ active }) => (active ? 1 : 0.4)};
  cursor: ${({ active }) => (active ? "default" : "pointer")};
`;

const TabBody = styled(TabPanel)`
  display: flex;
  flex-direction: column;
  flex: 1;
`;

export default TabsPrimitive;
