import React, { Dispatch, SetStateAction, Fragment, useCallback } from "react";

import Modal from "../Modal";
import Tabs from "../Tabs";

import AddGroup from "./AddGroup";
import RemoveGroup from "./RemoveGroup";

export interface GroupProps {
  tabs?: {
    titles: string[];
    contents: any;
  };
  title: string;
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export interface RemoveGroupForm {
  name: string;
}

const Groups = ({ tabs, title, open, setOpen }: GroupProps) => {
  const handleClose = useCallback(() => {
    return setOpen(false);
  }, []);

  return (
    <Modal
      height="auto"
      width={325}
      title={title}
      show={open}
      setShow={handleClose}
    >
      <Fragment key="Groups">{tabs ? <Tabs {...tabs} /> : null}</Fragment>
    </Modal>
  );
};

Groups.AddGroup = AddGroup;
Groups.RemoveGroup = RemoveGroup;

export default Groups;
