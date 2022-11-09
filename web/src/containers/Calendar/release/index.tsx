import React, { useEffect, useState } from "react";
import styled from "styled-components";
import { LazyLoadImage } from "react-lazy-load-image-component";
import moment, { Moment } from "moment";
import { motion } from "framer-motion";
import { ArrowRight } from "react-feather";

import Search from "../search";

import { Typography } from "../../../elements";

interface ReleaseProps {
  date: Moment;
  release: Release | null;
}

const extractStores = (stores: Stores) => {
  const { length } = Object.values(stores);

  if (length === 0) {
    return `No Stores`;
  }

  if (length === 1) {
    return `${length} Store`;
  }

  return `${length} Stores`;
};

const fadeArrowIn = {
  rest: {
    x: 0,
    transition: {
      duration: 0.175,
      type: "tween",
      ease: "easeIn",
    },
  },
  hover: {
    x: -2.5,
    transition: {
      duration: 0.175,
      type: "tween",
      ease: "easeOut",
    },
  },
};

const arrowMotion = {
  rest: { opacity: 0, ease: "easeOut", duration: 0.25, type: "tween" },
  hover: {
    opacity: 1,
    transition: {
      duration: 0.25,
      type: "tween",
      ease: "easeIn",
    },
  },
};

const List = ({ stores = {} }: any) => {
  return (
    <>
      {Object.entries(stores).map(([id, { name, date, datum }]: any) => {
        return (
          <Row initial="rest" whileHover="hover" animate="rest">
            <Flex fd="column">
              <Name>{name}</Name>
              <Time>{moment(Number(date)).format("h:mm A")}</Time>
            </Flex>
            <ActionContainer>
              <ActionText variants={fadeArrowIn}>Use Setup</ActionText>
              <ArrowContainer variants={arrowMotion}>
                <Arrow />
              </ArrowContainer>
            </ActionContainer>
          </Row>
        );
      })}
    </>
  );
};

const Release = ({ date, release }: ReleaseProps) => {
  const [list, setList] = useState(release?.stores);

  useEffect(() => {
    setList(release?.stores);
  }, [release]);

  if (!release) {
    return (
      <Container center>
        <Flex fd="column" jc="center" ai="center">
          <Flex>
            <Title>Please select a release.</Title>
          </Flex>
          <Flex>
            <Date>
              In order to create tasks, you must first select a release.
            </Date>
          </Flex>
        </Flex>
      </Container>
    );
  }

  const { name, image, stores } = release;

  return (
    <Container>
      <Flex fd="column" jc="center">
        <Flex jc="flex-start" ai="center" m="0 0 16px 0">
          <Image
            alt=""
            height={60}
            width={60}
            src={image || "https://nebulabots.s3.amazonaws.com/default.jpeg"}
          />
          <Flex fd="column">
            <Flex>
              <Title>{name || "Untitled Release"}</Title>
              <Pill>{extractStores(stores)}</Pill>
            </Flex>
            <Date>{date.format("MMM Do") || "Unknown Date"}</Date>
          </Flex>
        </Flex>
        <Search
          placeholder="Search for a store"
          list={stores}
          property={"name"}
          setList={setList}
        />
        <ListContainer>
          <List stores={list} />
        </ListContainer>
      </Flex>
    </Container>
  );
};

const Container = styled.div<{ center?: boolean }>`
  background-color: ${({ theme }) => theme.colors.sidebar};
  border-radius: 16px;
  padding: 16px;
  border: none;
  display: flex;
  margin-top: 16px;
  flex-direction: column;
  flex: 1;
  ${({ center }) =>
    center
      ? `
		align-items: center;
		justify-content: center;
	`
      : ""}
`;

const Image = styled(LazyLoadImage)`
  display: flex;
  cursor: default;
  object-fit: fill;
  border-radius: 8px;
  height: 60px;
  width: 60px;
`;

const Title = styled(Typography.H2)`
  display: flex;
  color: ${({ theme }) => theme.colors.subHeading};
  font-size: 16px;
  font-weight: 500;
  margin: 0 16px 2px 16px;
`;

const Date = styled(Typography.Paragraph)`
  display: flex;
  font-size: 12px;
  font-weight: 400;
  opacity: 0.8;
  color: ${({ theme }) => theme.colors.h1};
  margin: 0 16px;
`;

const Pill = styled.div`
  height: 20px;
  min-width: 48px;
  padding: 0px 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  border-radius: 4px;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.primary};
  background-color: ${({ theme }) => theme.colors.secondary};
`;

const Flex = styled.div<{
  fd?: string;
  jc?: string;
  ai?: string;
  m?: string;
  flex?: number;
}>`
  display: flex;
  cursor: pointer;
  ${({ fd }) => (fd ? `flex-direction: ${fd};` : "")}
  ${({ jc }) => (jc ? `justify-content: ${jc};` : "")}
	${({ ai }) => (ai ? `align-items: ${ai};` : "")}
	${({ m }) => (m ? `margin: ${m};` : "")}
	${({ flex }) => (flex ? `flex: ${flex};` : "")}
`;

const Row = styled(motion.div)<{ active?: boolean }>`
	display: flex;
	align-items: center;
	border-radius: 4px;
	cursor: pointer;
	background-color: ${({ theme, active }) =>
    active ? theme.colors.secondary : `${theme.colors.lightHue}, 0.1)`}};
	border: 1px solid ${({ theme, active }) =>
    active ? theme.colors.primary : `${theme.colors.lightHue}, 0.01)`}};
	padding: 6px 14px; 
	margin: 8px 0;
	flex: 1;
	
	&:nth-child(1) {
		margin: 0 0 8px 0;
	}
	
	&:hover {
		& > div {
	    color: ${({ theme }) => theme.colors.primary};
		}
		background-color: ${({ theme }) => theme.colors.lightRow};
	}
`;

const ListContainer = styled.div`
  display: flex;
  flex-direction: column;
  overflow: scroll;

  &::-webkit-scrollbar {
    display: none;
  }
`;

const Name = styled(Typography.H2)`
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.h1};
`;

const Time = styled(Typography.Paragraph)`
  display: flex;
  font-size: 10px;
  font-weight: 400;
  opacity: 0.8;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.h1};
  margin: 2px 0 0 0;
`;

const ActionContainer = styled(motion.div)`
  display: flex;
  margin-left: auto;
  justify-content: center;
  align-items: center;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.h1};
`;

const ActionText = styled(motion.div)`
  display: flex;
  cursor: pointer;
`;

const ArrowContainer = styled(motion.div)`
  display: flex;
  cursor: pointer;
`;

const Arrow = styled(ArrowRight)`
  display: flex;
  width: 14px;
  color: ${({ theme }) => theme.colors.primary};
  margin: 0 0 0 4px;
  cursor: pointer;
`;

export default Release;
