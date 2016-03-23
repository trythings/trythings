import React from 'react';

import Icon from './Icon.js';
import resetStyles from './resetStyles.js';
import theme from './theme.js';

export default class SearchField extends React.Component {
	static propTypes = {
		initialQuery: React.PropTypes.string,
		onQueryChange: React.PropTypes.func,
	};

	static styles = {
		container: {
			...resetStyles,
		},
		icon: {
			...resetStyles,
			color: theme.text.light.primary.color,
		},
	};

	constructor(props, ...args) {
		super(props, ...args);
		this.state = {
			query: props.initialQuery || '',
		};
	}

	onChange = (event) => {
		const query = event.target.value;
		this.setState({ query });
		if (this.props.onQueryChange) {
			this.props.onQueryChange(query);
		}
	};

	render() {
		return (
			<div style={SearchField.styles.container}>
				<Icon name="search" style={SearchField.styles.icon} />
				<input onChange={this.onChange} value={this.state.query} />
			</div>
		);
	}
}
