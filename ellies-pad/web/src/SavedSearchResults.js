import React from 'react';
import Relay from 'react-relay';

import TaskList from './TaskList.js';

class SavedSearchResults extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			task: React.PropTypes.shape({
				savedSearch: React.PropTypes.shape({
					// ...TaskList.propTypes.search
				}).isRequired,
			}).isRequired,
		}).isRequired,
	};

	render() {
		return <TaskList search={this.props.viewer.task.savedSearch} />;
	}
}

export default Relay.createContainer(SavedSearchResults, {
	initialVariables: {
		searchId: null,
	},

	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				task {
					savedSearch(id: $searchId) {
						${TaskList.getFragment('search')},
					},
				},
			}
		`,
	},
});
