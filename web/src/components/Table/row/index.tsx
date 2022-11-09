import React, { MouseEvent, useCallback, useRef } from "react";
import { DropTargetMonitor, useDrag, useDrop } from "react-dnd";
import { Cell } from "react-table";
import styled from "styled-components";

import { ArrowDown, ArrowRight } from "react-feather";

import { RowProps, DragItem } from "../types";

import { IconComponent, maskInput, renderArrayCell } from '../utils';
import TaskStatusCell from '../TaskStatusCell';

const DND_ITEM_TYPE = 'row';

const TableRow = ({ index, style, data }: RowProps) => {
	const {
		rowHeight,
		rows,
		moveRow,
		prepareRow,
		selected,
		setSelected,
		actions,
		openMenu,
		setAnchorPoint,
		setContextMenuRow,
	} = data;

	const row = rows[index];
	prepareRow(row);

	const isSelected = selected?.includes(row.id);
	const onRowClick = useCallback((e: MouseEvent) => {
		if (setSelected) {
			if (!isSelected) {
				if (e.shiftKey || e.ctrlKey || e.altKey) {
					setSelected((prev: string[]) => [...prev, row.id]);
				} else {
					setSelected([row.id]);
				}
			} else {
				if (e.shiftKey || e.ctrlKey || e.altKey) {
					setSelected((prev: string[]) => [...prev.filter((id) => id !== row.id)]);
				} else {
					setSelected([]);
				}
			}
		}
	}, [row.id, isSelected]);

	const dropRef = useRef(null);
	const dragRef = useRef(null);

	const [, drop] = useDrop({
		accept: DND_ITEM_TYPE,
		hover(item: DragItem, monitor: DropTargetMonitor) {
			if (!dropRef.current) {
				return;
			}
			const dragIndex = item.index;
			const hoverIndex = index;
			if (dragIndex === hoverIndex) {
				return;
			}
			// @ts-ignore
			const hoverBoundingRect = dropRef.current?.getBoundingClientRect();
			const hoverMiddleY =
				(hoverBoundingRect.bottom - hoverBoundingRect.top) / 2;
			const clientOffset = monitor.getClientOffset();
			// @ts-ignore
			const hoverClientY = clientOffset.y - hoverBoundingRect.top;
			if (dragIndex < hoverIndex && hoverClientY < hoverMiddleY) {
				return;
			}
			if (dragIndex > hoverIndex && hoverClientY > hoverMiddleY) {
				return;
			}
			moveRow(dragIndex, hoverIndex);
			item.index = hoverIndex;
		},
	});

	const [{ isDragging }, drag, preview] = useDrag({
		type: DND_ITEM_TYPE,
		item: { index },
		collect: (monitor: any) => ({
			isDragging: monitor.isDragging(),
		}),
	});

	const opacity = isDragging ? 0 : 1;
	preview(drop(dropRef));
	drag(dragRef);

	return (
	<div style={style}>
		<Container opacity={opacity} ref={dropRef} {...row.getRowProps()}>
			<TableRowData
				onContextMenu={e => {
					e.preventDefault();

					if (!isSelected) {
						if (e.shiftKey || e.ctrlKey || e.altKey) {
							setSelected((prev: string[]) => [...prev, row.id]);
						} else {
							setSelected([row.id]);
						}
					}

					// set the row that enabled context menu
					setContextMenuRow(row);

					setAnchorPoint({ x: e.clientX, y: e.clientY });
					openMenu();
				}}
				rowHeight={rowHeight}
				isSelected={isSelected}
			>
				<DragRowCell ref={dragRef}>
					<IconComponent name="Drag" />
				</DragRowCell>
				{row.cells.map((cell: Cell) => {
					/* @ts-ignore */
					if (cell.column.type === 'password') {
						return (
							/* @ts-ignore */
							<RowCell onClick={!cell.isGrouped ? onRowClick : () => {}}  cursor="pointer" usePadding {...cell.getCellProps()}>
								{cell.value === null ? 'None' : maskInput(cell.value.length)}
							</RowCell>
						)
					}

					/* @ts-ignore */
					if (cell.column.type === 'card') {
						return (
							/* @ts-ignore */
							<RowCell onClick={!cell.isGrouped ? onRowClick : () => {}} cursor="pointer" usePadding {...cell.getCellProps()}>
								{`${maskInput(cell.value.length - 4, true)} ${cell.value.slice(cell.value.length - 4, cell.value.length)}`}
							</RowCell>
						)
					}

					if (cell.column.id === 'id') {
						return (
							// @ts-ignore
							<RowCell onClick={!cell.isGrouped ? onRowClick : () => {}} cursor="pointer" usePadding {...cell.getCellProps()}>
								{(cell?.value || '').toString().slice(0, 5)}
							</RowCell>
						);
					}

					if (cell.column.id === 'status') {
						return (
							<TaskStatusCell
								id={row.original.id}
								cell={cell}
								onRowClick={onRowClick}
								/>
						);
					}

					return (
						/* @ts-ignore */
						<RowCell onClick={!cell.isGrouped ? onRowClick : () => {}} cursor="pointer" usePadding {...cell.getCellProps()}>
							{/* @ts-ignore */}
							{cell.isGrouped ? (
									<span {...row.getToggleRowExpandedProps()}>
								{cell.render("Cell")} ({row.subRows.length})
										{row.isExpanded ? <ArrowDown width={14} height={14} /> : <ArrowRight width={14} height={14} />}
									</span>
								) : /* @ts-ignore */
								cell.isAggregated ? (
										`${row.subRows.length} ${row.subRows.length === 1 ? "Value" : "Values"}`
									) : /* @ts-ignore */
									cell.isPlaceholder ? null : cell.value === null
										? 'None' :
										Array.isArray(cell.value)
											? renderArrayCell(cell.value)
											:
											typeof cell.value === 'string'
												? cell.render("Cell") :
												typeof cell.value === 'boolean' || typeof cell.value === 'number'
													? `${cell?.value}` :
													(
														<>
															<span>{cell?.value?.name}</span>
															<span>{cell?.value?.details}</span>
														</>
													)
							}
						</RowCell>
					)})}
				<RowCellColumn>
					{actions.map(props => <IconComponent {...props} state={row.original.state} id={row.id} />)}
				</RowCellColumn>
			</TableRowData>
		</Container>
	</div>
	)
};

