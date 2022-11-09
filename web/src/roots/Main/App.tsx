import React, { useMemo } from "react";
import styled from "styled-components";

import { useTaskStatus } from "../../hooks";
import Dragbar from "../../components/Dragbar";
import Sidebar from "../../containers/Sidebar";
import Routes from "./Routes";

import "../../i18n/config";
import { notify } from "../../utils";
import { useDispatch, useStore } from "react-redux";
import { makeUpdate, wsUpdate } from "../../stores/Main/reducers/update";

const App = () => {
  useTaskStatus();
  const dispatch = useDispatch();

  const store = useStore();
  const state = store.getState();
  let updateData = makeUpdate(state);

  if (state.Update) {
    setTimeout(() => {
      dispatch(wsUpdate(updateData));
      notify("A new version of Carina is available.", "success", {
        duration: 5000,
      });
    }, 500);
  }

  if (!(window as any).HasUpdateListener) {
    (window as any).HasUpdateListener = true;

    window.addEventListener("autoupdater", (evt: any) => {
      dispatch(wsUpdate(updateData));
      notify(evt.detail, "success", { duration: 5000 });
    });
  }

  return useMemo(
    () => (
      <AppContainer>
        <Dragbar />
        <Sidebar />
        <Routes />
      </AppContainer>
    ),
    []
  );
};

const AppContainer = styled.div`
  display: flex;
  flex: 1 0 0;
  min-width: 0;
  background-color: ${({ theme }) => theme.colors.background};
`;

export default App;
