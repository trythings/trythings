import React from 'react';
import Relay from 'react-relay';

import resetStyles from './resetStyles.js';
import TaskSearch from './TaskSearch.js';

export default class DefaultView extends React.Component {
	static propTypes = {
		space: React.PropTypes.shape({
			// ...TaskSearch.propTypes.space,
		}).isRequired,
	};

	static styles = {
		container: {
			...resetStyles,
			alignItems: 'stretch',
			flexDirection: 'column',
			overflow: 'visible',
		},
		spacer: {
			...resetStyles,
			padding: 12,
		},
	};

	render() {
		return (
			<div style={DefaultView.styles.container}>
				<TaskSearch
					name="#now"
					query="#now AND IsArchived: false"
					space={this.props.space}
				/>
				<div style={DefaultView.styles.spacer} />

				<TaskSearch
					name="Incoming"
					query="NOT #now AND NOT #next AND NOT #later AND IsArchived: false"
					space={this.props.space}
				/>
				<div style={DefaultView.styles.spacer} />

				<TaskSearch
					name="#next"
					query="#next AND NOT #now AND IsArchived: false"
					space={this.props.space}
				/>
				<div style={DefaultView.styles.spacer} />

				<TaskSearch
					name="#later"
					query="#later AND NOT #next AND NOT #now AND IsArchived: false"
					space={this.props.space}
				/>
				<div style={DefaultView.styles.spacer} />

				<TaskSearch
					name="Archived"
					query="IsArchived: true"
					space={this.props.space}
				/>
			</div>
		);
	}
}

export default Relay.createContainer(DefaultView, {
	fragments: {
		space: () => Relay.QL`
			fragment on Space {
				${TaskSearch.getFragment('space')},
			},
		`,
	},
});
