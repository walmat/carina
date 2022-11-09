import React from "react";
import { useSelector } from "react-redux";
import { motion } from "framer-motion";
import styled from "styled-components";

import { Typography } from "../../../elements";

import { makeUser } from "../../../stores/Main/reducers/user";

const root = {
  open: {
    marginLeft: 32,
    marginRight: 32,
    marginTop: "auto",
    marginBottom: 32,
    transition: {
      duration: 0.45,
    },
  },
  collapsed: {
    marginLeft: 16,
    marginRight: 16,
    marginTop: "auto",
    marginBottom: 32,
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
    x: 8,
    transition: {
      duration: 0.45,
    },
  },
};

const img = {
  open: {
    width: 40,
    height: 40,
    transition: {
      duration: 0.45,
    },
  },
  collapsed: {
    width: 32,
    height: 32,
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
      duration: 0.45,
    },
  },
  collapsed: {
    x: -55,
    opacity: 0,
    transition: {
      duration: 0.45,
    },
  },
};

const User = () => {
  const user = useSelector(makeUser);

  let displayName = "";
  let subtext = "";
  let avatarUrl = "";
  if (user) {
    const { id, avatar, email, type } = user;
    displayName = email;
    subtext = type;

    if (avatar) {
      if (avatar.startsWith("a_")) {
        avatarUrl = `https://cdn.discordapp.com/avatars/${id}/${avatar}.gif`;
      } else {
        avatarUrl = `https://cdn.discordapp.com/avatars/${id}/${avatar}.png`;
      }
    }
  }

  // if (!avatarUrl) {
  //   avatarUrl =
  //     "https://discord.com/assets/c09a43a372ba81e3018c3151d4ed4773.png";
  // }

  if (!displayName) {
    displayName = "default";
  }

  if (!subtext) {
    subtext = "user#0001";
  }
  
  return (
    <Container variants={root}>
      {!!avatarUrl && (
        <Icon variants={logo}>
          <Avatar variants={img} src={avatarUrl} />
        </Icon>
      )}
      <Text variants={typography}>
        <Username>{displayName}</Username>
        <SubText>{subtext}</SubText>
      </Text>
    </Container>
  );
};

const Container = styled(motion.div)`
  display: flex;
  margin-top: auto;
`;

const Icon = styled(motion.div)`
  display: flex;
  color: ${({ theme }) => theme.colors.h3};
`;

const Avatar = styled(motion.img)`
  cursor: default;
  background-color: transparent;
  border-radius: 50%;
`;

const Text = styled(motion.div)`
  margin: auto 0 auto 16px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
`;

const SubText = styled(Typography.H2)`
  display: flex;
  font-size: 12px;
  font-weight: 400;
  margin: 0;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const Username = styled(Typography.H2)`
  display: flex;
  font-size: 14px;
  font-weight: 700;
  text-transform: capitalize;
  margin: 0 0 4px 0;
  color: ${({ theme }) => theme.colors.h2};
`;

export default User;
