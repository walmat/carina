import { useRecoilValue } from 'recoil';
import { taskStatusAtomFamily } from '../../stores/Main/reducers/tasksAtom';
import React from 'react';
import { RowCell } from './row';

type Props = {
  id: string;
  cell: any;
  onRowClick: any;
};
const TaskStatusCell = ({ id, cell, onRowClick }: Props) => {
  const statusColor = useRecoilValue(taskStatusAtomFamily(id));

  return (
    <RowCell
      color={statusColor.color}
      onClick={!cell.isGrouped ? onRowClick : () => {}}
      cursor="pointer"
      usePadding
      {...cell.getCellProps()}
    >
      {statusColor.status}
    </RowCell>
  );
};

export default TaskStatusCell;
