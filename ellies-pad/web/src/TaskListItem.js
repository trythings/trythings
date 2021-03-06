import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import TaskCard from './TaskCard.js';
import TaskTile from './TaskTile.js';

class TaskListItem extends React.Component {
	static propTypes = {
		task: React.PropTypes.shape({
			// ...TaskCard.propTypes.task
			// ...TaskTile.propTypes.task
			id: React.PropTypes.string.isRequired,
			title: React.PropTypes.string,
			description: React.PropTypes.string,
			isArchived: React.PropTypes.bool,
		}).isRequired,
	};

	static styles = {
		item: {
			...resetStyles,
			maxWidth: '100%',
			overflow: 'visible',
		},
		taskCard: {
			...resetStyles,
			flex: '1 0 auto',
		},
		taskTile: {
			...resetStyles,
			flex: '1 1 auto',
		},
	};

	state = {
		hasFocus: false,
	};

	onFocus = () => {
		this.setState({ hasFocus: true });
	};

	onBlur = (event) => {
		if (event.relatedTarget && !event.currentTarget.contains(event.relatedTarget)) {
			this.taskCard.requestClose();
		}
	};

	close = () => {
		this.setState({ hasFocus: false });
	};

	taskCardRef = (taskCard) => {
		this.taskCard = taskCard && taskCard.refs['component'];
	};

	render() {
		return (
			<div
				onFocus={this.onFocus}
				onBlur={this.onBlur}
				style={TaskListItem.styles.item}
				tabIndex={-1}
			>
				{this.state.hasFocus
					? (
							<TaskCard
								autoFocus
								style={TaskListItem.styles.taskCard}
								onClose={this.close}
								ref={this.taskCardRef}
								task={this.props.task}
							/>
						)
					: <TaskTile style={TaskListItem.styles.taskTile} task={this.props.task} />
				}
			</div>
		);
	}
}

export default Relay.createContainer(TaskListItem, {
	fragments: {
		task: () => Relay.QL`
			fragment on Task {
				id,
				${TaskTile.getFragment('task')},
				${TaskCard.getFragment('task')},
			}
		`,
	},
});
