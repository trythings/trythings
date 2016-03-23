import React from 'react';

import Icon from './Icon.js';
import resetStyles from './resetStyles.js';
import theme from './theme.js';

export default class SearchField extends React.Component {
	static propTypes = {
		initialQuery: React.PropTypes.string,
		onQueryChange: React.PropTypes.func,
		style: React.PropTypes.shape({
			flex: React.PropTypes.string,
		}),
	};

	static styles = {
		container: {
			...resetStyles,
			alignItems: 'stretch',
			backgroundColor: theme.colors.primary.light,
			borderRadius: 2,
			paddingBottom: 4,

			paddingTop: 4,
		},
		spacer: {
			...resetStyles,
			paddingLeft: 16,
		},
		icon: {
			...resetStyles,
			...theme.text.light.primary,
		},
		input: {
			...resetStyles,
			...theme.text.light.primary,
			flex: '1 0 auto',
			fontWeight: 300,
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
	};

	onChange = (event) => {
		const query = event.target.value;
		this.setState({ query });
		if (this.props.onQueryChange) {
			this.props.onQueryChange(query);
		}
	};

	onMouseEnter = () => {
		this.setState({ isHovered: true });
	};

	onMouseLeave = () => {
		this.setState({ isHovered: false });
	};

	render() {
		let style = SearchField.styles.container;
		if (this.props.style && this.props.style.flex) {
			style = {
				...style,
				flex: this.props.style.flex,
			};
		}

		if (this.state.isHovered) {
			style = {
				...style,
				backgroundColor: theme.colors.primary.xlight,
			};
		}

		return (
			<div
				onMouseEnter={this.onMouseEnter}
				onMouseLeave={this.onMouseLeave}
				style={style}
			>
				<div style={SearchField.styles.spacer} />
				<Icon name="search" style={SearchField.styles.icon} />
				<div style={SearchField.styles.spacer} />
				<input
					onChange={this.onChange}
					style={SearchField.styles.input}
					value={this.state.query}
				/>
			</div>
		);
	}
}
