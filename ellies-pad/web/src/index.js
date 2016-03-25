import React from 'react';
import ReactDOM from 'react-dom';
import { Route, browserHistory } from 'react-router';
import { RelayRouter } from 'react-router-relay';
import Relay from 'react-relay';

import App from './App.js';
import TaskSearchResults from './TaskSearchResults.js';

Relay.injectNetworkLayer(new Relay.DefaultNetworkLayer('/graphql', {
	credentials: 'same-origin',
}));

const queries = {
	viewer: () => Relay.QL`
		query {
			viewer,
		}
	`,
};

const element = (
	<RelayRouter history={browserHistory}>
		<Route
			path="/"
			component={App}
			queries={queries}
		>
			<Route
				path="search/"
				component={TaskSearchResults}
				queries={queries}
			/>
			<Route
				path="search/:query"
				component={TaskSearchResults}
				queries={queries}
			/>
		</Route>
	</RelayRouter>
);

ReactDOM.render(element, document.getElementById('App'));
