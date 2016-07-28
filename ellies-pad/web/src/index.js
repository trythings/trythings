import gapi from 'gapi';
import React from 'react';
import ReactDOM from 'react-dom';
import { IndexRoute, Route, browserHistory } from 'react-router';
import { RelayRouter } from 'react-router-relay';
import Relay from 'react-relay';

import App from './App.js';
import Migrate from './Migrate.js';
import SignedInApp from './SignedInApp.js';
import SignIn from './SignIn.js';

Relay.injectNetworkLayer(new Relay.DefaultNetworkLayer('/graphql', {
	headers: {
		get Authorization() {
			if (gapi.auth2.getAuthInstance().isSignedIn.get()) {
				return `Bearer ${
					gapi.auth2.getAuthInstance().currentUser.get().getAuthResponse().id_token
				}`;
			}
			return null;
		},
	},
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
		<Route path="/" component={App}>
			<IndexRoute
				component={SignedInApp}
				onEnter={SignedInApp.onEnter}
				queries={queries}
			/>
			<Route
				path="search/(:query)"
				component={SignedInApp}
				onEnter={SignedInApp.onEnter}
				queries={queries}
			/>
			<Route
				path="signin"
				component={SignIn}
				onEnter={SignIn.onEnter}
			/>
			<Route
				path="migrate"
				component={Migrate}
				queries={queries}
			/>
		</Route>
	</RelayRouter>
);

gapi.load('auth2', () => {
	gapi.auth2.init({
		client_id: '695504958192-8k3tf807271m7jcllcvlauddeqhbr0hg.apps.googleusercontent.com',
	}).then(() => {
		ReactDOM.render(element, document.getElementById('App'));
	});
});
