import React from "react";
import { motion } from "framer-motion";
import styled from "styled-components";

import { Router } from "../../animations";

interface ContentProps {
  children: React.ReactElement;
}

const Content = ({ children }: ContentProps) => {
  return (
    <Container
      initial="exit"
      animate="enter"
      exit="exit"
      variants={Router.variants}
    >
      {children}
    </Container>
  );
};

const Container = styled(motion.div)`
  position: relative;
  display: flex;
  flex-direction: column;
  flex: 1 0 0;
  min-width: 0;
  margin: 32px;
`;

export default Content;
