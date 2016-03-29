import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import QuerySearchResults from './QuerySearchResults.js';
import theme from './theme.js';

class QuerySearchRoute extends Relay.Route {
	static routeName = 'QuerySearchRoute';

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

class QuerySearch extends React.Component {
	static propTypes = {
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
			display: 'inline',
			fontSize: 14,
		},
		spaceName: {
			...resetStyles,
			...theme.text.dark.secondary,

			display: 'inline',
			fontSize: 14,
			fontWeight: 500,
		},
	};

	renderLoading = () => (
			<span style={QuerySearch.styles.loading}>Loading...</span>
	);

	render() {
		return (
			<div style={QuerySearch.styles.container}>
				{this.props.query ?
					(
						<Relay.RootContainer
							Component={QuerySearchResults}
							route={new QuerySearchRoute({ query: this.props.query })}
							renderLoading={this.renderLoading}
						/>
					) :
					(
						<span style={QuerySearch.styles.noQuery}>
							Enter #tags, @usernames, or keywords to find tasks in
							<span style={QuerySearch.styles.spaceName}>&nbsp;{this.props.space.name}</span>
						</span>
					)
				}
			</div>
		);
	}
}

export default Relay.createContainer(QuerySearch, {
	fragments: {
		space: () => Relay.QL`
			fragment on Space {
				name,
			},
		`,
	},
});
