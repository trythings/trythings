import React from 'react';
import Relay from 'react-relay';
import TextareaAutosize from 'react-textarea-autosize';

import Card from './Card.js';
import FlatButton from './FlatButton.js';
import resetStyles from './resetStyles.js';
import theme from './theme.js';

class EditTaskMutation extends Relay.Mutation {
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
				editTask,
			}
		`;
	}

	getFatQuery() {
		return Relay.QL`
			fragment on EditTaskPayload {
				task {
					title,
					description,
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
			title: this.props.title,
			description: this.props.description,
			// isArchived: this.props.isArchived,
		};
	}

	getOptimisticResponse() {
		return {
			task: {
				id: this.props.task.id,
				title: this.props.title,
				description: this.props.description,
			},
		};
	}
}

class TaskCard extends React.Component {
	static propTypes = {
		autoFocus: React.PropTypes.bool,
		task: React.PropTypes.shape({
			title: React.PropTypes.string.isRequired,
			description: React.PropTypes.string,
			isArchived: React.PropTypes.boolean,
		}).isRequired,
		onClose: React.PropTypes.func,
	};

	state = {
		title: this.props.task.title,
		description: this.props.task.description,
	};

	onTitleChange = (event) => {
		this.setState({ title: event.target.value });
	};

	onDescriptionChange = (event) => {
		this.setState({ description: event.target.value });
	};

	requestClose = () => {
		const hasChanges = this.state.title !== this.props.task.title ||
			this.state.description !== this.props.task.description;
		if (hasChanges) {
			if (confirm('There are unsaved changes for this task. Would you like to save them?')) {
				this.saveChanges();
			}
		}
		this.props.onClose();
	};

	saveChanges = () => {
		Relay.Store.commitUpdate(
			new EditTaskMutation({
				task: this.props.task,
				title: this.state.title,
				description: this.state.description,
			}),
		);
	};

	static styles = {
		header: {
			...resetStyles,
			flexDirection: 'column',
			padding: 16,
		},
		title: {
			...resetStyles,
			...theme.text.dark.primary,

			fontSize: 24,
			fontWeight: 300,
		},
		description: {
			...resetStyles,
			...theme.text.dark.secondary,
			minHeight: 0,
		},
		titleSpacer: {
			...resetStyles,
			paddingTop: 16,
		},
		actionContainer: {
			...resetStyles,
			paddingBottom: 8,
			paddingLeft: 8,
			paddingRight: 8,
			paddingTop: 8,
		},
		actionSpacer: {
			...resetStyles,
			paddingLeft: 8,
		},
	};

	render() {
		return (
			<Card autoFocus={this.props.autoFocus}>
				<header style={TaskCard.styles.header}>
					<input
						placeholder="Title"
						value={this.state.title}
						onChange={this.onTitleChange}
						style={TaskCard.styles.title}
					/>

					<div style={TaskCard.styles.titleSpacer}/>

					<TextareaAutosize
						placeholder="Description"
						value={this.state.description}
						onChange={this.onDescriptionChange}
						style={TaskCard.styles.description}
					/>
				</header>

				<div style={TaskCard.styles.actionContainer}>
					<FlatButton
						color={theme.text.dark.primary.color}
						onClick={this.requestClose}
						label="Cancel"
					/>

					<div style={TaskCard.styles.actionSpacer}/>

					<FlatButton
						color={theme.colors.accentLight}
						onClick={this.saveChanges}
						label="Save"
					/>
				</div>
			</Card>
		);
	}
}

export default Relay.createContainer(TaskCard, {
	fragments: {
		task: () => Relay.QL`
			fragment on Task {
				${EditTaskMutation.getFragment('task')},
				title,
				description,
				isArchived,
			}
		`,
	},
});
