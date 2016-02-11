import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import TaskList from './TaskList.js';
import theme from './theme.js';

class TaskSearchRoute extends Relay.Route {
	static routeName = 'TaskSearchRoute';

	static paramDefinitions = {
		query: { required: true },
	};

	static queries = {
		viewer: () => Relay.QL`
			query {
				viewer,
			}
		`,
	};
}

export default class TaskSearch extends React.Component {
	static propTypes = {
		name: React.PropTypes.string,
		query: React.PropTypes.string,
	};

	static styles = {
		container: {
			...resetStyles,

			flexDirection: 'column',
			overflow: 'visible',
		},
		name: {
			...resetStyles,
			...theme.text.dark.secondary,

			fontSize: 14,
			paddingBottom: 8,
			paddingLeft: 16,
		},
	};

	render() {
		return (
			<div style={TaskSearch.styles.container}>
				<h1 style={TaskSearch.styles.name}>{this.props.name}</h1>
				<Relay.RootContainer
					Component={TaskList}
					route={new TaskSearchRoute({ query: this.props.query })}
				/>
			</div>
		);
	}
}
