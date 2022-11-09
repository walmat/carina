import { useState, useMemo, cloneElement } from "react";
import '@szhsin/react-menu/dist/index.css';
import {
  FocusableItem,
  ControlledMenu,
  MenuGroup,
  MenuItem,
  MenuHeader,
  MenuDivider,
  SubMenu
} from '@szhsin/react-menu';
import styled from 'styled-components';

import { Typography } from '../../elements';
import { Items } from "./types";
import * as React from "react";

type ContextMenuProps = {
  // id: string;
  row: any,
  addGroup: any;
  group: any;
  groups: any[];
  items: Items[];
  closeMenu: any;
  setAnchorPoint: any;
  anchorPoint: {
    x: number;
    y: number;
  };
  menuProps: any;
};

const getColors = (name: string) => {
  switch (name) {
    default:
      return "paragraph";
  }
};

const ContextMenu = (({ row, addGroup, group, groups, items, anchorPoint, menuProps, closeMenu, setAnchorPoint }: ContextMenuProps) => {

  const [filter, setFilter] = useState('');

  const data = Object.values(groups).filter(({ id }) => id !== group.id);

  return useMemo(() => (
    <ControlledMenu
      {...menuProps}
      onClose={() => {
        setFilter('');
        setAnchorPoint({ x: 0, y: 0 });
        closeMenu();
      }}
      anchorPoint={anchorPoint}
    >
      {
        items.map(({ name: groupName, actions }) => {

          if (groupName === 'Group' && !data.length) {
            return (
              <MenuGroup>
                <MenuHeader>{groupName} Actions</MenuHeader>
                <MenuDivider/>
                <MenuItem disabled>
                  <AddText
                    // @ts-ignore
                    onClick={(e: any) => {
                      e.preventDefault();
                      addGroup(true);
                      closeMenu();
                    }}
                  >
                    Add {items[0].name} Group
                  </AddText>
                </MenuItem>
              </MenuGroup>
            );
          }

          return (
            <MenuGroup>
              <MenuHeader>{groupName} Actions</MenuHeader>
              <MenuDivider/>
              {
                actions.map(({name, Icon, shortcut, onClick}) => {
                  const fillColor = getColors(name);

                  if (groupName === 'Group') {
                    return (
                      <SubMenu
                        label={
                          <>
                            <IconContainer
                              fillColor={fillColor}
                              component={<Icon/>}
                            />
                            <Text>{name}</Text>
                          </>
                        }
                      >
                        <MenuGroup>
                            <FocusableItem>
                              {({ref}) => (
                                <InputContainer showLabel={!!filter.length}>
                                  <label htmlFor="filter">Filter groups</label>
                                  <input
                                    autoFocus
                                    id="filter"
                                    name="filter"
                                    ref={ref}
                                    type="text"
                                    placeholder="Filter groups"
                                    value={filter}
                                    onChange={(e: any) => setFilter(e.target.value)}
                                  />
                                </InputContainer>
                              )}
                            </FocusableItem>
                            {
                              data
                                .filter(({ name }) => name.toUpperCase().includes(filter.trim().toUpperCase()))
                                .map(({ id, name }) => {
                                return (
                                  <MenuItem onClick={() => onClick(id)}>
                                    <Text>{name}</Text>
                                  </MenuItem>
                                )
                              })
                            }
                        </MenuGroup>
                      </SubMenu>
                    )
                  }

                  return (
                    <MenuItem
                      onClick={() => onClick(row.id)}
                    >
                      <IconContainer
                        fillColor={fillColor}
                        component={<Icon/>}
                      />
                      <Text>{name}</Text>
                      {shortcut ? <Shortcut>{shortcut.join(" ")}</Shortcut> : null}
                    </MenuItem>
                  )
                })
              }
            </MenuGroup>
          );
        })
      }
    </ControlledMenu>
  ), [row, anchorPoint, menuProps, filter]);
});

const InputContainer = styled.div<{ showLabel: boolean }>`
  position: relative;
  display: flex;
  flex: 1 1 auto;
  max-height: 44px;
  max-width: 140px;

  label {
    transition: all 0.2s ease-out;
    font-size: 11px;
    font-weight: 400;
    position: absolute;
    top: 0;
    width: 100%;
    padding: 0 0.9166666667em;
    z-index: 1;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    -webkit-transform: translateY(3px);
    transform: ${({ showLabel }) => (showLabel ? "none" : "translateY(2.5px)")};
    pointer-events: none;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    -webkit-box-sizing: border-box;
    box-sizing: border-box;
    color: ${({ theme }) => theme.colors.tooltip};
    opacity: ${({ showLabel }) => (showLabel ? 1 : 0)};
    margin: 0.5em 0;
    padding: 0 11px;
    padding-left: 9px;
  }

  div {
    color: #ccc;
    background: linear-gradient(
      to right,
      ${({ theme }) => theme.colors.fade}, 0) 0%,
      ${({ theme }) => theme.colors.fade}, 1) 20%
    );
    position: absolute;
    top: 50%;
    right: 4px;
    transform: translateY(-50%);
    padding-right: 8px;
    padding-left: 8px;
  }

  div,
  input {
    border-radius: 4px;
    font-size: 12px;
    font-family: Inter, sans-serif;
    font-weight: 400;
  }

  div,
  span {
    display: block;
    font-weight: 700;
  }

  input {
    max-width: 124px;
    outline: none;
    flex: 1;
    margin: 0;
    transition: all 0.2s ease-out;
    padding: 1.165em 8px;
    padding-right: 8px;
    padding-left: 8px;
    color: ${({ theme }) => theme.colors.h2};
    background-color: ${({ theme }) => theme.colors.sidebar};
    border: 1px solid ${({ theme }) => theme.colors.border};

    ${({ showLabel }) =>
  showLabel
    ? `
      padding-top: 1.845em;
      padding-bottom: .45em;
    `
    : ""}

    &:hover, &:focus {
      border-color: ${({ theme }) => theme.colors.h2};
    }

    &::placeholder {
      color: ${({ theme }) => theme.colors.placeholder};
    }
  }
`;

const IconContainer = styled(({ component, ...props }) =>
  cloneElement(component, props)
)`
  width: 14px;
  height: 14px;
  fill: ${({ theme, fillColor }) => theme.colors[fillColor]};
  color: ${({ theme, fillColor }) => theme.colors[fillColor]};
  fill-opacity: 0.05;
  margin-right: 8px;
  cursor: pointer;
`;

const Text = styled(Typography.Paragraph)`
  margin: 0 auto 0 8px;
  font-size: 14px;
  cursor: pointer;
  text-transform: capitalize;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const AddText = styled(Typography.Paragraph)`
  margin: 0 auto 0 8px;
  font-size: 14px;
  cursor: pointer;
  text-transform: capitalize;
  color: ${({ theme }) => theme.colors.primary};
`;

const Shortcut = styled.div`
  width: 32px;
  padding: 4px;
  text-align: center;
  border-radius: 4px;
  font-size: 12px;
  color: ${({ theme }) => theme.colors.primary};
  background-color: ${({ theme }) => theme.colors.secondary};
`;

export default ContextMenu;
