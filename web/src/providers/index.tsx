import React, { ReactElement, useMemo, Fragment } from "react";
import { Provider, useSelector } from "react-redux";
import { Store } from "redux";
import { HashRouter } from "react-router-dom";
import { ThemeProvider } from "styled-components";
import { Toaster } from "react-hot-toast";
import { RelayEnvironmentProvider } from "react-relay";
import { IntercomProvider } from 'react-use-intercom';

import { themes, GlobalStyles } from "../styles";
import { makeTheme } from "../stores/Main/reducers/theme";
import Environment from "../relay/Environment";
import { RecoilRoot } from "recoil";

interface ProviderProps {
  store?: Store;
  children: React.ReactElement;
}

interface ThemeProviderProps {
  children: React.ReactElement;
}

const getTheme = (index: number) => {
  return themes[index];
};

const ThemeInjector = ({ children }: ThemeProviderProps) => {
  const theme = useSelector(makeTheme);
  const chosen = useMemo(() => getTheme(theme), [theme]);

  return <ThemeProvider theme={chosen}>{children}</ThemeProvider>;
};

const Providers = ({ store, children }: ProviderProps): ReactElement => {

  if (typeof store !== 'undefined') {
    return (
      <Provider store={store}>
        {/* @ts-ignore*/}
        <RelayEnvironmentProvider environment={Environment}>
          <RecoilRoot>
            <ThemeInjector>
              <IntercomProvider appId="pv2jwmgn">
                <Fragment>
                  <Toaster
                    toastOptions={{
                      success: {
                        iconTheme: {
                          primary: '#786EF2',
                          secondary: '#fff',
                        },
                      },
                      loading: {
                        iconTheme: {
                          primary: '#FFC857',
                          secondary: '#fff',
                        },
                      },
                      error: {
                        iconTheme: {
                          primary: '#F26E86',
                          secondary: '#fff',
                        },
                      },
                    }}
                  />
                  <GlobalStyles />
                  <HashRouter>{children}</HashRouter>
                </Fragment>
              </IntercomProvider>
            </ThemeInjector>
          </RecoilRoot>
        </RelayEnvironmentProvider>
      </Provider>
    );
  }

  return (
    // @ts-ignore
    <RelayEnvironmentProvider environment={Environment}>
      <RecoilRoot>
        <IntercomProvider appId="pv2jwmgn">
          <Fragment>
            <Toaster />
            <GlobalStyles />
            <HashRouter>{children}</HashRouter>
          </Fragment>
        </IntercomProvider>
      </RecoilRoot>
    </RelayEnvironmentProvider>
  );
};

export default Providers;
