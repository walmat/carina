import React, { Dispatch, SetStateAction, useEffect, useState } from "react";
import styled from "styled-components";
import { motion, AnimatePresence } from "framer-motion";
import { Moment } from "moment";
import { ArrowRight } from "react-feather";

import { Typography } from "../../elements";

import Search from "./search";

interface ReleasesProps {
  date?: Moment;
  release: any;
  setRelease: Dispatch<SetStateAction<string>>;
  releases: Releases;
}

const extractDrops = (releases: Releases) => {
  const { length } = Object.values(releases);

  if (length === 0) {
    return `No Releases`;
  }

  if (length === 1) {
    return `${length} Release`;
  }

  return `${length} Releases`;
};

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

const transition = {
  opacity: {
    duration: 0.15,
    type: "linear",
    ease: [0.43, 0.13, 0.23, 0.96],
  },
};

const List = ({ release, setRelease, releases }: ReleasesProps) => {
  if (!Object.entries(releases).length) {
    return (
      <Center>
        <Flex fd="column" jc="center" ai="center">
          <Flex m="0 0 2px 0">
            <Title>No releases available.</Title>
          </Flex>
          <Flex>
            <Date>Please choose a different date to get started.</Date>
          </Flex>
        </Flex>
      </Center>
    );
  }

  return (
    <>
      {Object.entries(releases).map(([id, { name, image, stores }]) => {
        const active = id === release;

        const handleSetRelease = () => {
          if (!active) {
            setRelease(id);
          }
        };

        return (
          <Row
            onClick={handleSetRelease}
            active={active}
            initial="rest"
            whileHover="hover"
            animate="rest"
          >
            <Image
              src={image || "https://nebulabots.s3.amazonaws.com/default.jpeg"}
            />
            <Flex fd="column">
              <Name>{name}</Name>
              <Sites>{extractStores(stores)}</Sites>
            </Flex>
            <AnimatePresence exitBeforeEnter>
              {active && (
                <Flex flex={1} ai="center" jc="flex-end">
                  <ViewingRelease
                    key="viewing-release"
                    initial={{ opacity: 0 }}
                    onClick={handleSetRelease}
                    exit={{ opacity: 0 }}
                    transition={transition}
                    animate={{ opacity: active ? 1 : 0 }}
                  >
                    Viewing Release
                  </ViewingRelease>
                  <Arrow />
                </Flex>
              )}

              {!active && (
                <ViewRelease
                  key="view-release"
                  exit={{ opacity: 0 }}
                  transition={transition}
                  animate={{ opacity: !active ? 1 : 0 }}
                >
                  View Release
                </ViewRelease>
              )}
            </AnimatePresence>
          </Row>
        );
      })}
    </>
  );
};

const Releases = ({
  date,
  release,
  setRelease,
  releases = {},
}: ReleasesProps) => {
  const [list, setList] = useState(releases);

  useEffect(() => {
    setList(releases);
  }, [releases]);

  return (
    <Container>
      <Flex flex={1} of="hidden" fd="column" jc="center">
        <Flex flex={0} m="0 0 16px 0">
          <Title>{date?.format("MMMM Do")} Releases</Title>
          <Pill>{extractDrops(releases)}</Pill>
        </Flex>
        <Search
          placeholder="Search for a release"
          list={releases}
          property={"name"}
          setList={setList}
        />
        <ListContainer>
          <List release={release} releases={list} setRelease={setRelease} />
        </ListContainer>
      </Flex>
    </Container>
  );
};

const Container = styled.div`
  background-color: ${({ theme }) => theme.colors.sidebar};
  border-radius: 16px;
  border: none;
  display: flex;
  padding: 16px;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
`;

const Center = styled.div`
  background-color: ${({ theme }) => theme.colors.sidebar};
  border-radius: 16px;
  padding: 16px;
  border: none;
  display: flex;
  margin-top: 16px;
  flex-direction: column;
  flex: 1;
  align-items: center;
  justify-content: center;
`;

const Date = styled(Typography.Paragraph)`
  display: flex;
  font-size: 12px;
  font-weight: 400;
  opacity: 0.8;
  color: ${({ theme }) => theme.colors.h1};
  margin: 0 16px;
`;

const Title = styled(Typography.H2)`
  display: flex;
  color: ${({ theme }) => theme.colors.subHeading};
  margin: 0 8px 0 0;
  font-size: 16px;
  font-weight: 500;
  margin-right: 8px;
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
  flex?: number | string;
  jc?: string;
  ai?: string;
  m?: string;
  of?: string;
}>`
  display: flex;
  min-width: 0;
  ${({ of }) => (typeof of !== "undefined" ? `overflow: ${of};` : "")}
  ${({ flex }) => (typeof flex !== "undefined" ? `flex: ${flex};` : "")}
	${({ fd }) => (fd ? `flex-direction: ${fd};` : "")}
	${({ jc }) => (jc ? `justify-content: ${jc};` : "")}
	${({ ai }) => (ai ? `align-items: ${ai};` : "")}
	${({ m }) => (m ? `margin: ${m};` : "")}
`;

const Row = styled(motion.div)<{ active?: boolean }>`
	display: flex;
	align-items: center;
	border-radius: 4px;
	cursor: pointer;
	margin: 8px 0;
	background-color: ${({ theme, active }) =>
    active ? theme.colors.secondary : `${theme.colors.lightHue}, 0.1)`}};
	border: 1px solid ${({ theme, active }) =>
    active ? theme.colors.primary : `${theme.colors.lightHue}, 0.01)`}};
	height: 41px;
	padding: 16px;
	
	&:after {
	  content: "";
	  display: flex;
	  flex: 1;
	  position: absolute;
	  z-index: 1;
	  bottom: 0;
	  left: 0;
	  pointer-events: none;
	  background-image: linear-gradient(to bottom, rgba(255,255,255, 0), rgba(255,255,255, 1) 90%);
	  height: 1.2em;
	}
	
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

const Image = styled.img`
  height: 40px;
  width: 40px;
  object-fit: fill;
  border-radius: 4px;
  margin-right: 16px;
`;

const ListContainer = styled.div`
  display: flex;
  flex-direction: column;
  overflow: scroll;
  flex: 1;

  &::-webkit-scrollbar {
    display: none;
  }
`;

const Name = styled(Typography.H2)`
  margin: 0 0 2px 0;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.h1};
`;

const Sites = styled(Typography.Paragraph)`
  margin: 2px 0 0 0;
  font-size: 12px;
  font-weight: 400;
  opacity: 0.75;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.h1};
`;

const ViewRelease = styled(motion.div)`
  display: flex;
  margin: 0 0 0 auto;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.h1};
`;

const ViewingRelease = styled(motion.div)`
  display: flex;
  margin: 0 0 0 auto;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.viewing};
`;

const Arrow = styled(ArrowRight)`
  display: flex;
  width: 14px;
  color: ${({ theme }) => theme.colors.viewing};
  margin: 0 0 0 4px;
  cursor: pointer;
`;

export default Releases;
