import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import SavedSearch from './SavedSearch.js';

export default class View extends React.Component {
	static propTypes = {
		task: React.PropTypes.shape({
			searches: React.PropTypes.shape({
				edges: React.PropTypes.arrayOf(React.PropTypes.shape({
					node: React.PropTypes.shape({
						id: React.PropTypes.string.isRequired,
						// ...SavedSearch.propTypes.search,
					}).isRequired,
				})).isRequired,
			}).isRequired,
		}).isRequired,
	};

	static styles = {
		container: {
			...resetStyles,
			alignItems: 'stretch',
			flexDirection: 'column',
			overflow: 'visible',
		},
		spacer: {
			...resetStyles,
			padding: 12,
		},
	};

	constructor(...args) {
		super(...args);
		this.savedSearches = new Map();
	}

	savedSearchRef = (id) => (savedSearch) => {
		if (savedSearch) {
			this.savedSearches.set(id, savedSearch.refs.component);
		} else {
			this.savedSearches.delete(id);
		}
	};

	refetch = () => {
		this.savedSearches.forEach(s => s.refetch());
	}

	render() {
		const viewSearches = this.props.task.searches.edges.map((edge) => edge.node);
		return (
			<div style={View.styles.container}>
				{viewSearches.map((search, i, searches) => {
					if (i === searches.length - 1) {
						return (
							<SavedSearch
								key={search.id}
								search={search}
								ref={this.savedSearchRef(search.id)}
							/>
						);
					}
					return [
						<SavedSearch key={search.id} search={search} ref={this.savedSearchRef(search.id)} />,
						<div style={View.styles.spacer} />,
					];
				})}
			</div>
		);
	}
}

export default Relay.createContainer(View, {
	fragments: {
		task: () => Relay.QL`
			fragment on Task {
				searches(first: 1000) {
					edges {
						node {
							id,
							${SavedSearch.getFragment('search')},
						},
					},
				},
			},
		`,
	},
});
