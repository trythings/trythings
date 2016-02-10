import React from 'react';
import Relay from 'react-relay';

import Icon from './Icon.js';
import theme from './theme.js';

class ArchiveTaskMutation extends Relay.Mutation {
	static fragments = {
		task: () => Relay.QL`
			fragment on Task {
				id,
				isArchived,
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
			newIsArchived: !this.props.task.isArchived,
		};
	}

	getOptimisticResponse() {
		return {
			task: {
				isArchived: !this.props.task.isArchived,
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

	static styles = {
		tile: {
			boxSizing: 'border-box',
			display: 'flex',

			// paddingTop: 20,
			// paddingBottom: 20,
			paddingLeft: 16,
			paddingRight: 16,
			height: 60,
			alignItems: 'center',
			justifyContent: 'space-between',
		},
		title: {
			fontSize: 16,
			color: theme.text.dark.primary,
		},
		description: {
			fontSize: 14,
			color: theme.text.dark.secondary,
		},
		textContainer: {
			display: 'flex',
			flexDirection: 'column',
		},
		archive: {
		},
	};

	render() {
		return (
			<li
				style={Task.styles.tile}
				onMouseEnter={this.onMouseEnter}
				onMouseLeave={this.onMouseLeave}
			>
				<div style={Task.styles.textContainer}>
					<span style={Task.styles.title}>{this.props.task.title}</span>
					<span style={Task.styles.description}>{this.props.task.description}</span>
				</div>
				{this.state.isHovering ?
					<a href="#" style={Task.styles.archive} onClick={this.onArchiveClick}>
						{this.props.task.isArchived ?
							<Icon name="unarchive"/> :
							<Icon name="archive"/>
						}
					</a> :
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
