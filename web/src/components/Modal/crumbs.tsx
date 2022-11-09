import React, { Dispatch, SetStateAction } from "react";
import { motion } from "framer-motion";
import styled from "styled-components";

import { ChevronRight } from "react-feather";

import { Typography } from "../../elements";

interface StepMap {
  [id: string]: {
    done: string[];
    active: string;
    forward: string[];
  };
}

interface CrumbsProps {
  step: string;
  steps: string[];
  stepMap: StepMap;
  matches: boolean;
  setStep: Dispatch<SetStateAction<string>>;
  progressStepper: any;
}

const variants = {
  shown: {
    display: "flex",
    opacity: 1,
    transition: {
      duration: 0.45,
    },
  },
  hidden: {
    opacity: 0,
    transitionEnd: {
      display: "none",
    },
  },
};

const isForward = (step: string, going: string, stepMap: StepMap) => {
  return !stepMap[going].done.some((s) => s === step);
};

const Crumbs = ({
  step,
  steps,
  matches,
  stepMap,
  setStep,
  progressStepper,
}: CrumbsProps) => {
  return (
    <Container>
      <Crumb
        done={stepMap[steps[0]].done.some((s) => s === step)}
        active={stepMap[steps[0]].active === step}
        // @ts-ignore
        onClick={() => {
          if (isForward(step, steps[0], stepMap)) {
            return progressStepper();
          }

          return setStep(steps[0]);
        }}
      >
        {steps[0]}
      </Crumb>
      <Nav />
      <Payment
        initial="shown"
        animate={matches ? "hidden" : "shown"}
        variants={variants}
      >
        <Crumb
          done={stepMap[steps[1]].done.some((s) => s === step)}
          active={stepMap[steps[1]].active === step}
          // @ts-ignore
          onClick={() => {
            if (isForward(step, steps[1], stepMap)) {
              return progressStepper();
            }

            return setStep(steps[1]);
          }}
        >
          {steps[1]}
        </Crumb>
        <Nav />
      </Payment>

      <Crumb
        done={false}
        active={stepMap[steps[2]].active === step}
        // @ts-ignore
        onClick={() => {
          if (isForward(step, steps[2], stepMap)) {
            return progressStepper();
          }

          return setStep(steps[2]);
        }}
      >
        {steps[2]}
      </Crumb>
    </Container>
  );
};

const Container = styled(motion.div)`
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin: 0 8px;
`;

const Nav = styled(ChevronRight)`
  color: ${({ theme }) => theme.colors.paragraph};
  width: 24px;
  height: 16px;
  margin: 0 -4px;
`;

const Crumb = styled(Typography.H2)<{ active: boolean; done?: boolean }>`
  display: flex;
  flex-direction: column;
  justify-content: center;
  cursor: ${({ active, done }) => (active || done ? "default" : "pointer")};
  font-size: 14px;
  font-weight: 500;
  margin: 0 8px;
  text-transform: capitalize;
  color: ${({ theme, active, done }) =>
    active
      ? theme.colors.primary
      : done
      ? theme.colors.h3
      : theme.colors.description};

  &:hover {
    opacity: ${({ active, done }) => (active || done ? 1 : 0.6)};
  }
`;

const Payment = styled(motion.div)`
  display: flex;
  align-items: center;
`;

export default Crumbs;
