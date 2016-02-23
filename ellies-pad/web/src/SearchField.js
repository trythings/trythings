import React from 'react';

import Icon from './Icon.js';

export default class SearchField extends React.Component {
	static propTypes = {
		initialQuery: React.PropTypes.string,
		onQueryChange: React.PropTypes.func,
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
			<div>
				<Icon name="search"/>
				<input onChange={this.onChange} value={this.state.query}/>
			</div>
		);
	}
}
