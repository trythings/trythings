var path = require('path');
var webpack = require('webpack');

module.exports = {
	devtool: 'source-map',
	entry: path.resolve(__dirname, 'src', 'index.js'),
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
			exclude: /node_modules/,
			loaders: 'eslint',
		}],
		loaders: [{
			test: /\.js$/,
			exclude: /node_modules/,
			loaders: ['babel'],
		}],
	},
	devServer: {
		proxy: {
			'/graphql': 'http://localhost:8080',
		},
	},
};
