import React from 'react';
import Relay from 'react-relay';

import Card from './Card.js';
import resetStyles from './resetStyles.js';
import TaskTile from './TaskTile.js';
import theme from './theme.js';

class TaskList extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			// ...TaskTile.propTypes.task
			tasks: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
		}).isRequired,
	};

	static styles = {
		list: {
			...resetStyles,
			flexDirection: 'column',
			paddingTop: 4,
			paddingBottom: 4,
			overflow: 'visible',
		},
		divider: {
			...resetStyles,
			height: 1,
			backgroundColor: theme.text.dark.dividers.color,
		},
	};

	render() {
		return (
			<Card>
				<ol style={TaskList.styles.list}>
					{this.props.viewer.tasks.map((task, i, array) => {
						const tile = <TaskTile key={task.id} task={task}/>;
						if (i === array.length - 1) {
							return tile;
						}
						return [
							tile,
							<hr style={TaskList.styles.divider}/>,
						];
					})}
				</ol>
			</Card>
		);
	}
}

export default Relay.createContainer(TaskList, {
	initialVariables: {
		query: null,
	},

	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				tasks(query: $query) {
					id,
					${TaskTile.getFragment('task')},
				},
			}
		`,
	},
});
