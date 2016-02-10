import React from 'react';

import Icon from './Icon.js';
import theme from './theme.js';

export default class ActionButton extends React.Component {
	static propTypes = {
		onClick: React.PropTypes.func,
	};

	static styles = {
		button: {
			backgroundColor: theme.colors.accent,

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
		container: {
			padding: 24,
		},
	};

	render() {
		return (
			<div style={ActionButton.styles.container}>
				<button onClick={this.props.onClick} style={ActionButton.styles.button}>
					<Icon color={theme.text.light.primary} name="add"/>
				</button>
			</div>
		);
	}
}
