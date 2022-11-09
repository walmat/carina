const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");

const { merge } = require("webpack-merge");

const devConfig = require("./webpack.dev.config");

const cwd = process.cwd();
const outputPath = path.join(cwd, "build");

module.exports = merge(devConfig, {
  entry: {
    login: ["./src/roots/Login"],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: "./public/login.html",
    }),
  ],
  devServer: {
    static: {
      directory: outputPath,
    },
    allowedHosts: 'all',
    historyApiFallback: {
      disableDotRule: true,
    },
    hot: true,
    compress: true,
    port: '3001',
    client: {
      overlay: {
        errors: true,
        warnings: false,
      },
    },
  },
  cache: {
    type: 'filesystem',
    allowCollectingMemory: true,
  },
  stats: { warnings: false },
});
