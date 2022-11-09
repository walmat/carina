import styled from "styled-components";
import { motion, useCycle } from "framer-motion";

import { Typography } from "../../../elements";

import { Toggle } from "./toggle";
import Sections from "./items";

type Props = {
  collapsed: boolean;
};

const root = {
  open: {
    marginLeft: 32,
    marginRight: 0,
    transition: {
      duration: 0.45,
    },
  },
  collapsed: {
    marginLeft: 16,
    marginRight: 16,
    transition: {
      duration: 0.5,
    },
  },
};

const expand = {
  open: {
    x: 0,
    transition: {
      delay: 0.05,
      duration: 0.45,
    },
  },
  collapsed: {
    x: -14,
    transition: {
      delay: 0.05,
      duration: 0.35,
    },
  },
};

const heading = {
  open: {
    opacity: 1,
    transition: {
      delay: 0.1,
      duration: 0.45,
    },
  },
  collapsed: {
    opacity: 0,
    transition: {
      duration: 0.3,
    },
  },
};

const Settings = ({ collapsed }: Props) => {
  const [isOpen, toggleOpen] = useCycle(false, true);

  return (
    <Container variants={root}>
      <Flex collapsed={collapsed} isOpen={isOpen} onClick={() => toggleOpen()}>
        <Heading variants={heading}>
          <Header>Settings</Header>
        </Heading>
        <Expand variants={expand}>
          <Toggle isOpen={isOpen} />
        </Expand>
      </Flex>

      <Sections
        key="sections"
        collapsed={collapsed}
        isOpen={isOpen}
      />
    </Container>
  );
};

const Container = styled(motion.div)`
  display: flex;
  flex-direction: column;
  margin-left: 32px;
  margin-right: 0;
  flex: 1;
`;

const Flex = styled.div<{ collapsed: boolean; isOpen: boolean }>`
  display: flex;
  align-items: center;
  cursor: pointer;
  margin: 8px 0;
  padding: 0px 16px;
  border-radius: ${({ collapsed }) =>
    collapsed ? "4px 4px 4px 4px" : "4px 0 0 4px"};
  background-color: ${({ theme, isOpen }) =>
    isOpen ? `${theme.colors.lightHue}, 0.05)` : "transparent"};
  cursor: pointer;

  &:hover {
    background-color: ${({ theme }) => `${theme.colors.lightHue}, 0.05)`};

    * {
      cursor: pointer;
    }
  }
`;

const Heading = styled(motion.div)`
  display: flex;
  align-items: center;
  width: 100%;
`;

const Header = styled(Typography.H2)`
  font-size: 12px;
  z-index: 500;
  font-weight: 500;
  margin-right: auto;
  display: flex;
  margin: 8px 6px;
  flex-direction: column;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const Expand = styled(motion.div)`
  display: flex;
  margin-top: 2px;
`;

export default Settings;
