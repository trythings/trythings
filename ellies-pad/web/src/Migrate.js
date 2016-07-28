import gapi from 'gapi';
import React from 'react';
import Relay from 'react-relay';

import Icon from './Icon.js';
import resetStyles from './resetStyles.js';
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

class Migrate extends React.Component {
	static contextTypes = {
		router: React.PropTypes.shape({
			replace: React.PropTypes.func.isRequired,
		}).isRequired,
	};

	// Routing.
	static onEnter = (nextState, replace) => {
		if (!gapi.auth2.getAuthInstance().isSignedIn.get()) {
			replace('/signin');
		}
	};

	static propTypes = {
		style: React.PropTypes.shape({
			color: React.PropTypes.string,
		}),
		viewer: React.PropTypes.shape({
			isAdmin: React.PropTypes.bool.isRequired,
		}),
	};

	static styles = {
		container: {
			...resetStyles,
			flex: '1 1 auto',
		},
		buttonContainer: {
			...resetStyles,
			alignItems: 'center',
			backgroundColor: theme.colors.primary.default,
			paddingRight: 16,
		},
		button: {
			...resetStyles,

			borderRadius: '50%',

			paddingBottom: 8,
			paddingLeft: 8,
			paddingRight: 8,
			paddingTop: 8,
		},
		icon: {
			...resetStyles,
			...theme.text.light.primary,
		},
		label: {
			...resetStyles,
			...theme.text.light.primary,
		},
	};

	state = {
		isHovering: false,
	};

	componentWillMount() {
		if (!this.props.viewer.isAdmin) {
			this.context.router.replace('/');
		}
	}

	onClick = () => {
		Relay.Store.commitUpdate(
			new MigrateMutation({}),
		);
	};

	onEnter = () => {
		this.setState({ isHovering: true });
	};

	onLeave = () => {
		this.setState({ isHovering: false });
	};

	render() {
		let iconStyle = Migrate.styles.icon;
		if (this.props.style && this.props.style.color) {
			iconStyle = {
				...iconStyle,
				color: this.props.style.color,
			};
		}

		return (
			<div style={Migrate.styles.container}>
				<div style={Migrate.styles.buttonContainer}>
					<button
						style={{
							...Migrate.styles.button,
							backgroundColor: this.state.isHovering ?
								'rgba(255, 255, 255, 0.12)' :
								'rgba(255, 255, 255, 0)',
						}}
						onClick={this.onClick}
						onMouseEnter={this.onEnter}
						onMouseLeave={this.onLeave}
					>
						<Icon style={iconStyle} name="update" />
					</button>
					<span style={Migrate.styles.label}>Run Migrations</span>
				</div>
			</div>
		);
	}
}

const MigrateContainer = Relay.createContainer(Migrate, {
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				isAdmin,
			},
		`,
	},
});

MigrateContainer.onEnter = Migrate.onEnter;

export default MigrateContainer;
