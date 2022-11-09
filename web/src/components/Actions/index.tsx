import React from "react";
import { motion } from "framer-motion";
import styled from "styled-components";
import { X as Close, Minus } from "react-feather";

import { minimize, hide, close } from "../../utils";

interface ActionProps {
  useMargin?: boolean;
  useHide?: boolean;
  simple: boolean;
}

const Actions = ({ simple, useHide = false, useMargin = true }: ActionProps) => {
  return (
    <Container useMargin={useMargin} simple={simple}>
      <Buttons onClick={minimize}>
        <Minus width={20} height={20} />
      </Buttons>
      <Buttons onClick={useHide ? hide : close}>
        <Close width={20} height={20} />
      </Buttons>
    </Container>
  );
};

const Container = styled.div<{ simple: boolean; useMargin: boolean; }>`
  display: flex;
  align-items: center;
  margin: ${({ useMargin, simple }) => useMargin ? `auto -10px auto ${simple ? 'auto' : '16px'}` : ''};
`;

const Buttons = styled(motion.div)`
  display: flex;
  cursor: pointer;
  margin: 0 6px;
  color: ${({ theme }) => theme.colors.paragraph};
  align-items: center;
  justify-content: center;
`;

export default Actions;
