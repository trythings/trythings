import React from 'react';
import Relay from 'react-relay';

import Card from './Card.js';
import resetStyles from './resetStyles.js';
import TaskCard from './TaskCard.js';
import TaskTile from './TaskTile.js';
import theme from './theme.js';

class TaskList extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			// ...TaskTile.propTypes.task
			// ...TaskCard.propTypes.task
			tasks: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
		}).isRequired,
	};

	state = {
		focusedTaskId: null,
	};

	onTaskCardClose = () => {
		this.setState({ focusedTaskId: null });
	};

	onFocus = (event) => {
		this.setState({ focusedTaskId: event.currentTarget.dataset.taskid });
		event.stopPropagation();
	};

	onBlur = (event) => {
		if (!event.currentTarget.contains(event.relatedTarget) && this.taskCard) {
			this.taskCard.requestClose();
		}
	};

	taskCardRef = (taskCard) => {
		this.taskCard = taskCard && taskCard.refs.component;
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
				<ol style={TaskList.styles.list}>
					{this.props.viewer.tasks.map((task, i, array) => (
						<li
							key={task.id}
							data-taskid={task.id}
							onFocus={this.onFocus}
							onBlur={this.onBlur}
							style={TaskList.styles.listItem}
							tabIndex={-1}
						>
							{task.id === this.state.focusedTaskId
								? (
										<TaskCard
											autoFocus
											onClose={this.onTaskCardClose}
											ref={this.taskCardRef}
											task={task}
										/>
									)
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
				tasks(query: $query) {
					id,
					${TaskTile.getFragment('task')},
					${TaskCard.getFragment('task')},
				},
			}
		`,
	},
});
