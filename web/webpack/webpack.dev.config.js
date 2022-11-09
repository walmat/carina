const ReactRefreshPlugin = require('@pmmmwh/react-refresh-webpack-plugin');

const { merge } = require('webpack-merge');
const common = require('./webpack.common.config');

module.exports = merge(common, {
  mode: 'development',
  plugins: [
    new ReactRefreshPlugin(),
  ],
});
