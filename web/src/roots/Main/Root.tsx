import React, { useMemo } from 'react';
import { Store } from 'redux';
import { HashRouter as Router } from 'react-router-dom';

import Providers from '../../providers';
import App from './App';

type Props = {
  store: Store;
};

const Root = ({ store }: Props) => {
  return useMemo(() => (
    <Providers store={store}>
      <Router>
        <App />
      </Router>
    </Providers>
  ), []);
};

export default Root;
