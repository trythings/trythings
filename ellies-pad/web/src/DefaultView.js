import React from 'react';

import resetStyles from './resetStyles.js';
import TaskSearch from './TaskSearch.js';

export default class DefaultView extends React.Component {
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
				/>
				<div style={DefaultView.styles.spacer} />

				<TaskSearch
					name="Incoming"
					query="NOT #now AND NOT #next AND NOT #later AND IsArchived: false"
				/>
				<div style={DefaultView.styles.spacer} />

				<TaskSearch
					name="#next"
					query="#next AND NOT #now AND IsArchived: false"
				/>
				<div style={DefaultView.styles.spacer} />

				<TaskSearch
					name="#later"
					query="#later AND NOT #next AND NOT #now AND IsArchived: false"
				/>
				<div style={DefaultView.styles.spacer} />

				<TaskSearch
					name="Archived"
					query="IsArchived: true"
				/>
			</div>
		);
	}
}
