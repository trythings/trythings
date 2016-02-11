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

			backgroundColor: theme.colors.accent,

			alignItems: 'center',
			justifyContent: 'center',

			minWidth: 56,
			width: 56,
			minHeight: 56,
			height: 56,
			borderRadius: '50%',
		},
	};

	render() {
		return (
			<button onClick={this.props.onClick} style={ActionButton.styles.button}>
				<Icon color={theme.text.light.primary.color} name="add"/>
			</button>
		);
	}
}
