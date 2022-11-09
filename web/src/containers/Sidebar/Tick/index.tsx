import React from "react";
import { motion } from "framer-motion";
import styled from "styled-components";

import { MoreVertical } from "react-feather";

type Props = {
  collapsed: boolean;
  handleCollapse: () => void;
};

const root = {
  open: {
    left: 215,
    transition: {
      duration: 0.45,
    },
  },
  collapsed: {
    left: 80,
    transition: {
      duration: 0.55,
    },
  },
};

const Tick = ({ collapsed, handleCollapse }: Props) => {
  return (
    <Container
      initial="open"
      whileHover={{ scale: 1.2 }}
      whileTap={{ scale: 0.95, x: -2.5 }}
      animate={collapsed ? "collapsed" : "open"}
      variants={root}
      onClick={handleCollapse}
    >
      <Collapse />
    </Container>
  );
};

const Container = styled(motion.div)`
  position: absolute;
  bottom: 16px;
  background-color: ${({ theme }) => theme.colors.sidebar};
  height: 48px;
  width: 16px;
  cursor: pointer;
  border-top-right-radius: 8px;
  border-bottom-right-radius: 8px;
`;

const Collapse = styled(MoreVertical)`
  width: 16px;
  margin: auto;
  height: 100%;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.paragraph};
  align-items: center;
`;

export default Tick;
