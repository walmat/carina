import { useCallback } from "react";
import { useDispatch } from "react-redux";
import styled from "styled-components";
import { X } from "react-feather";

import { focus, remove } from "../../../stores/Collective/reducers/harvester";
import { Harvester as HarvesterType } from "../types";

const Harvester = (harvester: HarvesterType) => {
  const { id, index, name, ...props } = harvester;

  const dispatch = useDispatch();

  const onSelect = useCallback(() => {
    if (props.focused) {
      return;
    }

    dispatch(focus(id));
  }, [props.focused, id]);

  const handleRemove = useCallback(() => {
    dispatch(remove(id));
  }, [id]);

  return (
    <Container onClick={onSelect} {...props}>
      {name}
      {props.focused && id !== "default" ? (
        <Remove onClick={handleRemove} />
      ) : null}
    </Container>
  );
};

const Container = styled.div<{ focused: boolean }>`
  display: flex;
  justify-content: center;
  align-items: center;
  text-transform: capitalize;
  white-space: pre;
  opacity: ${({ focused }) => (focused ? 1 : 0.6)};
  background: ${({ theme }) => theme.colors.primary};
  color: #fff;
  padding: 4px 10px;
  margin: 0 4px;
  font-size: 12px;
  font-weight: ${({ focused }) => (focused ? 500 : 400)};
  cursor: ${({ focused }) => (focused ? "default" : "pointer")};
  border-radius: 3.5px;
  flex: 1;
`;

const Remove = styled(X)`
  color: #fff;
  width: 14px;
  height: auto;
  margin: 0 0 0 6px;
  cursor: pointer;
`;

export default Harvester;
