import React from 'react';

import Icon from './Icon.js';
import theme from './theme.js';

export default class ActionButton extends React.Component {
	static styles = {
		button: {
			backgroundColor: theme.colors.accent1,

			border: 'none',
			outline: 0,

			display: 'flex',
			alignItems: 'center',
			justifyContent: 'center',

			minWidth: 56,
			width: 56,
			minHeight: 56,
			height: 56,
			borderRadius: '50%',

			boxShadow: [
				'0 1px 18px 0 rgba(0, 0, 0, 0.12)', // Ambient.
				'0 6px 10px 0 rgba(0, 0, 0, 0.14)', // Penumbra.
				'0 3px 5px -1px rgba(0, 0, 0, 0.20)', // Umbra.
			].join(','),
			zIndex: 6,
		},
		icon: {
			color: theme.text.light.color,
			opacity: theme.text.light.opacity.primary,
			fontSize: 24,
		},
		padding: {
			padding: 24,
		},
	};

	render() {
		return (
			<div style={ActionButton.styles.padding}>
				<button style={ActionButton.styles.button}>
					<Icon style={ActionButton.styles.icon} name="add"/>
				</button>
			</div>
		);
	}
}
