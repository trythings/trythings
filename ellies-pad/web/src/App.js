import React from 'react';
import Relay from 'react-relay';

import ActionButton from './ActionButton.js';
import AddTaskCard from './AddTaskCard.js';
import AppBar from './AppBar.js';
import resetStyles from './resetStyles.js';
import TaskSearch from './TaskSearch.js';
import theme from './theme.js';

class App extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			// ...AddTaskCard.propTypes.viewer
			// ...TaskList.propTypes.viewer
		}).isRequired,
	};

	state = {
		tags: '',
		title: '',
		description: '',

		isAddTaskFormVisible: true,
	};

	onCancelClick = () => {
		this.setState({ isAddTaskFormVisible: false });
	};

	onPlusClick = () => {
		this.setState({ isAddTaskFormVisible: true });
	};

	static styles = {
		app: {
			...resetStyles,
			backgroundColor: theme.colors.canvas,
			flexDirection: 'column',
			height: '100%',
			width: '100%',
		},
		container: {
			...resetStyles,
			flex: '1 1 auto',
		},
		addTaskButton: {
			...resetStyles,
			overflow: 'visible',
			position: 'absolute',

			right: 24,
			top: 24,
		},
		content: {
			...resetStyles,
			flex: '1 1 auto',
			flexDirection: 'column',
			padding: 24,
			overflow: 'scroll',
		},
		contentSpacer: {
			...resetStyles,
			padding: 12,
		},
	};

	render() {
		return (
			<div style={App.styles.app}>
				<AppBar/>

				<div style={App.styles.container}>
					{!this.state.isAddTaskFormVisible ?
						(
							<div style={App.styles.addTaskButton}>
								<ActionButton onClick={this.onPlusClick}/>
							</div>
						) :
						null
					}
					<div style={App.styles.content}>
						{this.state.isAddTaskFormVisible ?
							(
								<AddTaskCard
									autoFocus
									viewer={this.props.viewer}
									onCancelClick={this.onCancelClick}
								/>
							) :
							null
						}

						{this.state.isAddTaskFormVisible ?
							<div style={App.styles.contentSpacer}/> :
							null
						}

						<TaskSearch name="#now" query="#now"/>
						<div style={App.styles.contentSpacer}/>

						<TaskSearch name="Incoming" query="NOT #now AND NOT #next AND NOT #later"/>
						<div style={App.styles.contentSpacer}/>

						<TaskSearch name="#next" query="#next AND NOT #now"/>
						<div style={App.styles.contentSpacer}/>

						<TaskSearch name="#later" query="#later AND NOT #next AND NOT #now"/>
						<div style={App.styles.contentSpacer}/>
					</div>
				</div>
			</div>
		);
	}
}

export default Relay.createContainer(App, {
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				${AddTaskCard.getFragment('viewer')},
			}
		`,
	},
});
