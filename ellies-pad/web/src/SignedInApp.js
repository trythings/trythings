import _debounce from 'lodash/debounce';
import gapi from 'gapi';
import React from 'react';
import Relay from 'react-relay';

import AddTask from './AddTask.js';
import AppBar from './AppBar.js';
import NavigationDrawer from './NavigationDrawer.js';
import QuerySearch from './QuerySearch.js';
import resetStyles from './resetStyles.js';
import SearchField from './SearchField.js';
import theme from './theme.js';
import View from './View.js';

class SignedInApp extends React.Component {
	static contextTypes = {
		router: React.PropTypes.shape({
			push: React.PropTypes.func.isRequired,
		}).isRequired,
	};

	// Routing.
	static onEnter = (nextState, replace) => {
		if (!gapi.auth2.getAuthInstance().isSignedIn.get()) {
			replace('/signin');
		}
	};

	static propTypes = {
		children: React.PropTypes.node,

		// Routing.
		location: React.PropTypes.shape({
			pathname: React.PropTypes.string.isRequired,
		}).isRequired,
		params: React.PropTypes.shape({
			query: React.PropTypes.string,
		}).isRequired,

		viewer: React.PropTypes.shape({
			space: React.PropTypes.shape({
				// ...AddTask.propTypes.space,
				// ...QuerySearch.propTypes.space,
				view: React.PropTypes.shape({
					// ...View.propTypes.view,
				}).isRequired,
			}).isRequired,
			spaces: React.PropTypes.arrayOf(React.PropTypes.shape({
				// ...NavigationDrawer.propTypes.space,
			})).isRequired,
		}).isRequired,
	};

	static styles = {
		rootContainer: {
			...resetStyles,
			alignItems: 'stretch',
			flex: '1 1 auto',
		},
		appBarContainer: {
			...resetStyles,
			alignItems: 'stretch',
			flexDirection: 'column',
			flex: '1 1 auto',
		},
		appBar: {
			...resetStyles,
			...theme.text.light.primary,
			backgroundColor: theme.colors.primary.default,
		},
		searchField: {
			...resetStyles,
			...theme.text.light.primary,
			backgroundColor: theme.colors.primary.light,
			flex: '1 0 auto',
		},
		overlayContainer: {
			...resetStyles,
			alignItems: 'stretch',
			flex: '1 1 auto',
		},
		addTaskContainer: {
			...resetStyles,

			alignItems: 'flex-end',
			justifyContent: 'flex-end',

			position: 'absolute',
			bottom: 0,
			left: 0,
			right: 0,
			top: 0,

			pointerEvents: 'none',
			flexDirection: 'column',
		},
		addTaskSpacer: {
			...resetStyles,
			paddingTop: 56,
		},
		addTask: {
			...resetStyles,
			flex: '1 1 auto',
			pointerEvents: 'auto',
		},
		contentContainer: {
			...resetStyles,
			alignItems: 'stretch',
			flex: '1 1 auto',
			flexDirection: 'column',
			overflow: 'scroll',
		},
		content: {
			...resetStyles,
			alignItems: 'stretch',
			flex: '1 0 auto',
			flexDirection: 'column',
			paddingBottom: 24,
			paddingLeft: 24,
			paddingRight: 24,
			paddingTop: 24,
			overflow: 'visible',
		},
		contentSpacer: {
			...resetStyles,
			padding: 12,
		},
	};

	constructor(props, ...args) {
		super(props, ...args);

		let searchQuery = this.props.params.query;
		if (this.props.location.pathname === '/search/') {
			searchQuery = '';
		}
		this.state = { searchQuery };
	}

	onSearchFocus = () => {
		if (this.state.searchQuery === undefined) {
			this.setState({ searchQuery: '' });
			this.context.router.push('/search/');
		}
	};

	onSearchBlur = () => {
		if (!this.state.searchQuery) {
			this.setState({ searchQuery: undefined });
			this.context.router.push('/');
		}
	};

	onSearchQueryChange = (query) => {
		this.updateSearchResults(query);
		this.updateSearchPath(query);
	};

	refetch = () => {
		if (this.view) {
			this.view.refetch();
		}
	};

	updateSearchPath = _debounce((query) => {
		// This undefined check exists because the user may have
		// emptied and blurred the search field in quick succession,
		// which could cause us to redirect from / to /search/${query}.
		if (this.state.searchQuery !== undefined) {
			this.context.router.push(`/search/${encodeURIComponent(query)}`);
		}
	}, 500);

	updateSearchResults = _debounce((query) => {
		this.setState({ searchQuery: query });
	}, 200);

	viewRef = (view) => {
		if (view) {
			this.view = view.refs.component;
		} else {
			this.view = null;
		}
	};

	renderContent() {
		if (this.state.searchQuery !== undefined) {
			return (
				<QuerySearch
					name="Search results"
					query={this.state.searchQuery}
					space={this.props.viewer.space}
				/>
			);
		}

		return <View ref={this.viewRef} view={this.props.viewer.space.view} />;
	}

	render() {
		let appBarStyle = SignedInApp.styles.appBar;
		if (this.state.searchQuery !== undefined) {
			appBarStyle = {
				...appBarStyle,
				...theme.text.dark.primary,
				backgroundColor: theme.colors.card,
			};
		}

		let searchFieldStyle = SignedInApp.styles.searchField;
		if (this.state.searchQuery !== undefined) {
			searchFieldStyle = {
				...searchFieldStyle,
				...theme.text.dark.primary,
				backgroundColor: theme.colors.card,
			};
		}

		return (
			<div style={SignedInApp.styles.rootContainer}>
				<NavigationDrawer spaces={this.props.viewer.spaces} />

				<div style={SignedInApp.styles.appBarContainer}>
					<AppBar style={appBarStyle}>
						<SearchField
							autoFocus={this.state.searchQuery === ''}
							initialQuery={this.state.searchQuery}
							onQueryChange={this.onSearchQueryChange}
							onFocus={this.onSearchFocus}
							onBlur={this.onSearchBlur}
							style={searchFieldStyle}
						/>
					</AppBar>

					<div style={SignedInApp.styles.overlayContainer}>
						<div style={SignedInApp.styles.contentContainer}>
							<div style={SignedInApp.styles.content}>
								{this.renderContent()}
							</div>
						</div>

						<div style={SignedInApp.styles.addTaskContainer}>
							<div style={SignedInApp.styles.addTaskSpacer} />
							<AddTask
								refetch={this.refetch}
								space={this.props.viewer.space}
								style={SignedInApp.styles.addTask}
							/>
						</div>
					</div>
				</div>
			</div>
		);
	}
}

const SignedInAppContainer = Relay.createContainer(SignedInApp, {
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				space {
					name,
					${AddTask.getFragment('space')},
					${QuerySearch.getFragment('space')},
					view {
						${View.getFragment('view')},
					},
				},
				spaces {
					${NavigationDrawer.getFragment('spaces')},
				},
			}
		`,
	},
});

SignedInAppContainer.onEnter = SignedInApp.onEnter;

export default SignedInAppContainer;
