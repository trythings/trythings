import React from 'react';

import './MaterialIcons.css';

export default class Icon extends React.Component {
	static propTypes = {
		name: React.PropTypes.string.isRequired,
		color: React.PropTypes.string,
	};

	render() {
		return (
			<i
				className="material-icons"
				style={{ color: this.props.color }}
			>
				{this.props.name}
			</i>
		);
	}
}
