import React from 'react';
import Relay from 'react-relay';

import Task from './Task.js';

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

class App extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			// ...AddTaskMutation.propTypes.viewer
			tasks: React.PropTypes.arrayOf(React.PropTypes.object).isRequired, // ...Task.propTypes.task
		}).isRequired,
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

	render() {
		return (
			<div>
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
				<button onClick={this.onAddClick}>
					Add task
				</button>

				<ol>
					{this.props.viewer.tasks.map(task => <Task key={task.id} task={task}/>)}
				</ol>
			</div>
		);
	}
}

export default Relay.createContainer(App, {
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				${AddTaskMutation.getFragment('viewer')},
				tasks {
					id,
					${Task.getFragment('task')},
				},
			}
		`,
	},
});
