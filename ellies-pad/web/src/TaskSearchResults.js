import React from 'react';
import Relay from 'react-relay';

import TaskList from './TaskList.js';

class TaskSearchResults extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			space: React.PropTypes.shape({
				// ...TaskList.propTypes.tasks
				tasks: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
			}).isRequired,
		}).isRequired,
	};

	render() {
		return <TaskList tasks={this.props.viewer.space.tasks}/>;
	}
}

export default Relay.createContainer(TaskSearchResults, {
	initialVariables: {
		query: null,
	},

	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				space {
					tasks(query: $query) {
						${TaskList.getFragment('tasks')},
					},
				},
			}
		`,
	},
});
