const HtmlWebpackPlugin = require('html-webpack-plugin');

const { merge } = require('webpack-merge');

const prodConfig = require('./webpack.prod.config')

module.exports = merge(prodConfig, {
	entry: {
		collective: ['./src/roots/Collective'],
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: './public/collective.html',
			filename: 'collective.html',
		})
	],

});
