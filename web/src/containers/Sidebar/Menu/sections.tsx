import React, { Fragment } from "react";
import { motion } from "framer-motion";
import styled from "styled-components";
import { useHistory, useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { Typography } from "../../../elements";
import { getIcon } from "../Icons";

import { top, middle, bottom, isActive } from "../../../roots/Main/Routes";

interface MenuProps {
  collapsed: boolean;
}

interface MenuItemProps {
  path: string;
  collapsed: boolean;
  disabled?: boolean;
  name: string;
}

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
  open: {
    transition: {
      duration: 0.45,
    },
  },
  collapsed: {
    transition: {
      duration: 0.35,
    },
  },
};

const border = {
  open: {
    width: 151,
    marginLeft: 0,
    marginRight: 0,
    transition: {
      marginLeft: {
        delay: 0.1,
      },
      marginRight: {
        delay: 0.1,
      },
      duration: 0.45,
    },
  },
  collapsed: {
    width: 32,
    marginLeft: 8,
    marginRight: 8,
    transition: {
      width: { duration: 0.45 },
      duration: 0.35,
    },
  },
};

const typography = {
  open: {
    x: 0,
    opacity: 1,
    transition: {
      delay: 0.05,
      duration: 0.35,
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

const MenuItem = ({ name, path, collapsed }: MenuItemProps) => {
  const { pathname } = useLocation();
  const { t } = useTranslation();

  const history = useHistory();

  const active = isActive(path, pathname);
  const Icon = getIcon(name);

  if (active) {
    return (
      <Active collapsed={collapsed} key={t(`${name}`).trim()} variants={item}>
        <IconWrapper initial={collapsed ? "collapsed" : "open"} variants={logo}>
          <Icon height={14} />
        </IconWrapper>
        <motion.div
          initial={collapsed ? "collapsed" : "open"}
          variants={typography}
        >
          <Text>{t(`Pages.${name}`)}</Text>
        </motion.div>
      </Active>
    );
  }

  return (
    <motion.div key={`${name}`} whileTap={{ scale: 0.95 }} variants={item}>
      <Item collapsed={collapsed} onClick={() => history.push(path)}>
        <IconWrapper initial={collapsed ? "collapsed" : "open"} variants={logo}>
          <Icon height={14} />
        </IconWrapper>
        <motion.div
          initial={collapsed ? "collapsed" : "open"}
          variants={typography}
        >
          <Text>{t(`Pages.${name}`)}</Text>
        </motion.div>
      </Item>
    </motion.div>
  );
};

const sections = [top, middle, bottom];

const Sections = ({ collapsed }: MenuProps) => {
  return (
    <Fragment key="Menu-Sections">
      {sections.map((section, index) => {
        const children = Object.values(section).map((item) =>
          MenuItem({ ...item, collapsed })
        );

        return (
          <Fragment key={index}>
            {children}
            <Border variants={border} />
          </Fragment>
        );
      })}
    </Fragment>
  );
};

const Border = styled(motion.div)`
  margin: 8px 0;
  border-bottom: 1px solid ${({ theme }) => theme.colors.separator};
`;

const Active = styled(motion.div)<{ collapsed: boolean }>`
  display: flex;
  align-items: center;
  padding: 2px 16px;
  margin: 4px 0;
  border-radius: ${({ collapsed }) =>
    collapsed ? "4px 4px 4px 4px" : "4px 0 0 4px"};
  background-color: ${({ theme }) => theme.colors.secondary};
  -webkit-user-drag: none;
  user-drag: none;

  * {
    color: ${({ theme }) => theme.colors.primary} !important;
  }
`;

const IconWrapper = styled(motion.div)`
  display: flex;
  flex-direction: column;
  z-index: 500;
  width: 20px;
  color: ${({ theme }) => theme.colors.h2};
`;

const Text = styled(Typography.H2)`
  display: flex;
  flex-direction: column;
  margin-left: 16px;
  font-size: 12px;
  z-index: 499;
  font-weight: 400;
  margin: 8px 16px;
  color: ${({ theme }) => theme.colors.h2};
`;

const Item = styled.div<{ collapsed: boolean }>`
  display: flex;
  align-items: center;
  padding: 2px 16px;
  margin: 4px 0;
  cursor: pointer;
  border-radius: ${({ collapsed }) =>
    collapsed ? "4px 4px 4px 4px" : "4px 0 0 4px"};
  -webkit-user-drag: none;
  user-drag: none;

  &:hover {
    * {
      cursor: pointer;
      color: ${({ theme }) => theme.colors.primary} !important;
    }
    background-color: ${({ theme }) => theme.colors.secondary};
  }
`;

export default Sections;
