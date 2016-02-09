import React from 'react';

import './MaterialIcons.css';

import theme from './theme.js';

export default class Icon extends React.Component {
	static propTypes = {
		name: React.PropTypes.string.isRequired,
	};

	static styles = {
		icon: {
			color: theme.text.light.color,
			opacity: theme.text.light.opacity.primary,
		},
	};

	render() {
		return <i
			className="material-icons"
			style={Icon.styles.icon}
		>{this.props.name}</i>;
	}
}
