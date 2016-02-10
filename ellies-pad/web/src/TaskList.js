import React from 'react';
import Relay from 'react-relay';

import Card from './Card.js';
import Task from './Task.js';
import theme from './theme.js';

class TaskList extends React.Component {
	static propTypes = {
		viewer: React.PropTypes.shape({
			tasks: React.PropTypes.arrayOf(React.PropTypes.object).isRequired, // ...Task.propTypes.task
		}).isRequired,
	};

	static styles = {
		list: {
			paddingLeft: 0,
			paddingRight: 0,
			paddingTop: 4,
			paddingBottom: 4,
			margin: 0,
			listStyle: 'none',
		},
		divider: {
			margin: 0,
			border: 'none',
			height: 1,
			backgroundColor: theme.text.dark.dividers,
		},
	};

	render() {
		return (
			<Card>
				<ol style={TaskList.styles.list}>
					{this.props.viewer.tasks.map((task, i, array) => {
						const tile = <Task key={task.id} task={task}/>;
						if (i === array.length - 1) {
							return tile;
						}
						return [
							tile,
							<hr style={TaskList.styles.divider}/>,
						];
					})}
				</ol>
			</Card>
		);
	}
}

export default Relay.createContainer(TaskList, {
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				tasks {
					id,
					${Task.getFragment('task')},
				},
			}
		`,
	},
});
