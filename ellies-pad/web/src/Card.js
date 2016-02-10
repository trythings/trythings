import React from 'react';

import theme from './theme.js';

export default class Card extends React.Component {
	static propTypes = {
		children: React.PropTypes.node,
	};

	static styles = {
		card: {
			borderRadius: 2,
			backgroundColor: theme.colors.card,
			boxShadow: [
				'0 1px 5px 0 rgba(0, 0, 0, 0.12)', // Ambient.
				'0 2px 2px 0 rgba(0, 0, 0, 0.14)', // Penumbra.
				'0 3px 1px -2px rgba(0, 0, 0, 0.20)', // Umbra.
			].join(','),
			zIndex: 2,

			boxSizing: 'border-box',
			display: 'flex',
			flexDirection: 'column',
			minHeight: 'min-content',
			minWidth: 'min-content',
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
