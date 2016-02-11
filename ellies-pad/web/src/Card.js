import React from 'react';

import resetStyles from './resetStyles.js';
import theme from './theme.js';

export default class Card extends React.Component {
	static propTypes = {
		children: React.PropTypes.node,
	};

	static styles = {
		card: {
			...resetStyles,
			...theme.elevation[2],
			borderRadius: 2,
			backgroundColor: theme.colors.card,
			flexDirection: 'column',
		},
	};

	render() {
		return (
			<div style={Card.styles.card}>
				{this.props.children}
			</div>
		);
	}
}
