import _pick from 'lodash/pick';
import color from 'color';
import React from 'react';

import Icon from './Icon.js';
import resetStyles from './resetStyles.js';

export default class IconButton extends React.Component {
	static propTypes = {
		iconName: React.PropTypes.string.isRequired,
		onClick: React.PropTypes.func,
		style: React.PropTypes.shape({
			color: React.PropTypes.string.isRequired,
		}).isRequired,
	};

	static styles = {
		button: {
			...resetStyles,

			borderRadius: '50%',
			paddingBottom: 8,
			paddingLeft: 8,
			paddingRight: 8,
			paddingTop: 8,
		},
	};

	state = {
		hasFocus: false,
		isActive: false,
		isHovering: false,
	};

	onBlur = (event) => {
		if (event.relatedTarget && !event.currentTarget.contains(event.relatedTarget)) {
			this.setState({ hasFocus: false, isActive: false });
		}
	};

	onFocus = () => {
		this.setState({ hasFocus: true });
	};

	onMouseEnter = () => {
		this.setState({ isHovering: true });
	};

	onMouseLeave = () => {
		this.setState({ isHovering: false, isActive: false });
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
					...IconButton.styles.button,
					...this.buttonStateStyle(),
				}}
				onBlur={this.onBlur}
				onClick={this.props.onClick}
				onFocus={this.onFocus}
				onMouseEnter={this.onMouseEnter}
				onMouseLeave={this.onMouseLeave}
			>
				<Icon
					name={this.props.iconName}
					style={_pick(this.props.style, 'color')}
				/>
			</button>
		);
	}
}
