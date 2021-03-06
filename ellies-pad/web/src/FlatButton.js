import color from 'color';
import React from 'react';

import resetStyles from './resetStyles.js';
import theme from './theme.js';

export default class FlatButton extends React.Component {
	static propTypes = {
		label: React.PropTypes.string.isRequired,
		onClick: React.PropTypes.func,
		style: React.PropTypes.shape({
			color: React.PropTypes.string.isRequired,
		}).isRequired,
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

	state = {
		hasFocus: false,
		isHovering: false,
		isActive: false,
	};

	onMouseEnter = () => {
		this.setState({ isHovering: true });
	};

	onMouseLeave = () => {
		this.setState({ isHovering: false, isActive: false });
	};

	onFocus = () => {
		this.setState({ hasFocus: true });
	};

	onBlur = (event) => {
		if (event.relatedTarget && !event.currentTarget.contains(event.relatedTarget)) {
			this.setState({ hasFocus: false, isActive: false });
		}
	};

	onMouseDown = () => {
		this.setState({ isActive: true });
	};

	onMouseUp = () => {
		this.setState({ isActive: false });
	};

	buttonStateStyle() {
		if (this.state.isActive) {
			return {
				backgroundColor: color(this.props.style.color).alpha(0.38).rgbString(),
			};
		}
		if (this.state.hasFocus) {
			return {
				backgroundColor: color(this.props.style.color).alpha(0.24).rgbString(),
			};
		}
		if (this.state.isHovering) {
			return {
				backgroundColor: color(this.props.style.color).alpha(0.12).rgbString(),
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
						color: this.props.style.color,
					}}
				>
					{this.props.label}
				</span>
			</button>
		);
	}
}
