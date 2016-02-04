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
	plugins: [
		new webpack.EnvironmentPlugin(['NODE_ENV']),
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
				include: path.resolve(__dirname, 'src'),
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
