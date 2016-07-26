import React from 'react';
import Relay from 'react-relay';

import TaskList from './TaskList.js';

class QuerySearchResults extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			space: React.PropTypes.shape({
				tasks: React.PropTypes.shape({
					// ...TaskList.propTypes.tasks,
				}).isRequired,
			}).isRequired,
		}).isRequired,
	};

	render() {
		return <TaskList tasks={this.props.viewer.space.tasks} />;
	}
}

export default Relay.createContainer(QuerySearchResults, {
	initialVariables: {
		query: null,
	},

	// TODO#xcxc: It's really unpleasant that this has to specify the pagination argument.
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				space {
					tasks(query: $query, first: 10) {
						${TaskList.getFragment('tasks')},
					},
				},
			}
		`,
	},
});
