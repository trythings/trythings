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
			backgroundColor: theme.colors.dividers.dark,
		},
		showAll: {
			...resetStyles,
			alignSelf: 'stretch',
			cursor: 'pointer',
			justifyContent: 'center',
			paddingBottom: 8,
			paddingTop: 8,
		},
		showAllText: {
			...resetStyles,
			...theme.text.dark.secondary,
			// alignSelf: 'center',
			fontSize: 14,
		},
		empty: {
			...resetStyles,
			...theme.text.dark.secondary,

			alignSelf: 'center',
			fontSize: 14,
		},
	};

	state = {
		isShowingAll: false,
	};

	onFocus = (event) => {
		event.stopPropagation();
	};

	onShowAllClick = () => {
		this.setState({ isShowingAll: true });
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

		const defaultNumShowing = 10;

		let tasks = this.props.tasks;
		if (!this.state.isShowingAll) {
			tasks = tasks.slice(0, defaultNumShowing);
		}

		const isShowAllVisible = !this.state.isShowingAll &&
			this.props.tasks.length > defaultNumShowing;

		return (
			<Card>
				<ol style={TaskList.styles.list} onFocus={this.onFocus} tabIndex={-1}>
					{tasks.map((task, i, array) => (
						<li key={task.id} style={TaskList.styles.listItem}>
							<TaskListItem task={task} />
							{i < array.length - 1
								? <hr style={TaskList.styles.divider} />
								: null
							}
						</li>
					))}

					{
						isShowAllVisible ?
							<hr style={TaskList.styles.divider} /> :
							null
					}
					{
						isShowAllVisible ?
							(
								<li style={TaskList.styles.showAll} onClick={this.onShowAllClick}>
									<span style={TaskList.styles.showAllText}>
										{`Show ${this.props.tasks.length - defaultNumShowing} remaining tasks`}
									</span>
								</li>
							) :
							null
					}
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
