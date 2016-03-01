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
		},
	};

	render() {
		return (
			<nav style={NavigationDrawer.styles.nav}>
				<ul>
					{this.props.spaces.map(space => (
						<li key={space.id}>
							{space.name}
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
