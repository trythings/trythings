import React from 'react';

import Icon from './Icon.js';
import resetStyles from './resetStyles.js';
import theme from './theme.js';

export default class ActionButton extends React.Component {
	static propTypes = {
		onClick: React.PropTypes.func,
	};

	static styles = {
		button: {
			...resetStyles,
			...theme.elevation[6],

			alignItems: 'center',
			display: 'flex',
			justifyContent: 'center',

			backgroundColor: theme.colors.accent,
			borderRadius: '50%',
			height: 56,
			width: 56,
		},
		icon: {
			...resetStyles,
			...theme.text.light.primary,
		},
	};

	render() {
		return (
			<button onClick={this.props.onClick} style={ActionButton.styles.button}>
				<Icon style={ActionButton.styles.icon} name="add" />
			</button>
		);
	}
}
