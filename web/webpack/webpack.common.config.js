const path = require("path");
const webpack = require("webpack");
const DotEnv = require("dotenv-webpack");

const ForkTsCheckerWebpackPlugin = require("fork-ts-checker-webpack-plugin");

const nodeFallbacks = require("./nodeFallbacks");

const cwd = process.cwd();
const outputPath = path.join(cwd, "build");

module.exports = {
  target: "web",
  mode: "development",
  devtool: "cheap-module-source-map",
  context: path.resolve(cwd, "./"),
  output: {
    path: outputPath,
    publicPath: "/",
    pathinfo: false,
    filename: "[name].js",
  },
  plugins: [
    new ForkTsCheckerWebpackPlugin(),
    new DotEnv({
      path: "./.env",
    }),
    new webpack.optimize.LimitChunkCountPlugin({
      maxChunks: 1,
    }),
    new webpack.ContextReplacementPlugin(/moment[\/\\]locale$/, /est/),
  ],
  experiments: {
    // enable if you feel it is right
    lazyCompilation: false,
  },
  resolve: {
    extensions: [".ts", ".tsx", ".js", ".json", ".mjs"],
    fallback: {
      ...nodeFallbacks,
    },
  },
  cache: {
    // type: 'memory',
    type: "filesystem",
    // persistent cache
    buildDependencies: {
      config: [__filename], // you may omit this when your CLI automatically adds it
    },
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx|ts|tsx)?$/,
        exclude: [/node_modules/],
        use: [
          // {
          //   loader: 'cache-loader',
          // },
          {
            loader: "babel-loader?cacheDirectory",
          },
        ],
      },
      {
        test: /\.(jpe?g|png|gif|svg|pdf|csv|xlsx|ttf|woff(2)?)$/i,
        use: [
          {
            loader: "file-loader",
            options: {
              name: "[name].[ext]",
              outputPath: "img/",
            },
          },
        ],
      },
      {
        test: /\.mp3$/,
        use: [
          {
            loader: "file-loader",
            options: {
              name: "[name].[ext]",
              outputPath: "sounds/",
            },
          },
        ],
      },
      {
        test: /\.css$/,
        use: ["style-loader", "css-loader"],
      },
    ],
  },
};
