import React from "react";
import { motion, AnimatePresence } from "framer-motion";
import styled from "styled-components";

type CircleProps = {
  show: boolean;
  character: string;
  index: number;
};

const Container = ({ index, show, character }: CircleProps) => {
  return (
    <AnimatePresence exitBeforeEnter>
      {!show && (
        <Circle
          key={`${index}-circle`}
          animate={{ opacity: !show ? 1 : 0 }}
          transition={{
            opacity: { duration: 0.15, type: "linear" },
          }}
          exit={{ opacity: 0 }}
        />
      )}

      {show && (
        <Character
          key={`${index}-character`}
          initial={{ opacity: 0 }}
          animate={{ opacity: show ? 1 : 0 }}
          transition={{
            opacity: { duration: 0.15, type: "linear" },
          }}
          exit={{ opacity: 0 }}
        >
          {character}
        </Character>
      )}
    </AnimatePresence>
  );
};

const Circle = styled(motion.div)`
  border-radius: 50%;
  width: 8px;
  height: 8px;
  margin-right: 5.5px;
  background-color: ${({ theme }) => theme.colors.paragraph};
`;

const Character = styled(motion.p)`
  font-size: 14px;
  letter-spacing: 0.15em;
  color: ${({ theme }) => theme.colors.paragraph};
  margin: 0;
`;

export default Container;
