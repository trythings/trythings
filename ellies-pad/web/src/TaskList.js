import React from 'react';
import Relay from 'react-relay';

import Card from './Card.js';
import Task from './Task.js';

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

	};

	render() {
		return (
			<Card>
				<ol style={TaskList.styles.list}>
					{this.props.viewer.tasks.map(task => <Task key={task.id} task={task}/>)}
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
