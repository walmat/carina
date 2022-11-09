import {Dispatch, SetStateAction} from "react";

export interface Items {
  name: string;
  actions: Action[];
}

export interface Action {
  name: string;
  Icon: any;
  shortcut?: string[];
  onClick: (id: string) => void;
}

export interface RowProps {
  index: number;
  style: any;
  data: {
    rowHeight: number;
    rows: any;
    moveRow: any;
    prepareRow: any;
    selected: any;
    setSelected: any;
    getRowId: any;
    actions: { name: string; onClick: any }[];
    items: Items[];
    openMenu: any;
    closeMenu: any;
    menuProps: any;
    anchorPoint: {x: number, y: number};
    setAnchorPoint: Dispatch<SetStateAction<{x: number, y: number}>>;
    setContextMenuRow: any,
  };
}

export interface IconComponentProps {
  id?: string;
  state?: number;
  name: string;
  onClick?: any;
}

export interface DragItem {
  index: number;
  type: string;
}
