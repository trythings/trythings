import color from 'color';
import React from 'react';

import resetStyles from './resetStyles.js';
import theme from './theme.js';

export default class FlatButton extends React.Component {
	static propTypes = {
		color: React.PropTypes.string.isRequired,
		label: React.PropTypes.string.isRequired,
		onClick: React.PropTypes.func,
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

	onBlur = (event) => {
		if (event.relatedTarget && !event.currentTarget.contains(event.relatedTarget)) {
			this.setState({ isFocused: false, isActive: false });
		}
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
			justifyContent: 'center',
			paddingLeft: 8,
			paddingRight: 8,
		},
		label: {
			...resetStyles,
			...theme.text,

			fontSize: 14,
			fontWeight: 500,
			textTransform: 'uppercase',
		},
	};

	buttonStateStyle() {
		if (this.state.isActive) {
			return {
				backgroundColor: color(this.props.color).alpha(0.38).rgbString(),
			};
		}
		if (this.state.isFocused) {
			return {
				backgroundColor: color(this.props.color).alpha(0.24).rgbString(),
			};
		}
		if (this.state.isHovered) {
			return {
				backgroundColor: color(this.props.color).alpha(0.12).rgbString(),
			};
		}
		return {};
	}

	render() {
		return (
			<button
				style={{
					...FlatButton.styles.button,
					...this.buttonStateStyle(),
					color: this.props.color,
				}}
				onClick={this.props.onClick}

				onMouseEnter={this.onMouseEnter}
				onMouseLeave={this.onMouseLeave}
				onFocus={this.onFocus}
				onBlur={this.onBlur}
				onMouseDown={this.onMouseDown}
				onMouseUp={this.onMouseUp}
			>
				<span
					style={{
						...FlatButton.styles.label,
						color: this.props.color,
					}}
				>
					{this.props.label}
				</span>
			</button>
		);
	}
}
