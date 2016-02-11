import React from 'react';
import Relay from 'react-relay';

import AddTaskCard from './AddTaskCard.js';
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

	state = {
		focusedTaskId: null,
	};

	onFocus = (event) => {
		this.setState({ focusedTaskId: event.currentTarget.dataset.taskid });
		event.stopPropagation();
	};

	onBlur = (event) => {
		if (!event.currentTarget.contains(event.relatedTarget)) {
			this.setState({ focusedTaskId: null });
		}
	};

	static styles = {
		list: {
			...resetStyles,
			flexDirection: 'column',
			paddingTop: 4,
			paddingBottom: 4,
			overflow: 'visible',
		},
		listItem: {
			...resetStyles,
			// backgroundColor: 'red',
			flexDirection: 'column',
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
				<ol style={TaskList.styles.list} onBlur={this.onBlur}>
					{this.props.viewer.tasks.map((task, i, array) => (
						<li
							key={task.id}
							data-taskid={task.id}
							onFocus={this.onFocus}
							style={TaskList.styles.listItem}
							tabIndex={-1}
						>
							{task.id === this.state.focusedTaskId
								? <AddTaskCard autoFocus viewer={this.props.viewer}/>
								: <TaskTile task={task}/>
							}
							{i < array.length - 1
								? <hr style={TaskList.styles.divider}/>
								: null
							}
						</li>
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
				${AddTaskCard.getFragment('viewer')},
				tasks(query: $query) {
					id,
					${TaskTile.getFragment('task')},
				},
			}
		`,
	},
});
