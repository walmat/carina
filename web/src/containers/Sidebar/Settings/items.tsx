import React, {useCallback} from "react";
import {AnimatePresence, motion} from "framer-motion";
import styled from "styled-components";
import { useTranslation } from "react-i18next";
import { useHistory, useLocation } from "react-router-dom";

import { getIcon } from "../Icons";
import { Typography } from "../../../elements";

import { settings, isActive } from "../../../roots/Main/Routes";
import { ReCAPTCHA } from "../../../icons";

type Props = {
  collapsed: boolean;
  isOpen: boolean;
};

type MenuItem = {
  name: string;
  path: string;
  isOpen: boolean;
  collapsed: boolean;
};

const parent = {
  initial: {
    display: "flex",
    flex: 1,
    transition: {
      staggerDirection: -1,
      staggerChildren: 0.1,
      when: "afterChildren",
    },
  },
  animate: {
    display: "flex",
    flex: 1,
    transition: {
      when: "beforeChildren",
      staggerChildren: 0.1,
      staggerDirection: 1,
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
    x: -4,
    transition: {
      duration: 0.45,
    },
  },
};

const item = {
  initial: {
    opacity: 0,
  },
  animate: {
    opacity: 1,
  },
};

const typography = {
  open: {
    x: 0,
    opacity: 1,
    transition: {
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

const recapText = {
  open: {
    x: 0,
    opacity: 1,
    transition: {
      duration: 0.35,
    },
  },
  collapsed: {
    x: -55,
    opacity: 0,
    transition: {
      duration: 0.15,
    },
  },
};

const recap = {
  open: {
    x: 0,
    transition: {
      duration: 0.45,
    },
  },
  collapsed: {
    x: 34,
    transition: {
      duration: 0.45,
    },
  },
};

const Item = ({ name, path, isOpen, collapsed }: MenuItem) => {
  const { pathname } = useLocation();
  const { t } = useTranslation();
  const history = useHistory();

  const active = isActive(path, pathname);
  const Icon = getIcon(name);

  if (active) {
    return (
      <Container
        key={`${name}`}
        collapsed={collapsed}
        initial={isOpen ? "animate" : "initial"}
        variants={item}
      >
        <>
          <IconWrapper
            initial={collapsed ? "collapsed" : "open"}
            variants={logo}
          >
            <Icon height={16} />
          </IconWrapper>
          <motion.div
            initial={collapsed ? "collapsed" : "open"}
            variants={typography}
          >
            <Text active={active}>{t(`Pages.${name}`)}</Text>
          </motion.div>
        </>
      </Container>
    );
  }

  return (
    <motion.div
      key={`${name}`}
      whileTap={{ scale: 0.95 }}
      initial={isOpen ? "animate" : "initial"}
      variants={item}
    >
      <Condensed collapsed={collapsed} onClick={() => history.push(path)}>
        <IconWrapper initial={collapsed ? "collapsed" : "open"} variants={logo}>
          <Icon height={16} />
        </IconWrapper>
        <motion.div
          initial={collapsed ? "collapsed" : "open"}
          variants={typography}
        >
          <Text active={active}>{t(`Pages.${name}`)}</Text>
        </motion.div>
      </Condensed>
    </motion.div>
  );
};

const MenuItems = ({ isOpen, collapsed, ...props }: Props) => {
  let children = settings.map((item) => Item({ ...item, isOpen, collapsed }));

  const showCollective = useCallback(() => {
    astilectron.sendMessage({ type: 'showCollective' })
  }, []);

  if (!isOpen) {
    // @ts-ignore
    children = (
      <FlexColumn
        initial={false}
        animate={{ marginRight: collapsed ? 0 : 32 }}
      >
        <Harvester
          key="harvester"
          onClick={showCollective}
          initial={{ opacity: isOpen ? 1 : 0 }}
          animate={{ opacity: isOpen ? 0 : 1 }}
          exit={{ opacity: 0 }}
          whileHover={{ scale: 1.03 }}
          whileTap={{ scale: 0.97 }}
        >
          <IconWrapper initial={collapsed ? "collapsed" : "open"} variants={recap}>
            <ReCAPTCHA />
          </IconWrapper>
          <motion.div
            initial={collapsed ? "collapsed" : "open"}
            variants={recapText}
          >
            <HarvesterText>
              Collective
            </HarvesterText>
          </motion.div>
        </Harvester>
      </FlexColumn>
    )
  }

  return (
    <FlexNavColumn
      {...props}
      variants={parent}
      initial={["initial", collapsed ? "collapsed" : "open"]}
      animate={
        isOpen
          ? ["animate", collapsed ? "collapsed" : "open"]
          : ["animate", collapsed ? "collapsed" : "open"]
      }
    >
      <AnimatePresence exitBeforeEnter>
        {children}
      </AnimatePresence>
    </FlexNavColumn>
  );
};

const Container = styled(motion.div)<{ collapsed: boolean }>`
  display: flex;
  align-items: center;
  padding: ${({ theme }) => `${theme.space.XS}px ${theme.space.L}px`};
  padding: 2px 16px;
  cursor: pointer;
  margin: ${({ theme }) => `${theme.space.S}px 0`};
  border-radius: ${({ collapsed }) =>
    collapsed ? "4px 4px 4px 4px" : "4px 0 0 4px"};
  color: ${({ theme }) => theme.colors.primary};
  background-color: ${({ theme }) => theme.colors.secondary};
  -webkit-user-drag: none;
  user-drag: none;
  transition: color 0.24s ease-in-out, background-color 0.24s ease-in-out;

  * {
    cursor: pointer;
    transition: color 0.24s ease-in-out, background-color 0.24s ease-in-out;
    color: ${({ theme }) => theme.colors.primary} !important;
  }
`;

const Condensed = styled.div<{ collapsed: boolean }>`
  display: flex;
  margin: 4px 0;
  align-items: center;
  padding: 2px 16px;
  border-radius: ${({ collapsed }) =>
    collapsed ? "4px 4px 4px 4px" : "4px 0 0 4px"};
  color: ${({ theme }) => theme.colors.h2};
  -webkit-user-drag: none;
  user-drag: none;
  cursor: pointer;

  &:hover {
    * {
      cursor: pointer;
      color: ${({ theme }) => theme.colors.primary} !important;
    }
    background-color: ${({ theme }) => theme.colors.secondary};
  }
`;

const IconWrapper = styled(motion.div)`
  display: flex;
  flex-direction: column;
  z-index: 500;
  width: 20px;
  color: ${({ theme }) => theme.colors.h2};
`;

const Text = styled(Typography.H2)<{ active: boolean }>`
  margin-left: 16px;
  font-size: 12px;
  z-index: 499;
  display: flex;
  font-weight: 400;
  padding: 2px 16px;
  margin: 4px 0;
  white-space: nowrap;
  flex-direction: column;
  color: ${({ theme, active }) =>
    active ? theme.colors.primary : theme.colors.h2};
  -webkit-user-drag: none;
  user-drag: none;

  &:hover {
    * {
      cursor: pointer;
      color: ${({ theme }) => theme.colors.primary} !important;
    }
  }
`;

const FlexNavColumn = styled(motion.nav)`
  display: flex;
  flex: 1;
  flex-direction: column;
`;

const FlexColumn = styled(motion.div)`
  display: flex;
  flex: 1;
  margin-bottom: 32px;
  justify-content: flex-end;
  flex-direction: column;
`;

const Harvester = styled(motion.button)`
  border-radius: 4px;
  border: none;
  cursor: pointer;
  padding: 12px 16px;
  background-color: ${({ theme }) => theme.colors.secondary};
  display: flex;
  align-items: center;
  justify-content: center;
`;

const HarvesterText = styled(Typography.Paragraph)`
  font-size: 14px;
  font-weight: 700;
  margin: 0 0 0 8px;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.primary};
`;

export default MenuItems;
