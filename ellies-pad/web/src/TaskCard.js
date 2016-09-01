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
					body,
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
			body: this.props.body,
			// isArchived: this.props.isArchived,
		};
	}

	getOptimisticResponse() {
		return {
			task: {
				id: this.props.task.id,
				title: this.props.title,
				body: this.props.body,
			},
		};
	}
}

class TaskCard extends React.Component {
	static propTypes = {
		autoFocus: React.PropTypes.bool,
		style: React.PropTypes.shape({
			flex: React.PropTypes.string,
		}),
		task: React.PropTypes.shape({
			title: React.PropTypes.string.isRequired,
			body: React.PropTypes.string,
			isArchived: React.PropTypes.bool,
		}).isRequired,
		onClose: React.PropTypes.func,
	};

	static styles = {
		header: {
			...resetStyles,
			alignItems: 'stretch',
			flexDirection: 'column',
			padding: 16,
		},
		title: {
			...resetStyles,
			...theme.text.dark.primary,

			fontSize: 24,
			fontWeight: 300,
		},
		body: {
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
		save: {
			...resetStyles,
			color: theme.colors.accent,
		},
	};

	state = {
		title: this.props.task.title,
		body: this.props.task.body,
	};

	onTitleChange = (event) => {
		this.setState({ title: event.target.value });
	};

	onbodyChange = (event) => {
		this.setState({ body: event.target.value });
	};

	onSaveClick = () => {
		if (this.hasUnsavedChanges()) {
			this.saveChanges();
		}
		this.props.onClose();
	};

	hasUnsavedChanges() {
		return this.state.title !== this.props.task.title ||
			this.state.body !== this.props.task.body;
	}

	requestClose = () => {
		if (this.hasUnsavedChanges()) {
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
				body: this.state.body,
			}),
		);
	};

	render() {
		return (
			<Card autoFocus={this.props.autoFocus} style={this.props.style}>
				<header style={TaskCard.styles.header}>
					<input
						placeholder="Title"
						value={this.state.title}
						onChange={this.onTitleChange}
						style={TaskCard.styles.title}
					/>

					<div style={TaskCard.styles.titleSpacer} />

					<TextareaAutosize
						placeholder="Description"
						value={this.state.body}
						onChange={this.onDescriptionChange}
						style={TaskCard.styles.body}
					/>
				</header>

				<div style={TaskCard.styles.actionContainer}>
					<FlatButton
						style={TaskCard.styles.save}
						onClick={this.onSaveClick}
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
				body,
				isArchived,
			}
		`,
	},
});
