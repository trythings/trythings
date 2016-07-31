import React from 'react';
import Relay from 'react-relay';

import Card from './Card.js';
import resetStyles from './resetStyles.js';
import TaskListItem from './TaskListItem.js';
import theme from './theme.js';

class TaskList extends React.Component {
	static propTypes = {
		search: React.PropTypes.shape({
			numResults: React.PropTypes.number.isRequired,
			results: React.PropTypes.shape({
				edges: React.PropTypes.arrayOf(React.PropTypes.shape({
					node: React.PropTypes.shape({
						// ...TaskListItem.propTypes.task
					}).isRequired,
				})).isRequired,
				pageInfo: React.PropTypes.shape({
					hasNextPage: React.PropTypes.bool.isRequired,
				}),
			}).isRequired,
		}).isRequired,
		relay: React.PropTypes.shape({
			setVariables: React.PropTypes.func.isRequired,
			variables: React.PropTypes.shape({
				numTasksToShow: React.PropTypes.number.isRequired,
			}).isRequired,
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
			// Since we can't specify that we want all of the results,
			// we fetch 100 more than we know about.
			numTasksToShow: this.props.search.numResults + 100,
		}, (readyState) => {
			if (readyState.ready) {
				this.setState({ isShowingAll: true });
			}
			// TODO#Errors: Figure out how we want to handle client errors (readyState.errors != null).
		});
	};

	renderEmpty() {
		return (
			<span style={TaskList.styles.empty}>No results</span>
		);
	}

	render() {
		const tasks = this.props.search.results.edges.map((edge) => edge.node);

		if (!tasks.length) {
			return this.renderEmpty();
		}

		const isShowAllVisible = !this.state.isShowingAll &&
			this.props.search.results.pageInfo.hasNextPage;

		let numRemainingTasks =
			this.props.search.numResults - this.props.relay.variables.numTasksToShow;
		numRemainingTasks = numRemainingTasks < 0 ? 0 : numRemainingTasks;

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
										{`Show ${numRemainingTasks} remaining tasks`}
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
				numResults,
				results(first: $numTasksToShow) {
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
