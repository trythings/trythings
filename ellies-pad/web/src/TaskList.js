import React from 'react';
import Relay from 'react-relay';

import Card from './Card.js';
import resetStyles from './resetStyles.js';
import TaskListItem from './TaskListItem.js';
import theme from './theme.js';

class TaskList extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.object.shape({
			space: React.PropTypes.object.shape({
				// ...TaskListItem.propTypes.task
				tasks: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
			}),
		}).isRequired,
	};

	onFocus = (event) => {
		event.stopPropagation();
	};

	static styles = {
		list: {
			...resetStyles,
			alignItems: 'stretch',
			flexDirection: 'column',
			paddingTop: 4,
			paddingBottom: 4,
			overflow: 'visible',
		},
		listItem: {
			...resetStyles,
			alignItems: 'stretch',
			flexDirection: 'column',
			overflow: 'visible',
		},
		divider: {
			...resetStyles,
			height: 1,
			backgroundColor: theme.text.dark.dividers.color,
		},
		empty: {
			...resetStyles,
			...theme.text.dark.secondary,

			alignSelf: 'center',
			fontSize: 14,
		},
	};

	renderEmpty() {
		return (
			<span style={TaskList.styles.empty}>No results</span>
		);
	}

	render() {
		if (!this.props.viewer.space.tasks.length) {
			return this.renderEmpty();
		}

		return (
			<Card>
				<ol style={TaskList.styles.list} onFocus={this.onFocus} tabIndex={-1}>
					{this.props.viewer.space.tasks.map((task, i, array) => (
						<li key={task.id} style={TaskList.styles.listItem}>
							<TaskListItem task={task}/>
							{i < array.length - 1
								? <hr style={TaskList.styles.divider}/>
								: null
							}
						</li>
					))}
				</ol>
			</Card>
		);
	}
}

export default Relay.createContainer(TaskList, {
	initialVariables: {
		query: null,
	},

	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				space {
					tasks(query: $query) {
						id,
						${TaskListItem.getFragment('task')},
					},
				},
			}
		`,
	},
});
