var path = require('path');
var webpack = require('webpack');

module.exports = {
	devtool: 'source-map',
	entry: [
		path.resolve(__dirname, 'src', 'index.js'),
	],
	output: {
		path: path.resolve(__dirname, 'static'),
		filename: 'bundle.js',
	},
	externals: {
		gapi: 'gapi',
	},
	plugins: [
		new webpack.EnvironmentPlugin(['NODE_ENV']),
		new webpack.NoErrorsPlugin(),
	],
	module: {
		preLoaders: [{
			test: /\.js$/,
			include: path.resolve(__dirname, 'src'),
			loaders: 'eslint',
		}],
		loaders: [
			{
				test: /\.js$/,
				include: path.resolve(__dirname, 'src'),
				loaders: ['babel'],
			},
			{
				test: /\.css$/,
				// TODO Figure out how to restrict this loader to only
				// src and node_modules/normalize.css
				// include: [
				// 	path.resolve(__dirname, 'src'),
				// 	path.resolve(__dirname, 'node_modules', 'normalize.css'),
				// ],
				loaders: ['style', 'css'],
			},
		],
	},
	devServer: {
		proxy: {
			'*': 'http://localhost:8080',
		},
	},
};
