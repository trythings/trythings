import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import TaskSearchResults from './TaskSearchResults.js';
import theme from './theme.js';

class TaskSearchRoute extends Relay.Route {
	static routeName = 'TaskSearchRoute';

	static paramDefinitions = {
		query: { required: true },
	};

	static queries = {
		viewer: () => Relay.QL`
			query {
				viewer,
			}
		`,
	};
}

class TaskSearch extends React.Component {
	static propTypes = {
		name: React.PropTypes.string,
		query: React.PropTypes.string,
		space: React.PropTypes.shape({
			name: React.PropTypes.string.isRequired,
		}).isRequired,
	};

	static styles = {
		container: {
			...resetStyles,

			alignItems: 'stretch',
			flexDirection: 'column',
			overflow: 'visible',
		},
		name: {
			...resetStyles,
			...theme.text.dark.secondary,

			fontSize: 14,
			paddingBottom: 8,
			paddingLeft: 16,
		},
		loading: {
			...resetStyles,
			...theme.text.dark.secondary,

			alignSelf: 'center',
			fontSize: 14,
		},
		noQuery: {
			...resetStyles,
			...theme.text.dark.secondary,

			alignSelf: 'center',
			fontSize: 14,
		},
		spaceName: {
			...resetStyles,
			...theme.text.dark.secondary,

			fontSize: 14,
			fontWeight: 500,
		},
	};

	renderLoading = () => (
			<span style={TaskSearch.styles.loading}>Loading...</span>
	);

	render() {
		if (!this.props.query) {
			return (
				<div style={TaskSearch.styles.container}>
					<span style={TaskSearch.styles.noQuery}>
						Enter #tags, @usernames, or keywords to find tasks in&nbsp;
						<span style={TaskSearch.styles.spaceName}>{this.props.space.name}</span>
					</span>
				</div>
			);
		}

		return (
			<div style={TaskSearch.styles.container}>
				<h1 style={TaskSearch.styles.name}>{this.props.name}</h1>
				<Relay.RootContainer
					Component={TaskSearchResults}
					route={new TaskSearchRoute({ query: this.props.query })}
					renderLoading={this.renderLoading}
				/>
			</div>
		);
	}
}

export default Relay.createContainer(TaskSearch, {
	fragments: {
		space: () => Relay.QL`
			fragment on Space {
				name,
			},
		`,
	},
});
