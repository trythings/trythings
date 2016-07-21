import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import SavedSearch from './SavedSearch.js';

export default class View extends React.Component {
	static propTypes = {
		view: React.PropTypes.shape({
			searches: React.PropTypes.arrayOf(React.PropTypes.shape({
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
		this.savedSearches = [];
	}

	savedSearchRef = (savedSearch) => {
		console.log(this);
		this.savedSearches.push(savedSearch.refs.component);
	};

	refetch = () => {
		this.savedSearches.forEach(s => s.refetch());
	}

	render() {
		return (
			<div style={View.styles.container}>
				{this.props.view.searches.map((search, i, searches) => {
					if (i === searches.length - 1) {
						return <SavedSearch search={search} ref={this.savedSearchRef} />;
					}
					return [
						<SavedSearch search={search} ref={this.savedSearchRef} />,
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
					${SavedSearch.getFragment('search')},
				},
			},
		`,
	},
});
