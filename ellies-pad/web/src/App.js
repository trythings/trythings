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
		};
	}
	
	getOptimisticResponse() {
		return {
			task: {
				title: this.props.title,
			},
			viewer: {
				id: this.props.viewer.id,
			},
		};
	}
}

const Task = Relay.createContainer((props) => {
	return <li>{props.task.title}</li>;
}, {
	fragments: {
		task: () => Relay.QL`
			fragment on Task {
				title,
			}
		`,
	},
});


class App extends React.Component {
	state = {
		title: '',
	};

	render() {
		return (
			<div>
				<input
					placeholder="New task"
					value={this.state.title}
					onChange={(event) => {
						this.setState({title: event.target.value});
					}}
					onKeyPress={(event) => {
						if (event.key === 'Enter') {
							Relay.Store.commitUpdate(
								new AddTaskMutation({
									title: event.target.value,
									viewer: this.props.viewer,
								}),
							);
							this.setState({title: ''});
						}
					}}
				/>
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
