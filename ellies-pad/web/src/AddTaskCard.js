import React from 'react';
import Relay from 'react-relay';
import TextareaAutosize from 'react-textarea-autosize';

import Card from './Card.js';
import FlatButton from './FlatButton.js';
import resetStyles from './resetStyles.js';
import theme from './theme.js';

class AddTaskMutation extends Relay.Mutation {
	static fragments = {
		viewer: () => Relay.QL`
			fragment on User {
				id,
			}
		`,
	};

	getMutation() {
		return Relay.QL`
			mutation {
				addTask,
			}
		`;
	}

	getFatQuery() {
		return Relay.QL`
			fragment on AddTaskPayload {
				task,
				viewer {
					tasks,
				},
			}
		`;
	}

	getConfigs() {
		return [{
			type: 'FIELDS_CHANGE',
			fieldIDs: {
				viewer: this.props.viewer.id,
			},
		}];
	}

	getVariables() {
		return {
			title: this.props.title,
			description: this.props.description,
		};
	}

	getOptimisticResponse() {
		return {
			task: {
				title: this.props.title,
				description: this.props.description,
			},
			viewer: {
				id: this.props.viewer.id,
			},
		};
	}
}

class AddTaskCard extends React.Component {
	static propTypes = {
		autoFocus: React.PropTypes.bool,
		viewer: React.PropTypes.shape({
			// ...AddTaskMutation.propTypes.viewer
		}).isRequired,
		onCancelClick: React.PropTypes.func,
	};

	state = {
		title: '',
		description: '',
	};

	onTitleChange = (event) => {
		this.setState({ title: event.target.value });
	};

	onDescriptionChange = (event) => {
		this.setState({ description: event.target.value });
	};

	onAddClick = () => {
		Relay.Store.commitUpdate(
			new AddTaskMutation({
				title: `${this.state.title}`,
				description: this.state.description || null,
				viewer: this.props.viewer,
			}),
		);
		this.setState({
			title: '',
			description: '',
		});
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
			<Card>
				<header style={AddTaskCard.styles.header}>
					<input
						autoFocus={this.props.autoFocus}
						placeholder="Title"
						value={this.state.title}
						onChange={this.onTitleChange}
						style={AddTaskCard.styles.title}
					/>

					<div style={AddTaskCard.styles.titleSpacer}/>

					<TextareaAutosize
						placeholder="Description"
						value={this.state.description}
						onChange={this.onDescriptionChange}
						style={AddTaskCard.styles.description}
					/>
				</header>

				<div style={AddTaskCard.styles.actionContainer}>
					<FlatButton
						color={theme.text.dark.primary.color}
						onClick={this.props.onCancelClick}
						label="Cancel"
					/>

					<div style={AddTaskCard.styles.actionSpacer}/>

					<FlatButton
						color={theme.colors.accentLight}
						onClick={this.onAddClick}
						label="Add Task"
					/>
				</div>
			</Card>
		);
	}
}

export default Relay.createContainer(AddTaskCard, {
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				${AddTaskMutation.getFragment('viewer')},
			}
		`,
	},
});
