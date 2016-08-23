import _pick from 'lodash/pick';
import color from 'color';
import React from 'react';

import Icon from './Icon.js';
import resetStyles from './resetStyles.js';
// import theme from './theme.js';

export default class IconButton extends React.Component {
	static propTypes = {
		iconName: React.PropTypes.string.isRequired,
		onClick: React.PropTypes.func,
		style: React.PropTypes.shape({
			backgroundColor: React.PropTypes.string,
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
		isHovering: false,
		hasFocus: false,
	};

	onBlur = () => {
		this.setState({ hasFocus: false });
	};

	onFocus = () => {
		this.setState({ hasFocus: true });
	};

	onMouseEnter = () => {
		this.setState({ isHovering: true });
	};

	onMouseLeave = () => {
		this.setState({ isHovering: false });
	};

	render() {
		let backgroundColor = this.props.style.backgroundColor ||
			color(this.props.style.color).alpha(0).rgbString();
		if (this.state.isHovering || this.state.hasFocus) {
			backgroundColor = color(backgroundColor).
					mix(color(this.props.style.color), 1 - 0.12).
					rgbString();
		}

		return (
			<button
				style={{
					...IconButton.styles.button,
					backgroundColor,
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
