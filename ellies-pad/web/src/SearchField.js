import _pick from 'lodash/pick';
import color from 'color';
import React from 'react';

import Icon from './Icon.js';
import resetStyles from './resetStyles.js';
import theme from './theme.js';

export default class SearchField extends React.Component {
	static propTypes = {
		autoFocus: React.PropTypes.bool,
		initialQuery: React.PropTypes.string,
		onFocus: React.PropTypes.func,
		onBlur: React.PropTypes.func,
		onQueryChange: React.PropTypes.func,
		style: React.PropTypes.shape({
			backgroundColor: React.PropTypes.string.isRequired,
			color: React.PropTypes.string.isRequired,
			flex: React.PropTypes.string,
		}).isRequired,
	};

	static styles = {
		container: {
			...resetStyles,
			alignItems: 'stretch',
			borderRadius: 2,
			cursor: 'text',
			paddingBottom: 4,

			paddingTop: 4,
		},
		spacer: {
			...resetStyles,
			paddingLeft: 16,
		},
		icon: {
			...resetStyles,
			...theme.text,
		},
		input: {
			...resetStyles,
			...theme.text,
			flex: '1 0 auto',
		},
	};

	constructor(props, ...args) {
		super(props, ...args);
		this.state = {
			query: props.initialQuery || '',
		};
	}

	state = {
		isHovered: false,
		hasFocus: false,
	};

	onChange = (event) => {
		const query = event.target.value;
		this.setState({ query });
		if (this.props.onQueryChange) {
			this.props.onQueryChange(query);
		}
	};

	onClick = () => {
		if (this.input) {
			this.input.focus();
		}
	};

	onFocus = (...args) => {
		if (this.props.onFocus) {
			this.props.onFocus(...args);
		}
		this.setState({ hasFocus: true });
	};

	onBlur = (...args) => {
		if (this.props.onBlur) {
			this.props.onBlur(...args);
		}
		this.setState({ hasFocus: false });
	};

	onMouseEnter = () => {
		this.setState({ isHovered: true });
	};

	onMouseLeave = () => {
		this.setState({ isHovered: false });
	};

	inputRef = (input) => {
		this.input = input;
	};

	render() {
		let backgroundColor = this.props.style.backgroundColor;
		if (this.state.isHovered || this.state.hasFocus) {
			// Bring the background color closer to the text color.
			backgroundColor = color(this.props.style.backgroundColor).
					mix(color(this.props.style.color), 1 - 0.12).
					hexString();
		}

		return (
			<div
				onClick={this.onClick}
				onMouseEnter={this.onMouseEnter}
				onMouseLeave={this.onMouseLeave}
				style={{
					...SearchField.styles.container,
					..._pick(this.props.style, ['flex']),
					backgroundColor,
				}}
			>
				<div style={SearchField.styles.spacer} />
				<Icon
					name="search"
					style={{
						...SearchField.styles.icon,
						color: this.props.style.color,
					}}
				/>
				<div style={SearchField.styles.spacer} />
				<input
					autoFocus={this.props.autoFocus}
					onChange={this.onChange}
					onFocus={this.onFocus}
					onBlur={this.onBlur}
					style={{
						...SearchField.styles.input,
						color: this.props.style.color,
					}}
					value={this.state.query}
					ref={this.inputRef}
				/>
			</div>
		);
	}
}
