import React from "react";
import { motion } from "framer-motion";
import styled from "styled-components";

import Sections from "./sections";

type Props = {
  collapsed: boolean;
};

const root = {
  open: {
    maxWidth: 215,
    height: "auto",
    marginLeft: 32,
    marginRight: 0,
    transition: {
      duration: 0.45,
    },
  },
  collapsed: {
    maxWidth: 80,
    height: "auto",
    marginLeft: 16,
    marginRight: 16,
    transition: {
      duration: 0.45,
    },
  },
};

const Menu = ({ collapsed }: Props) => {
  return (
    <Container variants={root}>
      <Sections collapsed={collapsed} />
    </Container>
  );
};

const Container = styled(motion.div)`
  height: 100%;
  display: flex;
  flex-direction: column;
`;

export default Menu;
