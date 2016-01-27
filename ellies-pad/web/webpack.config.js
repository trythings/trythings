var path = require('path');

module.exports = {
	devtool: 'source-map',
	entry: path.resolve(__dirname, 'src', 'index.js'),
	output: {
		path: path.resolve(__dirname, 'static'),
		publicPath: '/static/',
		filename: 'bundle.js',
	},
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
};
