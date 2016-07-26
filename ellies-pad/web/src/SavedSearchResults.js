import React from 'react';
import Relay from 'react-relay';

import TaskList from './TaskList.js';

class SavedSearchResults extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			space: React.PropTypes.shape({
				search: React.PropTypes.shape({
					tasks: React.PropTypes.shape({
						// ...TaskList.propTypes.tasks
					}).isRequired,
				}).isRequired,
			}).isRequired,
		}).isRequired,
	};

	render() {
		return <TaskList tasks={this.props.viewer.space.search.tasks} />;
	}
}

export default Relay.createContainer(SavedSearchResults, {
	initialVariables: {
		searchId: null,
	},

	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				space {
					search(id: $searchId) {
						tasks(first: 10) {
							${TaskList.getFragment('tasks')},
						}
					},
				},
			}
		`,
	},
});
