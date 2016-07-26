import React from 'react';

import resetStyles from './resetStyles.js';
import theme from './theme.js';

class App extends React.Component {
	static propTypes = {
		children: React.PropTypes.node,
	};

	static styles = {
		app: {
			...resetStyles,
			alignItems: 'stretch',
			backgroundColor: theme.colors.canvas,
			height: '100%',
			width: '100%',
		},
	};

	render() {
		return (
			<div style={App.styles.app} tabIndex={-1}>
				{this.props.children}
			</div>
		);
	}
}

export default App;
