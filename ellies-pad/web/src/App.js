import React from 'react';
import Relay from 'react-relay';

import ActionButton from './ActionButton.js';
import AddTaskCard from './AddTaskCard.js';
import AppBar from './AppBar.js';
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
			backgroundColor: theme.colors.canvas,
			display: 'flex',
			flex: 1,
			flexDirection: 'column',
			overflowX: 'hidden',
		},
		addTaskButton: {
			position: 'absolute',
			right: 24,
		},
		contentScroll: {
			overflow: 'scroll',
		},
		content: {
			display: 'flex',
			flexDirection: 'column',
			padding: 24,
			minHeight: 'min-content',
		},
		contentSpacer: {
			padding: 12,
		},
	};

	render() {
		return (
			<div style={App.styles.app}>
				<AppBar/>

				<div style={App.styles.contentScroll}>
					<div style={App.styles.content}>

						{this.state.isAddTaskFormVisible ?
							(
								<AddTaskCard
									viewer={this.props.viewer}
									onCancelClick={this.onCancelClick}
								/>
							) :
							(
								<div style={App.styles.addTaskButton}>
									<ActionButton onClick={this.onPlusClick}/>
								</div>
							)
						}

						{this.state.isAddTaskFormVisible ?
							<div style={App.styles.contentSpacer}/> :
							null
						}

						<TaskSearch name="#now" query="#now"/>
						<div style={App.styles.contentSpacer}/>

						<TaskSearch name="incoming" query="NOT #now AND NOT #next AND NOT #later"/>
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
