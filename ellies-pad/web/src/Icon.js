import React from 'react';

import './MaterialIcons.css';

export default class Icon extends React.Component {
	static propTypes = {
		name: React.PropTypes.string.isRequired,
	};

	render() {
		return <i {...this.props} className="material-icons">{this.props.name}</i>;
	}
}
