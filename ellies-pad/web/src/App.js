import 'normalize.css';
import React from 'react';
import Relay from 'react-relay';

import './Roboto.css';

import ActionButton from './ActionButton.js';
import AddTaskCard from './AddTaskCard.js';
import Icon from './Icon.js';
import TaskList from './TaskList.js';
import theme from './theme.js';

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
		isMigrateHovering: false,
	};

	onMigrateClick = () => {
		Relay.Store.commitUpdate(
			new MigrateMutation({}),
		);
	};

	onCancelClick = () => {
		this.setState({ isAddTaskFormVisible: false });
	};

	onPlusClick = () => {
		this.setState({ isAddTaskFormVisible: true });
	};

	onMigrateMouseEnter = () => {
		this.setState({ isMigrateHovering: true });
	};

	onMigrateMouseLeave = () => {
		this.setState({ isMigrateHovering: false });
	};

	static styles = {
		app: {
			backgroundColor: theme.colors.canvas,
			display: 'flex',
			flex: 1,
			flexDirection: 'column',
			overflowX: 'hidden',
		},
		appBar: {
			backgroundColor: theme.colors.primary,

			alignItems: 'center',
			justifyContent: 'space-between',

			boxSizing: 'border-box',
			display: 'flex',
			height: 56,
			minHeight: 56,
			paddingLeft: 16,
			paddingRight: 16,

			// App bar hovers above other sheets.
			boxShadow: [
				'0 1px 5px 0 rgba(0, 0, 0, 0.12)', // Ambient.
				'0 2px 2px 0 rgba(0, 0, 0, 0.14)', // Penumbra.
				'0 3px 1px -2px rgba(0, 0, 0, 0.20)', // Umbra.
			].join(','),
			zIndex: 4,

			title: {
				// Light primary text.
				color: theme.text.light.color,
				opacity: theme.text.light.opacity.primary,

				// Title text.
				fontFamily: 'Roboto, sans-serif',
				fontSize: 20,
				fontWeight: 600,
				lineHeight: '44px',
			},
			migrateButton: {
				display: 'flex',
				alignItems: 'center',

				// backgroundColor: 'rgba(255, 255, 255, 0)',
				border: 'none',
				outline: 0,

				borderRadius: '50%',
				padding: 8,
			},
		},
		addTaskButton: {
			position: 'absolute',
			right: 0,
		},
		content: {
			display: 'flex',
			flexDirection: 'column',
			overflow: 'scroll',
			padding: 24,
		},
		contentSpacer: {
			padding: 12,
		},
	};

	render() {
		return (
			<div style={App.styles.app}>
				<div style={App.styles.appBar}>
					<span style={App.styles.appBar.title}>Ellie's Pad</span>

					<button
						style={{
							...App.styles.appBar.migrateButton,
							backgroundColor: this.state.isMigrateHovering ?
								'rgba(255, 255, 255, 0.12)' :
								'rgba(255, 255, 255, 0)',
						}}
						onClick={this.onMigrateClick}
						onMouseEnter={this.onMigrateMouseEnter}
						onMouseLeave={this.onMigrateMouseLeave}
					>
						<Icon color={theme.text.light.primary} name="update"/>
					</button>
				</div>

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

					<div style={App.styles.contentSpacer}/>

					<TaskList viewer={this.props.viewer}/>
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
				${TaskList.getFragment('viewer')},
			}
		`,
	},
});
