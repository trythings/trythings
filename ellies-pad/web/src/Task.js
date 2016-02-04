import React from 'react';
import Relay from 'react-relay';

class ArchiveTaskMutation extends Relay.Mutation {
	static fragments = {
		task: () => Relay.QL`
			fragment on Task {
				id,
			}
		`,
	};

	getMutation() {
		return Relay.QL`
			mutation {
				archiveTask,
			}
		`;
	}

	getFatQuery() {
		return Relay.QL`
			fragment on ArchiveTaskPayload {
				task {
					isArchived,
				},
			}
		`;
	}

	getConfigs() {
		return [{
			type: 'FIELDS_CHANGE',
			fieldIDs: {
				task: this.props.task.id,
			},
		}];
	}

	getVariables() {
		return {
			taskId: this.props.task.id,
		};
	}

	getOptimisticResponse() {
		return {
			task: {
				isArchived: true,
			},
		};
	}
}

class Task extends React.Component {
	static propTypes = {
		task: React.PropTypes.shape({
			title: React.PropTypes.string.isRequired,
			description: React.PropTypes.string,
			isArchived: React.PropTypes.boolean,
		}).isRequired,
	};

	state = {
		isHovering: false,
	};

	onMouseEnter = () => {
		this.setState({ isHovering: true });
	};

	onMouseLeave = () => {
		this.setState({ isHovering: false });
	};

	onArchiveClick = () => {
		Relay.Store.commitUpdate(
			new ArchiveTaskMutation({
				task: this.props.task,
			}),
		);
	};

	render() {
		return (
			<li onMouseEnter={this.onMouseEnter} onMouseLeave={this.onMouseLeave}>
				{this.props.task.isArchived ?
					<del>{this.props.task.title}</del> :
					<strong>{this.props.task.title}</strong>
				}
				{' ≫ '}
				<span>{this.props.task.description}</span>
				{' '}
				{this.state.isHovering ?
					<a href="#" onClick={this.onArchiveClick}>archive</a> :
					null
				}
			</li>
		);
	}
}

export default Relay.createContainer(Task, {
	fragments: {
		task: () => Relay.QL`
			fragment on Task {
				${ArchiveTaskMutation.getFragment('task')},
				title,
				description,
				isArchived,
			}
		`,
	},
});
