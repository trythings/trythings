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

export default class AppBar extends React.Component {
	state = {
		isMigrateHovering: false,
	};

	onMigrateClick = () => {
		Relay.Store.commitUpdate(
			new MigrateMutation({}),
		);
	};

	onMigrateMouseEnter = () => {
		this.setState({ isMigrateHovering: true });
	};

	onMigrateMouseLeave = () => {
		this.setState({ isMigrateHovering: false });
	};

	static styles = {
		appBar: {
			...resetStyles,
			...theme.elevation[4],

			backgroundColor: theme.colors.primary,

			alignItems: 'center',
			justifyContent: 'space-between',

			height: 56,
			minHeight: 56,
			paddingLeft: 16,
			paddingRight: 16,
		},
		title: {
			...resetStyles,
			...theme.text.light.primary,

			fontSize: 20,
		},
		migrateButton: {
			...resetStyles,

			alignItems: 'center',
			borderRadius: '50%',

			paddingBottom: 8,
			paddingLeft: 8,
			paddingRight: 8,
			paddingTop: 8,
		},
	};

	render() {
		return (
			<div style={AppBar.styles.appBar}>
				<span style={AppBar.styles.title}>Ellie's Pad</span>

				<button
					style={{
						...AppBar.styles.migrateButton,
						backgroundColor: this.state.isMigrateHovering ?
							'rgba(255, 255, 255, 0.12)' :
							'rgba(255, 255, 255, 0)',
					}}
					onClick={this.onMigrateClick}
					onMouseEnter={this.onMigrateMouseEnter}
					onMouseLeave={this.onMigrateMouseLeave}
				>
					<Icon color={theme.text.light.primary.color} name="update"/>
				</button>
			</div>
		);
	}
}