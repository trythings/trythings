// import React from 'react';
// import Relay from 'react-relay';

// import FlatButton from './FlatButton.js';
// import resetStyles from './resetStyles.js';
// import theme from './theme.js';

// // TODO#Rewrite: The NavigationDrawer is no longer in use. Repurpose it or remove it.
// class NavigationDrawer extends React.Component {
// 	static propTypes = {
// 		onSignOutClick: React.PropTypes.func.isRequired,
// 		onViewNameClick: React.PropTypes.func.isRequired,

// 		spaces: React.PropTypes.arrayOf(React.PropTypes.shape({
// 			id: React.PropTypes.string.isRequired,
// 			name: React.PropTypes.string.isRequired,
// 			views: React.PropTypes.arrayOf(React.PropTypes.shape({
// 				id: React.PropTypes.string.isRequired,
// 				name: React.PropTypes.string.isRequired,
// 			})).isRequired,
// 		})).isRequired,

// 		selectedViewId: React.PropTypes.string,
// 	};

// 	static styles = {
// 		nav: {
// 			...resetStyles,
// 			backgroundColor: theme.colors.card,
// 			borderLeft: `1px solid ${theme.colors.dividers.dark}`,
// 			borderRight: `1px solid ${theme.colors.dividers.dark}`,
// 			flexDirection: 'column',
// 			paddingBottom: 16,
// 			paddingLeft: 16,
// 			paddingRight: 16,
// 			paddingTop: 16,
// 			width: 240,
// 		},
// 		profile: {
// 			...resetStyles,
// 			height: 56, // 56 to align with the app bar.
// 		},
// 		signOutButton: {
// 			...resetStyles,
// 			color: theme.colors.primary.light,
// 		},
// 		spaces: {
// 			...resetStyles,
// 			flexDirection: 'column',
// 		},
// 		space: {
// 			...resetStyles,
// 			flexDirection: 'column',
// 		},
// 		spaceName: {
// 			...resetStyles,
// 			...theme.text.dark.secondary,
// 			fontSize: 14,
// 			fontWeight: 500,
// 			paddingBottom: 16,
// 		},
// 		views: {
// 			...resetStyles,
// 			flexDirection: 'column',
// 			paddingLeft: 16,
// 		},
// 		view: {
// 			...resetStyles,
// 			flexDirection: 'column',
// 		},
// 		viewName: {
// 			...resetStyles,
// 			...theme.text.dark.secondary,
// 			cursor: 'pointer',
// 			fontSize: 14,
// 			paddingBottom: 16,
// 		},
// 		selectedViewName: {
// 			...resetStyles,
// 			...theme.text.dark.primary,
// 			cursor: 'pointer',
// 			fontSize: 14,
// 			fontWeight: 500,
// 			paddingBottom: 16,
// 		},
// 		searches: {
// 			...resetStyles,
// 			flexDirection: 'column',
// 			paddingLeft: 16 * 2,
// 		},
// 		search: {
// 			...resetStyles,
// 			flexDirection: 'column',
// 		},
// 		searchName: {
// 			...resetStyles,
// 			...theme.text.dark.secondary,
// 			fontSize: 14,
// 			paddingBottom: 16,
// 		},
// 	};

// 	onViewNameClick = (event) => {
// 		this.props.onViewNameClick(event.currentTarget.dataset.viewId);
// 	}

// 	render() {
// 		return (
// 			<nav style={NavigationDrawer.styles.nav}>
// 				<div style={NavigationDrawer.styles.profile}>
// 					<FlatButton
// 						label="Sign out"
// 						onClick={this.props.onSignOutClick}
// 						style={NavigationDrawer.styles.signOutButton}
// 					/>
// 				</div>

// 				<ul style={NavigationDrawer.styles.spaces}>
// 					{this.props.spaces.map(space => (
// 						<li style={NavigationDrawer.styles.space} key={space.id}>
// 							<span style={NavigationDrawer.styles.spaceName}>{space.name}</span>
// 							<ul style={NavigationDrawer.styles.views}>
// 								{space.views.map((view, index) => {
// 									let nameStyle = NavigationDrawer.styles.viewName;
// 									if ((!this.props.selectedViewId && !index) ||
// 										this.props.selectedViewId === view.id) {
// 										nameStyle = NavigationDrawer.styles.selectedViewName;
// 									}
// 									return (
// 										<li style={NavigationDrawer.styles.view} key={view.id}>
// 											<a onClick={this.onViewNameClick} data-view-id={view.id}>
// 												<span style={nameStyle}>{view.name}</span>
// 											</a>
// 											<ul style={NavigationDrawer.styles.searches}>
// 												{view.searches.map(search => (
// 													<li style={NavigationDrawer.styles.search} key={search.id}>
// 														<span style={NavigationDrawer.styles.searchName}>
// 															{search.name}
// 														</span>
// 													</li>
// 												))}
// 											</ul>
// 										</li>
// 									);
// 								})}
// 							</ul>
// 						</li>
// 					))}
// 				</ul>
// 			</nav>
// 		);
// 	}
// }

// // export default Relay.createContainer(NavigationDrawer, {
// // 	fragments: {
// // 		spaces: () => Relay.QL`
// // 			fragment on Space @relay(plural: true) {
// // 				id,
// // 				name,
// // 				views {
// // 					id,
// // 					name,
// // 					searches {
// // 						id,
// // 						name,
// // 					},
// // 				},
// // 			},
// // 		`,
// // 	},
// // });
