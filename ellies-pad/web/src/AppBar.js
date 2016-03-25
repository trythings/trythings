import _debounce from 'lodash/debounce';
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
	static contextTypes = {
		router: React.PropTypes.object.isRequired,
	};

	static propTypes = {
		searchQuery: React.PropTypes.string,
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
			...theme.text.light.primary,
			backgroundColor: theme.colors.primary.light,
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

	onSearchFocus = () => {
		if (this.props.searchQuery === undefined) {
			this.context.router.push('/search/');
		}
	};

	onSearchBlur = () => {
		if (!this.props.searchQuery) {
			this.context.router.push('/');
		}
	};

	onSearchQueryChange = _debounce((query) => {
		this.context.router.push(`/search/${encodeURIComponent(query)}`);
	}, 200);

	render() {
		let style = AppBar.styles.appBar;
		if (this.props.searchQuery !== undefined) {
			style = {
				...style,
				backgroundColor: theme.colors.card,
			};
		}

		let titleStyle = AppBar.styles.title;
		if (this.props.searchQuery !== undefined) {
			titleStyle = {
				...titleStyle,
				...theme.text.dark.primary,
			};
		}

		let searchFieldStyle = AppBar.styles.searchField;
		if (this.props.searchQuery !== undefined) {
			searchFieldStyle = {
				...searchFieldStyle,
				...theme.text.dark.primary,
				backgroundColor: theme.colors.card,
			};
		}

		let migrateIconStyle = AppBar.styles.migrateIcon;
		if (this.props.searchQuery !== undefined) {
			migrateIconStyle = {
				...migrateIconStyle,
				...theme.text.dark.primary,
			};
		}

		return (
			<div style={style}>
				<span style={titleStyle}>Ellie's Pad</span>

				<div style={AppBar.styles.spacer} />

				<SearchField
					initialQuery={this.props.searchQuery}
					onQueryChange={this.onSearchQueryChange}
					onFocus={this.onSearchFocus}
					onBlur={this.onSearchBlur}
					style={searchFieldStyle}
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
					<Icon style={migrateIconStyle} name="update" />
				</button>
			</div>
		);
	}
}

export default Relay.createContainer(AppBar, {
	fragments: {},
});
