import React from 'react';

import './MaterialIcons.css';

import resetStyles from './resetStyles.js';

export default class Icon extends React.Component {
	static propTypes = {
		name: React.PropTypes.string.isRequired,
		color: React.PropTypes.string,
	};

	static styles = {
		icon: {
			...resetStyles,
			fontFamily: 'default',
			fontSize: 'default',
		},
	};

	render() {
		return (
			<i
				className="material-icons"
				style={{
					...Icon.styles.icon,
					color: this.props.color,
				}}
			>
				{this.props.name}
			</i>
		);
	}
}
