import React, { useState } from "react";
import { ChevronLeft, ChevronRight } from "react-feather";
import styled from "styled-components";
import { AnimatePresence, motion } from "framer-motion";

const getTitle = (list: string) => {
  switch (list) {
    default:
      return list;
    case "quicktasks":
      return "Quick Tasks";
  }
};

type InfoBarProps = {
  info: any[];
  list: string;
};

const infoVariants = {
  open: {
    opacity: 1,
    x: 0,
    transition: {
      x: { type: "spring", stiffness: 250, damping: 20 },
    },
  },
  collapsed: {
    opacity: 0,
    x: -96,
    transition: {
      x: { type: "spring", stiffness: 250, damping: 20 },
    },
  },
};

const expandVariants = {
  open: {
    opacity: 1,
    x: 128,
    rotate: 180,
    transition: {
      rotate: { type: "linear", duration: 0.1 },
      x: { type: "spring", stiffness: 275, damping: 20 },
    },
  },
  collapsed: {
    opacity: 1,
    x: 0,
    transition: {
      x: { type: "spring", stiffness: 275, damping: 25 },
    },
  },
};

const InfoBar = ({ info, list }: InfoBarProps) => {
  const [collapsed, setCollapsed] = useState<boolean>(true);

  const useDrawer = list === "tasks";

  if (!useDrawer) {
    return (
      <Container useCursor={false}>
        <Number>{info.length}</Number>
        <Text>{getTitle(list)}</Text>
      </Container>
    );
  }

  return (
    <Container useCursor onClick={() => setCollapsed((prev) => !prev)}>
      <Number>{info.length}</Number>
      <Text>{list}</Text>
      <ExpandContainer
        key="expand"
        initial="collapsed"
        variants={expandVariants}
        animate={!collapsed ? "open" : "collapsed"}
        exit="collapsed"
      >
        <Collapse />
      </ExpandContainer>
      <AnimatePresence initial={false} exitBeforeEnter>
        {!collapsed && (
          <Info
            key="info-container"
            initial="collapsed"
            variants={infoVariants}
            animate={!collapsed ? "open" : "collapsed"}
            exit="collapsed"
          >
            <Number>{info.filter(({ state }) => state).length}</Number>
            <Text>Running</Text>
          </Info>
        )}
      </AnimatePresence>
    </Container>
  );
};

const Container = styled(motion.div)<{ useCursor: boolean }>`
  display: flex;
  align-items: center;
  cursor: ${({ useCursor }) => (useCursor ? "pointer" : "default")};

  & > * {
    cursor: ${({ useCursor }) => (useCursor ? "pointer" : "default")};
  }
`;

const Number = styled.div`
  color: ${({ theme }) => theme.colors.primary};
  padding: 6.5px 4px;
  background-color: ${({ theme }) => theme.colors.secondary};
  border-radius: 4px;
  min-width: 40px;
  max-height: 32px;
  font-size: 14px;
  font-weight: 700;
  display: flex;
  justify-content: center;
  align-items: center;
`;

const Text = styled.div`
  color: ${({ theme }) => theme.colors.h2};
  text-transform: capitalize;
  font-weight: 400;
  margin: 0 0 0 8px;
  font-size: 14px;
  display: flex;
  justify-content: center;
  align-items: center;
`;

const ExpandContainer = styled(motion.div)`
  width: 16px;
  height: 16px;
  display: flex;
  margin-left: 8px;
  justify-content: center;
  align-items: center;

  &:hover {
    opacity: 0.6;
  }
`;

const Collapse = styled(ChevronRight)`
  display: flex;
  color: ${({ theme }) => theme.colors.expand};
  justify-content: center;
  align-items: center;
  margin: auto 0;
  width: 14px;
  height: 14px;
`;

const Info = styled(motion.div)`
  display: flex;
`;

export default InfoBar;