const Container = styled.div<{ isDragging: boolean }>`
	display: flex;
	flex: 1;
	margin: 4px 4px 4px 0;
  padding: 4px 0;
  background-color: ${({ theme, isDragging }) => isDragging ? theme.colors.secondary : theme.colors.background};
`;

const TableRowData = styled.div<{ rowHeight: number; isSelected: boolean }>`
	display: flex;
	flex: 1;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 400;
  padding: 0 8px 0 0;
  height: ${({ rowHeight }) => `${rowHeight}`}px;
  background-color: ${({ theme, isSelected }) => (isSelected ? theme.colors.secondary : theme.colors.sidebar)};
  color: ${({ theme }) => theme.colors.paragraph};

  & > * {
    cursor: pointer;
  }
`;

export const RowCell = styled.div<{ usePadding: boolean; color?: string; cursor: string }>`
  padding: ${({ usePadding }) => !usePadding ? '0 0 0 4px' : '0 4px' };
  display: flex;
  overflow: hidden;
  text-overflow: ellipsis;
  flex-direction: column;
  justify-content: center;
  color: ${({ theme, color }) => color ? color : theme.colors.table};
	cursor: ${({ cursor }) => cursor};

  & > span {
    cursor: ${({ cursor }) => cursor};
    display: flex;
    align-items: center;
    & svg {
      margin-left: 10px;
    }
  }

  & > span:nth-child(2) {
    margin-top: 2px;
    font-size: 10px;
    font-weight: 400;
    color: ${({ theme }) => theme.colors.paragraph};
  }
`;

const DragRowCell = styled.div`
  padding: 0 4px 0 6px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  color: ${({ theme }) => theme.colors.table};
	cursor: grab;

	& > * {
		cursor: grab;
	}
`;

const RowCellColumn = styled.div`
  padding: 0 4px;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  color: ${({ theme }) => theme.colors.table};
`;

export default TableRow;
