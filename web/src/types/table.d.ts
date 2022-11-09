export interface TableItems {
  name: string;
  actions: TableAction[];
}

export interface TableAction {
  name: string;
  Icon: any;
  shortcut?: string[];
  onClick: (id: string) => void;
}
