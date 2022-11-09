import React, { Dispatch, SetStateAction, useEffect, useState, useCallback, useMemo } from "react";
import styled from "styled-components";
import {
	Column,
	useTable,
	useResizeColumns,
	useFlexLayout,
	useSortBy,
	useGroupBy,
	useExpanded
} from "react-table";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";
import AutoSizer from 'react-virtualized-auto-sizer';
import update from "immutability-helper";
import { FixedSizeList } from "react-window";
import memoize from 'memoize-one';
import { useMenuState } from '@szhsin/react-menu';
import { throttle } from 'lodash';

import TableHeader from './header';
import TableRow from './row';

import { Items } from './types';
import ContextMenu from './contextMenu';

type SortBy = {
	id: string;
	desc: boolean;
};

type TableProps = {
	columns: Column[];
	addGroup?: any;
	group: any;
	groups: any;
	data: any[];
	view: string;
	sortBy?: SortBy;
	setSortBy?: Dispatch<SetStateAction<SortBy>>;
	hiddenColumns?: string[];
	setHiddenColumns?: Dispatch<SetStateAction<string[]>>;
	filter?: string;
	setFilter?: Dispatch<SetStateAction<string>>;
	filterRow?: (row: any, search: string) => boolean;
	hotkeysEnabled?: boolean;
	selected?: string[];
	setSelected?: Dispatch<SetStateAction<string[]>>;
	groupBy?: string;
	setGroupBy?: Dispatch<SetStateAction<string>>;
	actions: { name: string; onClick: any }[];
	items: Items[];
};

const defaultColumn = {
	minWidth: 64,
	width: 100,
	maxWidth: 200,
};

const GUTTER_SIZE = 8;

const extractRowHeight = (view: string) => view === 'Basic' ? 36 : 48;

const createItemData = memoize((rowHeight, getRowId, moveRow, rows, prepareRow, selected, setSelected, actions, items, openMenu, closeMenu, menuProps, anchorPoint, setAnchorPoint, setContextMenuRow) => ({
	rowHeight,
	getRowId,
	moveRow,
	rows,
	prepareRow,
	selected,
	setSelected,
	actions,
	items,
	openMenu,
	closeMenu,
	menuProps,
	anchorPoint,
	setAnchorPoint,
	setContextMenuRow
}));

const Table = (props: TableProps) => {
	const {
		columns,
		addGroup,
		group,
		groups,
		data,
		view,
		filter,
		selected,
		setSelected,
		actions,
		items,
		sortBy,
		hiddenColumns,
		groupBy,
		filterRow,
		hotkeysEnabled = true
	} = props;

	const [anchorPoint, setAnchorPoint] = useState({ x: 0, y: 0 });
	const [contextMenuRow, setContextMenuRow] = useState(null);
	const { openMenu, closeMenu, ...menuProps } = useMenuState();

	const rowHeight = useMemo(() => extractRowHeight(view), [view]);
	const [tableData, setTableData] = useState<any[]>([]);

	useEffect(() => {
		if (data) {
			setTableData(data);
		}
	}, [data]);

	useEffect(() => {
		if (data) {
			if (filter && filterRow) {
				setTableData(
					data.filter((d) => filterRow(d, filter))
				);
			} else {
				setTableData(data);
			}
		}
	}, [filter, data]);

	const getRowId = useCallback((row) => {
		return row.id;
	}, []);

	// @ts-ignore
	const { getTableProps, headerGroups, rows, prepareRow, setSortBy, setHiddenColumns, setGroupBy } = useTable(
		{
			columns,
			data: tableData,
			defaultColumn,
			getRowId,
		},
		useResizeColumns,
		useFlexLayout,
		useGroupBy,
		useSortBy,
		useExpanded
	);

	const moveRow = useCallback(
		(dragIndex: number, hoverIndex: number) => {
			const dragRow = tableData[dragIndex];
			setTableData(
				update(tableData, {
					$splice: [
						[dragIndex, 1],
						[hoverIndex, 0, dragRow],
					],
				})
			);
		},
		[tableData]
	);

	useEffect(() => {
		if (sortBy) {
			setSortBy([sortBy])
		}
	}, [sortBy]);

	useEffect(() => {
		if (hiddenColumns) {
			setHiddenColumns(hiddenColumns)
		}
	}, [hiddenColumns]);

	useEffect(() => {
		if (groupBy) {
			setGroupBy([groupBy])
		}
	}, [groupBy]);

	const onScroll = useCallback(() => {
		if (menuProps.isOpen) {
			closeMenu();
		}
	}, [menuProps]);

	const itemData = createItemData(rowHeight, getRowId, moveRow, rows, prepareRow, selected, setSelected, actions, items, openMenu, closeMenu, menuProps, anchorPoint, setAnchorPoint, setContextMenuRow);

	return useMemo(() => (
		<DndProvider backend={HTML5Backend}>
			<ContextMenu {...{ row: contextMenuRow, addGroup, group, groups, items, closeMenu, setAnchorPoint, anchorPoint, menuProps }} />
			<Container {...getTableProps()}>
				<OverflowContainer>
					<TableHeader
						actions={actions}
						hotkeysEnabled={hotkeysEnabled}
						headerGroups={headerGroups}
						setSelected={setSelected}
						selected={selected}
						rows={rows}
					/>

					<TableBody>
						<AutoSizer>
							{throttle(({ height, width }) => (
								<List
									height={height}
									width={width}
									itemCount={rows.length}
									itemData={itemData}
									itemSize={rowHeight + GUTTER_SIZE}
									style={{ willChange: "auto" }}
									onScroll={onScroll}
								>
									{/* @ts-ignore */}
									{TableRow}
								</List>
							), 25)}
						</AutoSizer>
					</TableBody>
				</OverflowContainer>
			</Container>
		</DndProvider>
	), [itemData, group, groups]);
};

const Container = styled.div`
  margin: 16px 0 0 0;
  display: flex;
  flex-direction: column;
  flex: 1;
`;

const OverflowContainer = styled.div`
	display: flex;
	flex: 1;
	margin: 0 -12px 0 0;
  overflow: hidden;
  flex-direction: column;
`;

const TableBody = styled.div`
	position: relative;
	flex: 1 1 auto;
  width: 100%;
  margin: 0;
  display: flex;
  overflow: hidden;
`;

const List = styled(FixedSizeList)`
	overflow-y: scroll !important;
	will-change: auto;
`;

// the wrong will-change causes bad perf on windows

export default Table;
