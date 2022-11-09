import React, { Fragment, useState, useCallback, useMemo } from "react";
import { motion } from "framer-motion";
import styled from "styled-components";

import Logo from "./Logo";
import Menu from "./Menu";
import Settings from "./Settings";
import User from "./User";
import Tick from "./Tick";

const container = {
  open: {
    maxWidth: 215,
    transition: {
      duration: 0.45,
    },
  },
  collapsed: {
    maxWidth: 80,
    transition: {
      duration: 0.55,
    },
  },
};

const Sidebar = () => {
  const [collapsed, setCollapsed] = useState(false);

  const handleCollapse = useCallback(() => {
    setCollapsed(!collapsed);
  }, [collapsed]);

  return useMemo(
    () => (
      <Fragment key="sidebar">
        <Container
          initial="open"
          animate={collapsed ? "collapsed" : "open"}
          variants={container}
        >
          <Logo />
          <Menu collapsed={collapsed} />
          <Settings collapsed={collapsed} />
          <User />
        </Container>
        {/* <Tick collapsed={collapsed} handleCollapse={handleCollapse} /> */}
      </Fragment>
    ),
    [collapsed]
  );
};

const Container = styled(motion.div)`
  display: flex;
  flex-direction: column;
  flex: 1 0 0;
  min-width: 0;
  background-color: ${({ theme }) => theme.colors.sidebar};
`;

export default Sidebar;
