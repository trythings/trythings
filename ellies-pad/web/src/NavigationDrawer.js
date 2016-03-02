import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import theme from './theme.js';

class NavigationDrawer extends React.Component {
	static propTypes = {
		spaces: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
	};

	static styles = {
		nav: {
			...resetStyles,
			backgroundColor: theme.colors.card,
			borderLeft: `1px solid ${theme.text.dark.dividers.color}`,
			borderRight: `1px solid ${theme.text.dark.dividers.color}`,
			paddingLeft: 16,
			paddingRight: 16,
			width: 240,
		},
		list: {
			...resetStyles,
		},
		text: {
			...resetStyles,
			...theme.text.dark.primary,
			fontSize: 14,
			fontWeight: 500,
		},
	};

	render() {
		return (
			<nav style={NavigationDrawer.styles.nav}>
				<ul style={NavigationDrawer.styles.list}>
					{this.props.spaces.map(space => (
						<li key={space.id}>
							<span style={NavigationDrawer.styles.text}>{space.name}</span>
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
			}
		`,
	},
});
