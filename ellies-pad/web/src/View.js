import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import SavedSearch from './SavedSearch.js';

export default class View extends React.Component {
	static propTypes = {
		view: React.PropTypes.shape({
			searches: React.PropTypes.arrayOf(React.PropTypes.shape({
				id: React.PropTypes.string.isRequired,
				// ...SavedSearch.propTypes.search,
			})).isRequired,
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
			this.savedSearches.set(id, savedSearch.refs['component']);
		} else {
			this.savedSearches.delete(id);
		}
	};

	refetch = () => {
		this.savedSearches.forEach(s => s.refetch());
	}

	render() {
		return (
			<div style={View.styles.container}>
				{this.props.view.searches.map((search, i, searches) => {
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
		view: () => Relay.QL`
			fragment on View {
				searches {
					id,
					${SavedSearch.getFragment('search')},
				},
			},
		`,
	},
});
