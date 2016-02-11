import React from 'react';
import Relay from 'react-relay';

import Icon from './Icon.js';
import resetStyles from './resetStyles.js';
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
			...resetStyles,
			paddingLeft: 16,
			paddingRight: 16,
			height: 60,
			alignItems: 'center',
			justifyContent: 'space-between',
		},
		title: {
			...resetStyles,
			...theme.text.dark.primary,

			display: 'inline',
			fontSize: 16,
			textOverflow: 'ellipsis',
			whiteSpace: 'nowrap',
		},
		description: {
			...resetStyles,
			...theme.text.dark.secondary,

			display: 'inline',
			fontSize: 14,
			textOverflow: 'ellipsis',
			whiteSpace: 'nowrap',
		},
		textContainer: {
			...resetStyles,
			flexDirection: 'column',
		},
		archive: {
			...resetStyles,
			justifyContent: 'center',
		},
	};

	renderText() {
		const title = <span style={Task.styles.title}>{this.props.task.title}</span>;
		const description = <span style={Task.styles.description}>{this.props.task.description}</span>;
		if (this.props.task.isArchived) {
			return <del style={Task.styles.textContainer}>{title}{description}</del>;
		}
		return <div style={Task.styles.textContainer}>{title}{description}</div>;
	}

	render() {
		return (
			<li
				style={Task.styles.tile}
				onMouseEnter={this.onMouseEnter}
				onMouseLeave={this.onMouseLeave}
			>
				{this.renderText()}
				{this.state.isHovering ?
					<button style={Task.styles.archive} onClick={this.onArchiveClick}>
						{this.props.task.isArchived ?
							<Icon color={theme.text.dark.secondary.color} name="unarchive"/> :
							<Icon color={theme.text.dark.secondary.color} name="archive"/>
						}
					</button> :
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
