import React from 'react';
import Relay from 'react-relay';

import Card from './Card.js';
import resetStyles from './resetStyles.js';
import TaskListItem from './TaskListItem.js';
import theme from './theme.js';

class TaskList extends React.Component {
	static propTypes = {
		search: React.PropTypes.shape({
			tasks: React.PropTypes.shape({
				edges: React.PropTypes.arrayOf(React.PropTypes.shape({
					node: React.PropTypes.shape({
						// ...TaskListItem.propTypes.task
					}).isRequired,
				})).isRequired,
				pageInfo: React.PropTypes.shape({
					hasNextPage: React.PropTypes.bool,
				}),
			}).isRequired,
		}).isRequired,
		relay: React.PropTypes.shape({
			setVariables: React.PropTypes.func.isRequired,
		}).isRequired,
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
		this.props.relay.setVariables({
			numTasksToShow: 100, // TODO#xcxc: Use the last count that we had of all of the tasks.
		}, (readyState) => {
			if (readyState.ready) {
				this.setState({ isShowingAll: true });
			}

			if (readyState.error) {
				console.log('encountered a sad error: ', readyState.error);
			}
		});
	};

	renderEmpty() {
		return (
			<span style={TaskList.styles.empty}>No results</span>
		);
	}

	render() {
		const tasks = this.props.search.tasks.edges.map((edge) => edge.node);

		if (!tasks.length) {
			return this.renderEmpty();
		}

		const defaultNumShowing = 10;

		let tasks = this.props.tasks;
		if (!this.state.isShowingAll) {
			tasks = tasks.slice(0, defaultNumShowing);
		}

		const isShowAllVisible = !this.state.isShowingAll &&
			this.props.search.tasks.pageInfo.hasNextPage;

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
	initialVariables: {
		numTasksToShow: 10,
	},
	fragments: {
		search: () => Relay.QL`
			fragment on Search {
				tasks(first: $numTasksToShow) {
					edges {
						node {
							id,
							${TaskListItem.getFragment('task')},
						},
					},
					pageInfo {
						hasNextPage,
					},
				}
			}
		`,
	},
});
