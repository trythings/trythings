import React from 'react';
import Relay from 'react-relay';

class AddTaskMutation extends Relay.Mutation {
	static fragments = {
		viewer: () => Relay.QL`
			fragment on User {
				id,
				tasks,
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
		viewer: React.PropTypes.shape({
			// ...AddTaskMutation.propTypes.viewer
		}).isRequired,
		onCancelClick: React.PropTypes.func,
	};

	state = {
		tags: '',
		title: '',
		description: '',
	};

	onTagsChange = (event) => {
		this.setState({ tags: event.target.value });
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
				title: `${this.state.tags} ${this.state.title}`,
				description: this.state.description || null,
				viewer: this.props.viewer,
			}),
		);
		this.setState({
			tags: '',
			title: '',
			description: '',
		});
	};

	static styles = {
		card: {
			backgroundColor: '#ffffff',
			boxShadow: [
				'0 1px 5px 0 rgba(0, 0, 0, 0.12)', // Ambient.
				'0 2px 2px 0 rgba(0, 0, 0, 0.14)', // Penumbra.
				'0 3px 1px -2px rgba(0, 0, 0, 0.20)', // Umbra.
			].join(','),
			display: 'flex',
			minHeight: 'min-content',
			minWidth: 'min-content',
			zIndex: 6,
		},
	};

	render() {
		return (
			<div style={AddTaskCard.styles.card}>
				<input
					placeholder="Tags"
					value={this.state.tags}
					onChange={this.onTagsChange}
				/>
				<input
					placeholder="Title"
					value={this.state.title}
					onChange={this.onTitleChange}
				/>
				<textarea
					placeholder="Description"
					value={this.state.description}
					onChange={this.onDescriptionChange}
				/>

				<div>
					<button onClick={this.props.onCancelClick}>
						Cancel
					</button>
					<button onClick={this.onAddClick}>
						Add Task
					</button>
				</div>
			</div>
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
