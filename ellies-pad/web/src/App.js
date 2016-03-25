import React from 'react';
import Relay from 'react-relay';

import ActionButton from './ActionButton.js';
import AddTaskCard from './AddTaskCard.js';
import AppBar from './AppBar.js';
import NavigationDrawer from './NavigationDrawer.js';
import resetStyles from './resetStyles.js';
import theme from './theme.js';

class App extends React.Component {
	static propTypes = {
		children: React.PropTypes.node,
		location: React.PropTypes.shape({
			pathname: React.PropTypes.string.isRequired,
		}).isRequired,
		params: React.PropTypes.shape({
			query: React.PropTypes.string,
		}).isRequired,
		viewer: React.PropTypes.shape({
			space: React.PropTypes.shape({
				// ...AddTaskCard.propTypes.space
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

	state = {
		isAddTaskFormVisible: false,
	};

	onCancelClick = () => {
		this.setState({ isAddTaskFormVisible: false });
	};

	onPlusClick = () => {
		this.setState({ isAddTaskFormVisible: true });
	};

	render() {
		let searchQuery = this.props.params.query;
		if (this.props.location.pathname === '/search/') {
			searchQuery = '';
		}

		return (
			<div style={App.styles.app} tabIndex={-1}>
				<AppBar searchQuery={searchQuery} />

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

							{this.props.children}
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
				},
				spaces {
					${NavigationDrawer.getFragment('spaces')},
				},
			}
		`,
	},
});
