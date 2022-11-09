import React, { Dispatch, SetStateAction, useCallback } from "react";
import styled from "styled-components";
import { useHotkeys } from "react-hotkeys-hook";

type TableHeaderProps = {
  actions: { name: string; onClick: any }[];
  hotkeysEnabled: boolean;
  headerGroups: any[];
  selected?: string[];
  setSelected?: Dispatch<SetStateAction<string[]>>;
  rows: object[];
};

const TableHeader = (props: TableHeaderProps) => {
  const { actions, hotkeysEnabled, headerGroups, setSelected, selected, rows } = props;

  const onSelectAll = useCallback(() => {
    if (setSelected) {
      if (selected?.length === rows.length) {
        setSelected([]);
      } else {
        setSelected(rows.map((r: any) => r.id));
      }
    }
  }, [selected, rows]);

  // NOTE: MetaKey/Shift + A KeyPress Handler
  useHotkeys("shift+a,command+a,ctrl+a", onSelectAll, { enabled: hotkeysEnabled }, [onSelectAll]);

  return (
    <Container onClick={onSelectAll} padding={actions.length}>
      {headerGroups.map((headerGroup) => (
        <TableHeaderCell {...headerGroup.getHeaderGroupProps()}>
          {headerGroup.headers.map((column: any, index: number) => {
            return (
              <HeaderCell
                // @ts-ignore
                {...column.getHeaderProps(column.getSortByToggleProps())}
              >
                {/* @ts-ignore */}
                {index !== headerGroup.headers.length - 1 &&
                column.canResize ? (
                  // @ts-ignore
                  <ResizerContainer {...column.getResizerProps()}>
                    <Resizer />
                  </ResizerContainer>
                ) : null}
                {column.render("Header")}
              </HeaderCell>
            );
          })}
        </TableHeaderCell>
      ))}
    </Container>
  );
};

const Container = styled.div<{ padding: number }>`
  display: flex;
  flex-direction: column;
  flex: 1;
  cursor: pointer;
  max-height: 36px;
  margin-right: 12px;
  background-color: ${({ theme }) => theme.colors.sidebar};
  color: ${({ theme }) => theme.colors.h2};
  padding: 0 ${({ padding }) => (padding === 2 ? 69 : 94)}px 0 12px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
`;

const TableHeaderCell = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
`;

const HeaderCell = styled.div<{ flex?: number }>`
  padding: 8px 16px 8px 8px;
  ${({ flex }) => (typeof flex !== undefined ? `flex: ${flex};` : "")}
`;

const ResizerContainer = styled.div`
  padding: 0 8px;
  top: 0;
  right: 0;
  width: 1px;
  height: 100%;
  border-radius: 4px;
  position: absolute;
  cursor: col-resize;
  z-index: 1;
  display: flex;
  touch-action: none;
`;

const Resizer = styled.div`
  width: 1px;
  height: 16px;
  margin: auto 0;
  background: ${({ theme }) => theme.colors.h3};
  position: relative;
  cursor: col-resize;
  display: flex;
  cursor: resize;
  justify-content: center;
  align-items: center;
  touch-action: none;
`;

export default TableHeader;
