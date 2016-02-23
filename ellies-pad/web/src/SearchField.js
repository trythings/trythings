import _debounce from 'lodash/debounce';
import React from 'react';

import Icon from './Icon.js';

export default class SearchField extends React.Component {
	state = {
		query: '',
	};

	onChange = (event) => {
		const query = event.target.value;
		this.setState({ query });
		this.sendQuery(query);
	};

	sendQuery = _debounce((query) => {
		console.log(query);
	}, 200);

	render() {
		return (
			<div>
				<Icon name="search"/>
				<input onChange={this.onChange} value={this.state.query}/>
			</div>
		);
	}
}
