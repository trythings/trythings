import _pick from 'lodash/pick';
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
				editTask,
			}
		`;
	}

	getFatQuery() {
		return Relay.QL`
			fragment on EditTaskPayload {
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
			id: this.props.task.id,
			isArchived: !this.props.task.isArchived,
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

class TaskTile extends React.Component {
	static propTypes = {
		style: React.PropTypes.shape({
			flex: React.PropTypes.string,
		}),
		task: React.PropTypes.shape({
			title: React.PropTypes.string.isRequired,
			description: React.PropTypes.string,
			isArchived: React.PropTypes.bool,
		}).isRequired,
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
		textContainer: {
			...resetStyles,
			alignItems: 'stretch',
			flex: '0 1 auto',
			flexDirection: 'column',
		},
		title: {
			...resetStyles,
			...theme.text.dark.primary,

			display: 'block',
			fontSize: 16,
			textOverflow: 'ellipsis',
			whiteSpace: 'nowrap',
		},
		description: {
			...resetStyles,
			...theme.text.dark.secondary,

			display: 'block',
			fontSize: 14,
			textOverflow: 'ellipsis',
			whiteSpace: 'nowrap',
		},
		archive: {
			...resetStyles,
			justifyContent: 'center',
		},
		archiveIcon: {
			...resetStyles,
			...theme.text.dark.secondary,
		},
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

	onArchiveFocus = (event) => {
		event.stopPropagation();
	};

	renderText() {
		const title = <span style={TaskTile.styles.title}>{this.props.task.title}</span>;
		const description = (
			<span style={TaskTile.styles.description}>
				{this.props.task.description}
			</span>
		);
		if (this.props.task.isArchived) {
			return <del style={TaskTile.styles.textContainer}>{title}{description}</del>;
		}
		return <div style={TaskTile.styles.textContainer}>{title}{description}</div>;
	}

	render() {
		return (
			<div
				style={{
					...TaskTile.styles.tile,
					..._pick(this.props.style, ['flex']),
				}}
				onMouseEnter={this.onMouseEnter}
				onMouseLeave={this.onMouseLeave}
				tabIndex={-1}
			>
				{this.renderText()}
				{this.state.isHovering ?
					<button
						style={TaskTile.styles.archive}
						onClick={this.onArchiveClick}
						onFocus={this.onArchiveFocus}
					>
						{this.props.task.isArchived ?
							<Icon style={TaskTile.styles.archiveIcon} name="unarchive" /> :
							<Icon style={TaskTile.styles.archiveIcon} name="archive" />
						}
					</button> :
					null
				}
			</div>
		);
	}
}

export default Relay.createContainer(TaskTile, {
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
