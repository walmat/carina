import React, { Dispatch, SetStateAction, useCallback } from "react";
import { useDispatch } from "react-redux";
import styled from "styled-components";
import { AnimatePresence, motion } from "framer-motion";
import { ChevronUp } from "react-feather";
import { useHotkeys } from "react-hotkeys-hook";

import { Buttons } from "../../../elements";

import { removeTask, startTask, stopTask } from '../../../stores/Main/reducers/tasks';

interface ActionBarProps {
  selected: string[];
  group: string;
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

const menu = {
  open: {
    y: 0,
    marginTop: '8px',
    height: '64px',
    opacity: 1,
  },
  closed: {
    y: 32,
    marginTop: '-24px',
    height: '16px',
    opacity: 0.6,
  },
};

const cancel = (e: any) => e.stopPropagation();

const Actionbar = ({ selected, group, open, setOpen }: ActionBarProps) => {
  const dispatch = useDispatch();

  const start = useCallback(
    (e) => {
      cancel(e);

      dispatch(startTask({ ids: selected, group }));
    },
    [selected]
  );
  const stop = useCallback(
    (e) => {
      cancel(e);

      dispatch(stopTask({ ids: selected, group }));
    },
    [selected]
  );
  const remove = useCallback(
    (e) => {
      cancel(e);

      dispatch(removeTask({ ids: selected, group }));
    },
    [selected]
  );

  useHotkeys("command+r,ctrl+r", start, [selected]);
  useHotkeys("command+s,ctrl+s", stop, [selected]);
  useHotkeys("command+d,ctrl+d", remove, [selected]);

  return (
    <Container
      initial={open}
      drag="y"
      whileDrag={{ cursor: "drag" }}
      dragConstraints={{ top: 0, bottom: 0 }}
      whileHover={{ y: open ? 30 : 0, opacity: 1 }}
      dragElastic={0.15}
      animate={open ? "closed" : "open"}
      variants={menu}
      onClick={() => setOpen((prev) => !prev)}
    >
      <AnimatePresence exitBeforeEnter>
        {!open && (
          <Presence
            key="drawer-container"
            initial={{ opacity: 1 }}
            animate={{ opacity: open ? 0 : 1 }}
            transition={{ opacity: { duration: 0.15 } }}
            exit={{ opacity: 0 }}
          >
            <Col>
              <Buttons.Primary
                variant="IconButton"
                text="Start"
                command="⌘ R"
                width={144}
                height={36}
                onClick={start}
              />
            </Col>
            <Col>
              <Buttons.Secondary
                variant="IconButton"
                text="Stop"
                command="⌘ S"
                width={128}
                height={36}
                onClick={stop}
              />
            </Col>
            <Col shove>
              <Buttons.Tertiary
                variant="IconButton"
                text="Remove"
                command="⌘ D"
                width={144}
                height={36}
                onClick={remove}
              />
            </Col>
          </Presence>
        )}

        {open && (
          <Row
            initial={{ opacity: 0 }}
            animate={{ opacity: !open ? 0 : 1 }}
            transition={{ opacity: { duration: 0.15 } }}
            exit={{ opacity: 0 }}
          >
            <Expand />
          </Row>
        )}
      </AnimatePresence>
    </Container>
  );
};

const Container = styled(motion.div)`
  display: flex;
  align-items: center;
  border-radius: 4px;
  background-color: ${({ theme }) => theme.colors.sidebar};
`;

const Presence = styled(motion.div)`
  display: flex;
  flex: 1;
`;

const Col = styled.div<{ shove?: boolean; center?: boolean }>`
  display: flex;
  flex-direction: column;
  margin: 0 8px;
  ${({ center }) =>
    center
      ? `
		align-items: center;
		justify-content: center;
	`
      : ""}
  ${({ shove }) => (shove ? "margin-left: auto;" : "")}
`;

const Row = styled(motion.div)<{ color?: string; accent?: string }>`
  display: flex;
  flex: 1;
  cursor: pointer;
  align-items: center;
  justify-content: center;
`;

const Expand = styled(ChevronUp)`
  width: 14px;
  height: 14px;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.h2};
`;

export default Actionbar;
