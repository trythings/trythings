import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import SavedSearchResults from './SavedSearchResults.js';
import theme from './theme.js';

class SavedSearchRoute extends Relay.Route {
	static routeName = 'SavedSearchRoute';

	static paramDefinitions = {
		searchId: { required: true },
	};

	static queries = {
		viewer: () => Relay.QL`
			query {
				viewer,
			}
		`,
	};
}

class SavedSearch extends React.Component {
	static propTypes = {
		search: React.PropTypes.shape({
			id: React.PropTypes.string.isRequired,
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
	};

	renderLoading = () => (
			<span style={SavedSearch.styles.loading}>Loading...</span>
	);

	render() {
		return (
			<div style={SavedSearch.styles.container}>
				<h1 style={SavedSearch.styles.name}>{this.props.search.name}</h1>
				<Relay.RootContainer
					Component={SavedSearchResults}
					route={new SavedSearchRoute({ searchId: this.props.search.id })}
					renderLoading={this.renderLoading}
				/>
			</div>
		);
	}
}

export default Relay.createContainer(SavedSearch, {
	fragments: {
		search: () => Relay.QL`
			fragment on Search {
				id,
				name,
			},
		`,
	},
});

