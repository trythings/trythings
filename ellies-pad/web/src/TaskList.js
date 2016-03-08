import React from 'react';
import Relay from 'react-relay';

import Card from './Card.js';
import resetStyles from './resetStyles.js';
import TaskListItem from './TaskListItem.js';
import theme from './theme.js';

class TaskList extends React.Component {
	static propTypes = {
		tasks: React.PropTypes.arrayOf(React.PropTypes.shape({
			// ...TaskListItem.propTypes.task
		})).isRequired,
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

	onFocus = (event) => {
		event.stopPropagation();
	};

	renderEmpty() {
		return (
			<span style={TaskList.styles.empty}>No results</span>
		);
	}

	render() {
		if (!this.props.tasks.length) {
			return this.renderEmpty();
		}

		return (
			<Card>
				<ol style={TaskList.styles.list} onFocus={this.onFocus} tabIndex={-1}>
					{this.props.tasks.map((task, i, array) => (
						<li key={task.id} style={TaskList.styles.listItem}>
							<TaskListItem task={task} />
							{i < array.length - 1
								? <hr style={TaskList.styles.divider} />
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
	fragments: {
		tasks: () => Relay.QL`
			fragment on Task @relay(plural: true) {
				id,
				${TaskListItem.getFragment('task')},
			}
		`,
	},
});
