import React from 'react';
import Relay from 'react-relay';

import TaskList from './TaskList.js';

class SavedSearchResults extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			space: React.PropTypes.shape({
				savedSearch: React.PropTypes.shape({
					// ...TaskList.propTypes.search
				}).isRequired,
			}).isRequired,
		}).isRequired,
	};

	render() {
		return <TaskList search={this.props.viewer.space.savedSearch} />;
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
					savedSearch(id: $searchId) {
						${TaskList.getFragment('search')},
					},
				},
			}
		`,
	},
});
