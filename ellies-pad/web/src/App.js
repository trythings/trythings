import React from 'react';
import Relay from 'react-relay';

const App = (props) => {
	console.log(props);
	return <p>hello, world!</p>;
}

export default Relay.createContainer(App, {
	fragments: {
		viewer: () => Relay.QL`
			fragment on User {
				id,
				tasks {
					id,
					title,
				},
			}
		`,
	},
});
