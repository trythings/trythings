import React from 'react';
import Relay from 'react-relay';

import TaskList from './TaskList.js';

class SavedSearchResults extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			space: React.PropTypes.shape({
				savedSearch: React.PropTypes.shape({
					tasks: React.PropTypes.shape({
						// ...TaskList.propTypes.tasks
					}).isRequired,
				}).isRequired,
			}).isRequired,
		}).isRequired,
	};

	render() {
		return <TaskList tasks={this.props.viewer.space.savedSearch.tasks} />;
	}
}

export default Relay.createContainer(SavedSearchResults, {
	initialVariables: {
		searchId: null,
	},

	// TODO #xcxc: Pull tasks out into TaskList.
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				space {
					savedSearch(id: $searchId) {
						tasks(first: 10) {
							${TaskList.getFragment('tasks')},
						}
					},
				},
			}
		`,
	},
});
