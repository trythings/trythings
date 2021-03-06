import _pick from 'lodash/pick';
import React from 'react';
import Relay from 'react-relay';
import TextareaAutosize from 'react-textarea-autosize';

import FlatButton from './FlatButton.js';
import resetStyles from './resetStyles.js';
import theme from './theme.js';

class AddTaskMutation extends Relay.Mutation {
	static fragments = {
		space: () => Relay.QL`
			fragment on Space {
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
			}
		`;
	}

	getConfigs() {
		return [];
	}

	getVariables() {
		return {
			title: this.props.title,
			description: this.props.description,
			spaceId: this.props.space.id,
		};
	}

	getOptimisticResponse() {
		return {
			task: {
				title: this.props.title,
				description: this.props.description,
			},
		};
	}
}

class AddTaskMole extends React.Component {
	static propTypes = {
		autoFocus: React.PropTypes.bool,
		refetch: React.PropTypes.func.isRequired,
		space: React.PropTypes.shape({
			// ...AddTaskMutation.propTypes.space
		}).isRequired,
		onCancelClick: React.PropTypes.func,
		style: React.PropTypes.shape({
			flex: React.PropTypes.string,
		}),
	};

	static styles = {
		container: {
			...resetStyles,
			...theme.elevation[12],

			alignItems: 'stretch',
			backgroundColor: theme.colors.card,
			borderRadius: '5px 5px 0 0',
			flexDirection: 'column',

			width: 360,
		},
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
		cancel: {
			...resetStyles,
			...theme.text.dark.primary,
		},
		addTask: {
			...resetStyles,
			color: theme.colors.accent,
		},
	};

	state = {
		title: '',
		description: '',
	};

	componentWillUnmount() {
		if (this.timeout) {
			clearTimeout(this.timeout);
		}
	}

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
				space: this.props.space,
			}),
		);
		this.setState({
			title: '',
			description: '',
		});
		this.timeout = setTimeout(() => {
			this.props.refetch();
			this.timeout = null;
		}, 1000);
	};

	render() {
		return (
			<div
				style={{
					...AddTaskMole.styles.container,
					..._pick(this.props.style, ['flex']),
				}}
			>
				<header style={AddTaskMole.styles.header}>
					<input
						autoFocus={this.props.autoFocus}
						placeholder="Title"
						value={this.state.title}
						onChange={this.onTitleChange}
						style={AddTaskMole.styles.title}
					/>

					<div style={AddTaskMole.styles.titleSpacer} />

					<TextareaAutosize
						placeholder="Description"
						value={this.state.description}
						onChange={this.onDescriptionChange}
						style={AddTaskMole.styles.description}
					/>
				</header>

				<div style={AddTaskMole.styles.actionContainer}>
					<FlatButton
						style={AddTaskMole.styles.cancel}
						onClick={this.props.onCancelClick}
						label="Cancel"
					/>

					<div style={AddTaskMole.styles.actionSpacer} />

					<FlatButton
						style={AddTaskMole.styles.addTask}
						onClick={this.onAddClick}
						label="Add Task"
					/>
				</div>
			</div>
		);
	}
}

export default Relay.createContainer(AddTaskMole, {
	fragments: {
		space: () => Relay.QL`
			fragment on Space {
				${AddTaskMutation.getFragment('space')},
			}
		`,
	},
});
