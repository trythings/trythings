import React from 'react';
import Relay from 'react-relay';

import TaskList from './TaskList.js';

class QuerySearchResults extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			space: React.PropTypes.shape({
				querySearch: React.PropTypes.shape({
					// ...TaskList.propTypes.search,
				}).isRequired,
			}).isRequired,
		}).isRequired,
	};

	render() {
		return <TaskList search={this.props.viewer.space.querySearch} />;
	}
}

export default Relay.createContainer(QuerySearchResults, {
	initialVariables: {
		query: null,
	},

	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				space {
					querySearch(query: $query) {
						${TaskList.getFragment('search')},
					},
				},
			}
		`,
	},
});
