const HtmlWebpackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');

const { merge } = require('webpack-merge');

const prodConfig = require('./webpack.prod.config')

module.exports = merge(prodConfig, {
	entry: {
		main: ['./src/roots/Main'],
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: './public/index.html',
		}),
		new CleanWebpackPlugin()
	],
});
