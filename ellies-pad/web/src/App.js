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
				addTask
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

const Task = Relay.createContainer((props) => {
	return (
		<li>
			<strong>{props.task.title}</strong> ≫ <span>{props.task.description}</span>
		</li>
	);
}, {
	fragments: {
		task: () => Relay.QL`
			fragment on Task {
				title,
				description,
			}
		`,
	},
});

class App extends React.Component {
	state = {
		tags: '',
		title: '',
		description: '',
	};

	render() {
		return (
			<div>
				<input
					placeholder="Tags"
					value={this.state.tags}
					onChange={(event) => {
						this.setState({tags: event.target.value});
					}}
				/>
				<input
					placeholder="Title"
					value={this.state.title}
					onChange={(event) => {
						this.setState({title: event.target.value});
					}}
				/>
				<textarea
					placeholder="Description"
					value={this.state.description}
					onChange={(event) => {
						this.setState({description: event.target.value});
					}}
				/>
				<button
					onClick={(event) => {
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
					}}
				>
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
