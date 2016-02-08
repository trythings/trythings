import 'normalize.css';
import React from 'react';
import Relay from 'react-relay';

import './Roboto.css';

import Icon from './Icon.js';
import Task from './Task.js';


// TODO: This is a temporary solution to enable us to run all of our migrations.
class MigrateMutation extends Relay.Mutation {
	static fragments = {};

	getMutation() {
		return Relay.QL`
			mutation {
				migrate,
			}
		`;
	}

	// It's unclear how to specify a fragment with no fields.
	// We use the clientMutationId to give this fragment > 0 fields.
	getFatQuery() {
		return Relay.QL`
			fragment on MigratePayload {
				clientMutationId,
			}
		`;
	}

	getConfigs() {
		return [];
	}

	getVariables() {
		return {};
	}

	getOptimisticResponse() {
		return {};
	}
}

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

const deepPurple = {
	500: '#673ab7',
};

const colors = {
	primary1: deepPurple['500'],
	// primary2,
	// primary3,
	// accent,
};

const text = {
	light: {
		color: '#ffffff',
		opacity: {
			primary: '100%',
			secondary: '70%',
			disabled: '50%',
			dividers: '12%',
		},
	},
	dark: {
		color: '#000000',
		opacity: {
			primary: '87%',
			secondary: '54%',
			disabled: '38%',
			dividers: '12%',
		},
	},
};

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

		isAddTaskFormVisible: true,
		isMigrateHovering: false,
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

	onMigrateClick = () => {
		Relay.Store.commitUpdate(
			new MigrateMutation({}),
		);
	};

	onCancelClick = () => {
		this.setState({ isAddTaskFormVisible: false });
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

	onMigrateMouseEnter = () => {
		this.setState({ isMigrateHovering: true });
	};

	onMigrateMouseLeave = () => {
		this.setState({ isMigrateHovering: false });
	};

	static styles = {
		app: {
			backgroundColor: '#fafafa',
			display: 'flex',
			flex: 1,
			flexDirection: 'column',
		},
		appBar: {
			backgroundColor: colors.primary1,

			alignItems: 'center',
			justifyContent: 'space-between',

			display: 'flex',
			height: 56,
			minHeight: 'min-content',
			minWidth: 'min-content',
			paddingLeft: 16,
			paddingRight: 16,

			// position: 'fixed',

			// App bar hovers above other sheets.
			boxShadow: '0 2px 5px rgba(0, 0, 0, 0.26)',
			zIndex: 1,

			title: {
				// Light primary text.
				color: text.light.color,
				opacity: text.light.opacity.primary,

				// Title text.
				fontFamily: 'Roboto, sans-serif',
				fontSize: 20,
				fontWeight: 600,
				lineHeight: '44px',
			},
			migrate: {
				display: 'flex',
				alignItems: 'center',

				// backgroundColor: 'rgba(255, 255, 255, 0)',
				border: 'none',
				borderRadius: 24,

				color: text.light.color,
				opacity: text.light.opacity.primary,
				outline: 0,

				padding: 8,
			},
		},
		addTaskForm: {
			backgroundColor: '#ffffff',
			boxShadow: '0 2px 5px rgba(0, 0, 0, 0.26)',
			display: 'flex',
			minHeight: 'min-content',
			minWidth: 'min-content',
			zIndex: 1,
		},
	};

	render() {
		return (
			<div style={App.styles.app}>
				<div style={App.styles.appBar}>
					<span style={App.styles.appBar.title}>Ellie's Pad</span>

					<button
						style={{
							...App.styles.appBar.migrate,
							backgroundColor: this.state.isMigrateHovering ?
								'rgba(255, 255, 255, 0.12)' :
								'rgba(255, 255, 255, 0)',
						}}
						onClick={this.onMigrateClick}
						onMouseEnter={this.onMigrateMouseEnter}
						onMouseLeave={this.onMigrateMouseLeave}
					>
						<Icon name="update"/>
					</button>
				</div>

				{this.state.isAddTaskFormVisible ?
					(
						<div style={App.styles.addTaskForm}>
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
								<button onClick={this.onCancelClick}>
									Cancel
								</button>
								<button onClick={this.onAddClick}>
									Add Task
								</button>
							</div>
						</div>
					) :
					null
				}

				<div>
					<ol>
						{this.props.viewer.tasks.map(task => <Task key={task.id} task={task}/>)}
					</ol>
				</div>
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
