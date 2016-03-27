import _debounce from 'lodash/debounce';
import React from 'react';
import Relay from 'react-relay';

import ActionButton from './ActionButton.js';
import AddTaskCard from './AddTaskCard.js';
import AppBar from './AppBar.js';
import DefaultView from './DefaultView.js';
import NavigationDrawer from './NavigationDrawer.js';
import resetStyles from './resetStyles.js';
import SearchField from './SearchField.js';
import TaskSearch from './TaskSearch.js';
import theme from './theme.js';

class App extends React.Component {
	static contextTypes = {
		router: React.PropTypes.shape({
			push: React.PropTypes.func.isRequired,
		}).isRequired,
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
				// ...AddTaskCard.propTypes.space,
				// ...DefaultView.propTypes.space,
				// ...TaskSearch.propTypes.space,
			}),
			spaces: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
		}).isRequired,
	};

	static styles = {
		app: {
			...resetStyles,
			alignItems: 'stretch',
			backgroundColor: theme.colors.canvas,
			flexDirection: 'column',
			height: '100%',
			width: '100%',
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
		navigationContainer: {
			...resetStyles,
			alignItems: 'stretch',
			flex: '1 1 auto',
		},
		contentContainer: {
			...resetStyles,
			alignItems: 'stretch',
			flex: '1 1 auto',
			flexDirection: 'column',
			overflow: 'scroll',
		},
		addTaskButton: {
			...resetStyles,
			overflow: 'visible',
			position: 'absolute',

			right: 24,
			top: 24,
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
		this.state = {
			isAddTaskFormVisible: false,
			searchQuery,
		};
	}

	onCancelClick = () => {
		this.setState({ isAddTaskFormVisible: false });
	};

	onPlusClick = () => {
		this.setState({ isAddTaskFormVisible: true });
	};

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

	renderContent() {
		if (this.state.searchQuery !== undefined) {
			return (
				<TaskSearch
					name="Search results"
					query={this.state.searchQuery}
					space={this.props.viewer.space}
				/>
			);
		}

		return <DefaultView space={this.props.viewer.space} />;
	}

	render() {
		let appBarStyle = App.styles.appBar;
		if (this.state.searchQuery !== undefined) {
			appBarStyle = {
				...appBarStyle,
				...theme.text.dark.primary,
				backgroundColor: theme.colors.card,
			};
		}

		let searchFieldStyle = App.styles.searchField;
		if (this.state.searchQuery !== undefined) {
			searchFieldStyle = {
				...searchFieldStyle,
				...theme.text.dark.primary,
				backgroundColor: theme.colors.card,
			};
		}

		return (
			<div style={App.styles.app} tabIndex={-1}>
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

				<div style={App.styles.navigationContainer}>
					{!this.state.isAddTaskFormVisible ?
						(
							<div style={App.styles.addTaskButton}>
								<ActionButton onClick={this.onPlusClick} />
							</div>
						) :
						null
					}

					<NavigationDrawer spaces={this.props.viewer.spaces} />

					<div style={App.styles.contentContainer}>
						<div style={App.styles.content}>
							{this.state.isAddTaskFormVisible ?
								(
									<AddTaskCard
										autoFocus
										space={this.props.viewer.space}
										onCancelClick={this.onCancelClick}
									/>
								) :
								null
							}

							{this.state.isAddTaskFormVisible ?
								<div style={App.styles.contentSpacer} /> :
								null
							}

							{this.renderContent()}
						</div>
					</div>
				</div>
			</div>
		);
	}
}

export default Relay.createContainer(App, {
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				space {
					${AddTaskCard.getFragment('space')},
					${DefaultView.getFragment('space')},
					${TaskSearch.getFragment('space')},
				},
				spaces {
					${NavigationDrawer.getFragment('spaces')},
				},
			}
		`,
	},
});
