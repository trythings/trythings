
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

class AppBar extends React.Component {
	static propTypes = {
		children: React.PropTypes.node,
		style: React.PropTypes.shape({
			backgroundColor: React.PropTypes.string,
			color: React.PropTypes.string,
		}),
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
		},
		children: {
			...resetStyles,
			flex: '1 0 auto',
			paddingLeft: 24,
			paddingRight: 24,
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

	renderMigrateButton() {
		let migrateIconStyle = AppBar.styles.migrateIcon;
		if (this.props.style && this.props.style.color) {
			migrateIconStyle = {
				...migrateIconStyle,
				color: this.props.style.color,
			};
		}

		return (
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
				<Icon style={migrateIconStyle} name="update" />
			</button>
		);
	}

	render() {
		let style = AppBar.styles.appBar;
		if (this.props.style && this.props.style.backgroundColor) {
			style = {
				...style,
				backgroundColor: this.props.style.backgroundColor,
			};
		}

		let titleStyle = AppBar.styles.title;
		if (this.props.style && this.props.style.color) {
			titleStyle = {
				...titleStyle,
				color: this.props.style.color,
			};
		}

		return (
			<div style={style}>
				<span style={titleStyle}>Ellie's Pad</span>

				<div style={AppBar.styles.children}>
					{this.props.children}
				</div>

				{this.renderMigrateButton()}
			</div>
		);
	}
}

export default Relay.createContainer(AppBar, {
	fragments: {},
});
