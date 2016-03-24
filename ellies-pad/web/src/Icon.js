import React from 'react';

import './MaterialIcons.css';

import resetStyles from './resetStyles.js';

export default class Icon extends React.Component {
	static propTypes = {
		name: React.PropTypes.string.isRequired,
		style: React.PropTypes.shape({
			color: React.PropTypes.string,
		}),
	};

	static styles = {
		icon: {
			...resetStyles,
			cursor: 'default',
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
					color: this.props.style && this.props.style.color,
				}}
			>
				{this.props.name}
			</i>
		);
	}
}
