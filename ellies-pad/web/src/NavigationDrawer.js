import React from 'react';
import Relay from 'react-relay';

import FlatButton from './FlatButton.js';
import resetStyles from './resetStyles.js';
import theme from './theme.js';

class NavigationDrawer extends React.Component {
	static propTypes = {
		onSignOutClick: React.PropTypes.func.isRequired,

		spaces: React.PropTypes.arrayOf(React.PropTypes.shape({
			id: React.PropTypes.string.isRequired,
			name: React.PropTypes.string.isRequired,
			view: React.PropTypes.shape({
				id: React.PropTypes.string.isRequired,
				searches: React.PropTypes.arrayOf(React.PropTypes.shape({
					id: React.PropTypes.string.isRequired,
					name: React.PropTypes.string.isRequired,
				})).isRequired,
			}).isRequired,
			views: React.PropTypes.arrayOf(React.PropTypes.shape({
				id: React.PropTypes.string.isRequired,
				name: React.PropTypes.string.isRequired,
			})).isRequired,
		})).isRequired,
	};

	static styles = {
		nav: {
			...resetStyles,
			backgroundColor: theme.colors.card,
			borderLeft: `1px solid ${theme.colors.dividers.dark}`,
			borderRight: `1px solid ${theme.colors.dividers.dark}`,
			flexDirection: 'column',
			paddingBottom: 16,
			paddingLeft: 16,
			paddingRight: 16,
			paddingTop: 16,
			width: 240,
		},
		profile: {
			...resetStyles,
			height: 56, // 56 to align with the app bar.
		},
		signOutButton: {
			...resetStyles,
			color: theme.colors.primary.light,
		},
		spaces: {
			...resetStyles,
			flexDirection: 'column',
		},
		space: {
			...resetStyles,
			flexDirection: 'column',
		},
		spaceName: {
			...resetStyles,
			...theme.text.dark.secondary,
			fontSize: 14,
			fontWeight: 500,
			paddingBottom: 16,
		},
		views: {
			...resetStyles,
			flexDirection: 'column',
			paddingLeft: 16,
		},
		view: {
			...resetStyles,
			flexDirection: 'column',
		},
		viewName: {
			...resetStyles,
			...theme.text.dark.secondary,
			fontSize: 14,
			paddingBottom: 16,
		},
		selectedViewName: {
			...resetStyles,
			...theme.text.dark.primary,
			fontSize: 14,
			fontWeight: 500,
			paddingBottom: 16,
		},
		searches: {
			...resetStyles,
			flexDirection: 'column',
			paddingLeft: 16 * 2,
		},
		search: {
			...resetStyles,
			flexDirection: 'column',
		},
		searchName: {
			...resetStyles,
			...theme.text.dark.secondary,
			fontSize: 14,
			paddingBottom: 16,
		},
	};

	render() {
		return (
			<nav style={NavigationDrawer.styles.nav}>
				<div style={NavigationDrawer.styles.profile}>
					<FlatButton
						label="Sign out"
						onClick={this.props.onSignOutClick}
						style={NavigationDrawer.styles.signOutButton}
					/>
				</div>

				<ul style={NavigationDrawer.styles.spaces}>
					{this.props.spaces.map(space => (
						<li style={NavigationDrawer.styles.space} key={space.id}>
							<span style={NavigationDrawer.styles.spaceName}>{space.name}</span>
							<ul style={NavigationDrawer.styles.views}>
								{space.views.map(view => {
									const nameStyle = space.view.id === view.id ?
											NavigationDrawer.styles.selectedViewName :
											NavigationDrawer.styles.viewName;
									return (
										<li style={NavigationDrawer.styles.view} key={view.id}>
											<span style={nameStyle}>{view.name}</span>
											<ul style={NavigationDrawer.styles.searches}>
												{view.searches.map(search => (
													<li style={NavigationDrawer.styles.search} key={search.id}>
														<span style={NavigationDrawer.styles.searchName}>
															{search.name}
														</span>
													</li>
												))}
											</ul>
										</li>
									);
								})}
							</ul>
						</li>
					))}
				</ul>
			</nav>
		);
	}
}

export default Relay.createContainer(NavigationDrawer, {
	fragments: {
		spaces: () => Relay.QL`
			fragment on Space @relay(plural: true) {
				id,
				name,
				view {
					id,
				},
				views {
					id,
					name,
					searches {
						id,
						name,
					},
				},
			},
		`,
	},
});
