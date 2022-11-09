const HtmlWebpackPlugin = require('html-webpack-plugin');

const { merge } = require('webpack-merge');

const prodConfig = require('./webpack.prod.config')

module.exports = merge(prodConfig, {
	entry: {
		login: ['./src/roots/Login'],
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: './public/login.html',
			filename: 'login.html',
		})
	],
});
