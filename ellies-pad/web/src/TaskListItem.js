import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import TaskCard from './TaskCard.js';
import TaskTile from './TaskTile.js';

class TaskListItem extends React.Component {
	static propTypes = {
		// ...TaskCard.propTypes.task
		// ...TaskTile.propTypes.task
		task: React.PropTypes.shape({
			id: React.PropTypes.string.isRequired,
			title: React.PropTypes.string.isRequired,
			description: React.PropTypes.string,
			isArchived: React.PropTypes.boolean,
		}).isRequired,
	};

	state = {
		isFocused: false,
	};

	onFocus = () => {
		this.setState({ isFocused: true });
	};

	onBlur = (event) => {
		if (!event.currentTarget.contains(event.relatedTarget) && this.taskCard) {
			this.taskCard.requestClose();
		}
	};

	close = () => {
		this.setState({ isFocused: false });
	};

	taskCardRef = (taskCard) => {
		this.taskCard = taskCard && taskCard.refs.component;
	};

	static styles = {
		listItem: {
			...resetStyles,
			flexDirection: 'column',
			overflow: 'visible',
		},
	};

	render() {
		return (
			<li
				onFocus={this.onFocus}
				onBlur={this.onBlur}
				style={TaskListItem.styles.listItem}
				tabIndex={-1}
			>
				{this.state.isFocused
					? (
							<TaskCard
								autoFocus
								onClose={this.close}
								ref={this.taskCardRef}
								task={this.props.task}
							/>
						)
					: <TaskTile task={this.props.task}/>
				}
				{
					// i < array.length - 1
					// ? <hr style={TaskList.styles.divider}/>
					// : null
				}
			</li>
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
