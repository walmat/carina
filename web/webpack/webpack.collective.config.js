const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

const { merge } = require('webpack-merge');

const devConfig = require('./webpack.dev.config')

const cwd = process.cwd();
const outputPath = path.join(cwd, 'build');

module.exports = merge(devConfig, {
	entry: {
		collective: ['./src/roots/Collective'],
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: './public/collective.html',
		})
	],
	devServer: {
		allowedHosts: 'all',
		static: {
			directory: outputPath,
		},
		historyApiFallback: {
			disableDotRule: true,
		},
		hot: true,
		compress: true,
		port: '3002',
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
