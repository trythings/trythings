import React from 'react';
import Relay from 'react-relay';

import Icon from './Icon.js';
import resetStyles from './resetStyles.js';
import SearchField from './SearchField.js';
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

class AppBar extends React.Component {
	static propTypes = {
		initialSearchQuery: React.PropTypes.string,
		onSearchQueryChange: React.PropTypes.func,
	};

	static styles = {
		appBar: {
			...resetStyles,
			...theme.elevation[4],

			alignItems: 'center',
			backgroundColor: theme.colors.primary.default,
			height: 56,
			justifyContent: 'space-between',
			minHeight: 56,
			paddingLeft: 16,
			paddingRight: 16,
		},
		title: {
			...resetStyles,
			...theme.text.light.primary,

			fontSize: 20,
			width: 240 - 16, // Align with the navigation drawer.
		},
		spacer: {
			...resetStyles,
			paddingLeft: 24,
		},
		searchField: {
			...resetStyles,
			flex: '1 0 auto',
		},
		migrateButton: {
			...resetStyles,

			borderRadius: '50%',

			paddingBottom: 8,
			paddingLeft: 8,
			paddingRight: 8,
			paddingTop: 8,
		},
		migrateIcon: {
			...resetStyles,
			...theme.text.light.primary,
		},
	};

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

	render() {
		return (
			<div style={AppBar.styles.appBar}>
				<span style={AppBar.styles.title}>Ellie's Pad</span>

				<div style={AppBar.styles.spacer} />

				<SearchField
					initialQuery={this.props.initialSearchQuery}
					onQueryChange={this.props.onSearchQueryChange}
					style={AppBar.styles.searchField}
				/>

				<div style={AppBar.styles.spacer} />

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
					<Icon style={AppBar.styles.migrateIcon} name="update" />
				</button>
			</div>
		);
	}
}

export default Relay.createContainer(AppBar, {
	fragments: {},
});
