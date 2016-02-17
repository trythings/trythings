import React from 'react';
import Relay from 'react-relay';

import Card from './Card.js';
import resetStyles from './resetStyles.js';
import TaskListItem from './TaskListItem.js';
import theme from './theme.js';

class TaskList extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			// ...TaskListItem.propTypes.task
			tasks: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
		}).isRequired,
	};

	onFocus = (event) => {
		event.stopPropagation();
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
				<ol style={TaskList.styles.list} onFocus={this.onFocus} tabIndex={-1}>
					{this.props.viewer.tasks.map(task => (
						<TaskListItem
							key={task.id}
							task={task}
						/>
					))}
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
					${TaskListItem.getFragment('task')},
				},
			}
		`,
	},
});
