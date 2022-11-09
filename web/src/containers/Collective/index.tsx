import React, {
  useCallback,
  useState,
  forwardRef,
  cloneElement,
  useRef,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import { Popover } from "react-tiny-popover";
import styled from "styled-components";
import { RotateCw, Plus, User, Search, Trash } from "react-feather";
import { FixedSizeList as List } from "react-window";
import AutoSizer from "react-virtualized-auto-sizer";
import memoize from "memoize-one";
import { motion } from "framer-motion";

import { Typography, Input } from "../../elements";
import Footer from "./footer";

import {
  makeHarvesters,
  makeFocusedHarvester,
  add,
} from "../../stores/Collective/reducers/harvester";
import {
  add as addAccount,
  select as selectAccount,
  remove as removeAccount,
  select,
} from "../../stores/Collective/reducers/accounts";

import Pill from "./Harvester";
import { makeAccountsOptions } from "../../stores/Collective/reducers/accounts";
import Navigation from "./navigation";

const PADDING_SIZE = 4;

const createItemData = memoize((accounts) => ({
  accounts,
}));

const innerElementType = forwardRef(({ style, ...rest }: any, ref) => (
  <div
    ref={ref}
    style={{
      ...style,
      height: `${parseFloat(style.height) + PADDING_SIZE * 2}px`,
    }}
    {...rest}
  />
));

const AccountRow = ({ index, style, data }: any) => {
  const dispatch = useDispatch();

  const { accounts } = data;
  const account = accounts[index];

  const handleRemove = useCallback(() => {
    dispatch(removeAccount(account.id));
  }, [account]);

  const handleSelect = useCallback(() => {
    dispatch(selectAccount(account.id));
  }, [account]);

  return (
    <AccountRowComponent
      style={{
        ...style,
        top: `${parseFloat(style.top) + PADDING_SIZE}px`,
      }}
      onClick={handleSelect}
    >
      <AccountRowText>{account.email}</AccountRowText>
      <RemoveRow
        whileHover={{
          scale: 1.1,
        }}
        whileTap={{
          scale: 0.9,
        }}
      >
        <IconContainer onClick={handleRemove} component={<Trash />} />
      </RemoveRow>
    </AccountRowComponent>
  );
};

const Actions = () => {
  const accounts = useSelector(makeAccountsOptions);
  const [filter, setFilter] = useState("");
  const [open, setOpen] = useState(false);
  const dispatch = useDispatch();

  const handleAddHarvester = useCallback(() => {
    dispatch(add());
  }, []);

  const reload = useCallback(() => {
    // TODO: We might not want to just reload, we may need to reset stuff too?
  }, []);

  const handleOpenAccounts = useCallback(() => setOpen((p) => !p), []);

  if (open) {
    window.RPCAction("collective:hide");
  } else {
    window.RPCAction("collective:show");
  }

  const itemData = createItemData(accounts);

  return (
    <Row m="0 0 0 auto">
      <Refresh onClick={reload} />
      <Add onClick={handleAddHarvester} />
      <Popover
        isOpen={open}
        padding={4}
        align="end"
        positions={["bottom"]}
        onClickOutside={handleOpenAccounts}
        reposition={false}
        content={
          <PopoverContainer>
            <Input
              id="filter"
              name="filter"
              value={filter}
              Icon={Search}
              placeholder="Search accounts"
              onChange={(e) => setFilter(e.target.value)}
            />
            <AddAccount
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              onClick={() => {
                dispatch(addAccount());
              }}
            >
              Add Account
            </AddAccount>
            <OverflowList>
              <AutoSizer>
                {({ height, width }) => (
                  <StyledList
                    height={height}
                    width={width}
                    innerElementType={innerElementType}
                    itemCount={accounts.length}
                    itemData={itemData}
                    itemSize={24}
                  >
                    {AccountRow}
                  </StyledList>
                )}
              </AutoSizer>
            </OverflowList>
          </PopoverContainer>
        }
      >
        <Account onClick={handleOpenAccounts} />
      </Popover>
    </Row>
  );
};

const Collective = () => {
  const contentWrapper = useRef(null);

  const harvesters = useSelector(makeHarvesters);
  const focused = useSelector(makeFocusedHarvester);

  return (
    <Container>
      <Title>Harvester Collective</Title>
      <AltRow>
        <Navigation {...{ contentWrapper }} />
        <OverflowContainer ref={contentWrapper}>
          {harvesters
            .sort((a, b) => a.index - b.index)
            .map((h) => (
              <Pill {...h} />
            ))}
        </OverflowContainer>
        <Actions />
      </AltRow>
      {focused && <Footer {...focused} />}
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;
`;

const Title = styled(Typography.H2)`
  display: flex;
  font-size: 18px;
  margin: 16px;
  color: ${({ theme }) => theme.colors.h2};
  font-weight: 700;
`;

const AltRow = styled.div`
  display: flex;
  padding: 16px;
  margin: 0 -4px;
  background: ${({ theme }) => theme.colors.sidebar};
`;

const OverflowContainer = styled.div`
  display: flex;
  overflow-x: scroll;
  overflow-y: hidden;
  position: relative;

  &::-webkit-scrollbar {
    display: none;
  }
`;

const Row = styled.div<{ m?: string }>`
  display: flex;
  justify-content: center;
  align-items: center;
  ${({ m }) =>
    typeof m !== "undefined"
      ? `
		margin: ${m};
	`
      : ""}
`;

const Refresh = styled(RotateCw)`
  display: flex;
  width: auto;
  height: 16px;
  margin: 0 4px 0 12px;
  color: ${({ theme }) => theme.colors.h2};
`;

const Add = styled(Plus)`
  display: flex;
  width: auto;
  height: 16px;
  stroke-width: 2.5px;
  margin: 0 4px 0 8px;
  color: ${({ theme }) => theme.colors.h2};
`;

const Account = styled(User)`
  display: flex;
  width: auto;
  height: 16px;
  stroke-width: 2.5px;
  margin: 0 8px;
  padding: 4px;
  border-radius: 50%;
  justify-content: center;
  align-items: center;
  color: ${({ theme }) => theme.colors.h2};
  background-color: ${({ theme }) => theme.colors.day};
`;

const PopoverContainer = styled.div`
  height: fit-content;
  width: 168px;
  padding: 8px;
  border: 1px solid ${({ theme }) => theme.colors.border};
  box-shadow: rgba(0, 0, 0, 0.15) 0px 5px 15px 0px;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  color: ${({ theme }) => theme.colors.h2};
  background-color: ${({ theme }) => theme.colors.background};
  z-index: 100;
`;

const OverflowList = styled.div`
  position: relative;
  flex: 1 1 auto;
  width: 100%;
  height: 124px;
  overflow: hidden;
  font-size: 12px;
  color: ${({ theme }) => theme.colors.h2};
  font-weight: 400;
  margin: 0;
  display: flex;

  &::-webkit-scrollbar {
    display: none;
  }
`;

const StyledList = styled(List)`
  &::-webkit-scrollbar {
    display: none;
  }
`;

const AddAccount = styled(motion.div)`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 24px;
  margin: 8px 0;
  border-radius: 4px;
  cursor: pointer;
  background-color: ${({ theme }) => theme.colors.secondary};
  color: ${({ theme }) => theme.colors.primary};
  border: 1px solid ${({ theme }) => theme.colors.primary};
  font-size: 12px;
  font-weight: 500;
`;

const AccountRowComponent = styled.div`
  display: flex;
  align-items: center;
  border-radius: 2px;
  padding: 0 8px;
  cursor: pointer;

  &:hover {
    background-color: ${({ theme }) => theme.colors.primary};
    color: #fff;
  }
`;

const AccountRowText = styled.div`
  height: 100%;
  width: 80%;
  overflow: hidden;
  display: flex;
  justify-content: left;
  align-items: center;
`;

const RemoveRow = styled(motion.div)`
  display: flex;
  justify-content: center;
  align-items: center;
`;

const IconContainer = styled(({ component, ...props }) =>
  cloneElement(component, props)
)`
  width: auto;
  height: 12px;
  fill: ${({ theme }) => theme.colors.failed};
  color: ${({ theme }) => theme.colors.failed};
  fill-opacity: 0.25;
  margin-left: 8px;
  cursor: pointer;
`;

export default Collective;
