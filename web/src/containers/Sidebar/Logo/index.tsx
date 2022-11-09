import React from "react";
import { motion } from "framer-motion";
import styled from "styled-components";

import { Typography } from "../../../elements";

const root = {
  open: {
    width: 215,
    height: "auto",
    transition: {
      duration: 0.45,
    },
  },
  collapsed: {
    width: 80,
    height: "auto",
    transition: {
      duration: 0.45,
    },
  },
};

const logo = {
  open: {
    x: 0,
    transition: {
      duration: 0.45,
    },
  },
  collapsed: {
    x: 48,
    transition: {
      duration: 0.45,
    },
  },
};

const typography = {
  open: {
    x: 0,
    opacity: 1,
    transition: {
      delay: 0.05,
      opacity: { delay: 0.25 },
      x: { delay: 0.2 },
      duration: 0.35,
    },
  },
  collapsed: {
    x: -100,
    opacity: 0,
    transition: {
      x: { delay: 0.25 },
      duration: 0.3,
    },
  },
};

const Icon = () => (
  <SVG
    width="40"
    height="37"
    viewBox="0 0 48 37"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      d="M7.45772 25.0535L8.20662 22.3375L6.25949 20.8896L4.85156 23.7055L5.52057 25.3331L7.45772 25.0535Z"
      fill="#F0F0F0"
    />
    <path
      d="M35.4567 4.05371L38.7219 4.68279L38.5921 7.02933L38.0829 7.56853L34.9375 5.84108L35.4567 4.05371Z"
      fill="#F0F0F0"
    />
    <path
      d="M24.0537 36.0469C34.0078 36.0469 42.0771 27.9775 42.0771 18.0234C42.0771 8.06936 34.0078 0 24.0537 0C14.0996 0 6.03027 8.06936 6.03027 18.0234C6.03027 27.9775 14.0996 36.0469 24.0537 36.0469Z"
      fill="url(#paint0_linear)"
    />
    <path
      d="M35.4568 4.05402C36.5851 4.50335 37.4438 5.20233 37.993 6.13096C40.3795 10.2649 35.327 17.664 26.7097 22.6366C18.0924 27.6093 9.16552 28.2883 6.77904 24.1544C6.22985 23.2058 6.08008 22.0974 6.25981 20.8792C1.45689 25.3726 -0.949562 29.4466 0.348524 31.7033C2.59521 35.5975 15.0069 32.6419 28.0677 25.093C41.1284 17.5441 49.8955 8.27779 47.6488 4.39352C46.3507 2.15681 41.6876 2.17678 35.4568 4.05402Z"
      fill="#F0F0F0"
    />
    <defs>
      <linearGradient
        id="paint0_linear"
        x1="15.0411"
        y1="33.6314"
        x2="33.0667"
        y2="2.40913"
        gradientUnits="userSpaceOnUse"
      >
        <stop stopColor="#776AF2" />
        <stop offset="1" stopColor="#A097FD" />
      </linearGradient>
    </defs>
  </SVG>
);

const Logo = () => {
  return (
    <motion.div variants={root}>
      <Container>
        <IconWrapper variants={logo}>
          <Icon />
        </IconWrapper>
        <motion.div variants={typography}>
          <Name>
            Nebula<Bots>bots</Bots>
          </Name>
        </motion.div>
      </Container>
    </motion.div>
  );
};

const Container = styled.div`
  display: flex;
  margin: 32px 0 24px 0;
  align-items: center;
  justify-content: center;
`;

const IconWrapper = styled(motion.div)`
  display: flex;
  z-index: 999;
  flex-direction: column;

  & > svg,
  path {
    cursor: default;
  }
`;

const Name = styled(Typography.H2)`
  margin-left: 8px;
  font-size: 18px;
  font-weight: 700;
  color: ${({ theme }) => theme.colors.h2};
  display: inline;
  flex-direction: column;
  text-transform: lowercase;
`;

const Bots = styled.span`
  font-weight: 400;
`;

const SVG = styled.svg`
  cursor: default;
`;

export default Logo;
