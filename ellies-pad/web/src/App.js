import React from 'react';
import Relay from 'react-relay';

const Task = Relay.createContainer((props) => {
	return <li>{props.task.title}</li>;
}, {
	fragments: {
		task: () => Relay.QL`
			fragment on Task {
				title,
			}
		`,
	},
});


const App = (props) => {
	return (
		<ol>
			{props.viewer.tasks.map(task => <Task key={task.id} task={task}/>)}
		</ol>
	);
}

export default Relay.createContainer(App, {
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				tasks {
					id,
					${Task.getFragment('task')}
				},
			}
		`,
	},
});
