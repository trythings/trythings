import React from 'react';
import Relay from 'react-relay';

class Task extends React.Component {
	static propTypes = {
		task: React.PropTypes.shape({
			title: React.PropTypes.string.isRequired,
			description: React.PropTypes.string,
		}).isRequired,
	};

	render() {
		return (
			<li>
				<strong>{this.props.task.title}</strong> ≫ <span>{this.props.task.description}</span>
			</li>
		);
	}
}

export default Relay.createContainer(Task, {
	fragments: {
		task: () => Relay.QL`
			fragment on Task {
				title,
				description,
			}
		`,
	},
});
