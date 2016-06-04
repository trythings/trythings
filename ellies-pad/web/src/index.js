import gapi from 'gapi';
import React from 'react';
import ReactDOM from 'react-dom';
import { Route, browserHistory } from 'react-router';
import { RelayRouter } from 'react-router-relay';
import Relay from 'react-relay';

import App from './App.js';

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

const prepareParams = () => ({});

const element = (
	<RelayRouter history={browserHistory}>
		<Route
			path="/(search/(:query))"
			component={App}
			prepareParams={prepareParams}
			queries={queries}
		/>
	</RelayRouter>
);

gapi.load('auth2', () => {
	gapi.auth2.init({
		client_id: '695504958192-8k3tf807271m7jcllcvlauddeqhbr0hg.apps.googleusercontent.com',
	});
	ReactDOM.render(element, document.getElementById('App'));
});
