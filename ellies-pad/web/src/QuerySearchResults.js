import React from 'react';
import Relay from 'react-relay';

import TaskList from './TaskList.js';

class QuerySearchResults extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			task: React.PropTypes.shape({
				querySearch: React.PropTypes.shape({
					// ...TaskList.propTypes.search,
				}).isRequired,
			}).isRequired,
		}).isRequired,
	};

	render() {
		return <TaskList search={this.props.viewer.task.querySearch} />;
	}
}

export default Relay.createContainer(QuerySearchResults, {
	initialVariables: {
		query: null,
	},

	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				task {
					querySearch(query: $query) {
						${TaskList.getFragment('search')},
					},
				},
			}
		`,
	},
});
