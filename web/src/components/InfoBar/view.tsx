import React, { Dispatch, SetStateAction } from "react";
import styled from "styled-components";

type ViewProps = {
  view: string;
  setView: Dispatch<SetStateAction<string>>;
};

const View = ({ view, setView }: ViewProps) => {
  return (
    <Container>
      <Button active={view === "Basic"} onClick={() => setView("Basic")}>
        Basic
      </Button>
      <Button active={view === "Advanced"} onClick={() => setView("Advanced")}>
        Advanced
      </Button>
    </Container>
  );
};

const Container = styled.div`
  min-width: 160.672px;
  display: flex;
  align-items: center;
  margin-right: 16px;
`;

const Button = styled.div<{ active: boolean }>`
  height: 32px;
  font-size: 14px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 0 16px;
  font-weight: ${({ active }) => (active ? 700 : 400)};
  opacity: ${({ active }) => (active ? 1 : 0.6)};
  background: ${({ active, theme }) =>
    active ? theme.colors.view : "transparent"};
  color: ${({ theme }) => theme.colors.muted};
`;

export default View;
