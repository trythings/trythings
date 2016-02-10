import color from 'color';
import React from 'react';

import resetStyles from './resetStyles.js';

export default class FlatButton extends React.Component {
	static propTypes = {
		onClick: React.PropTypes.func,
		children: React.PropTypes.node.isRequired,
		color: React.PropTypes.string.isRequired,
	};

	state = {
		isFocused: false,
		isHovered: false,
		isActive: false,
	};

	onMouseEnter = () => {
		this.setState({ isHovered: true });
	};

	onMouseLeave = () => {
		this.setState({ isHovered: false, isActive: false });
	};

	onFocus = () => {
		this.setState({ isFocused: true });
	};

	onBlur = () => {
		this.setState({ isFocused: false, isActive: false });
	};

	onMouseDown = () => {
		this.setState({ isActive: true });
	};

	onMouseUp = () => {
		this.setState({ isActive: false });
	};

	static styles = {
		button: {
			...resetStyles,

			borderRadius: 2,
			height: 36,
			minWidth: 64,
			paddingLeft: 8,
			paddingRight: 8,

			fontSize: 14,
			fontWeight: 500,
			textTransform: 'uppercase',
		},
	};

	render() {
		let style = {
			...FlatButton.styles.button,
			color: this.props.color,
		};

		if (this.state.isActive) {
			style = {
				...style,
				backgroundColor: color(this.props.color).alpha(0.38).rgbString(),
			};
		} else if (this.state.isFocused) {
			style = {
				...style,
				backgroundColor: color(this.props.color).alpha(0.24).rgbString(),
			};
		} else if (this.state.isHovered) {
			style = {
				...style,
				backgroundColor: color(this.props.color).alpha(0.12).rgbString(),
			};
		}

		return (
			<button
				style={style}
				onClick={this.props.onClick}

				onMouseEnter={this.onMouseEnter}
				onMouseLeave={this.onMouseLeave}
				onFocus={this.onFocus}
				onBlur={this.onBlur}
				onMouseDown={this.onMouseDown}
				onMouseUp={this.onMouseUp}
			>
				{this.props.children}
			</button>
		);
	}
}
