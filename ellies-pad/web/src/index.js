import gapi from 'gapi';
import React from 'react';
import ReactDOM from 'react-dom';
import { IndexRoute, Route, browserHistory } from 'react-router';
import { RelayRouter } from 'react-router-relay';
import Relay from 'react-relay';

import App from './App.js';
import AppContainer from './AppContainer.js';
import SignInModal from './SignInModal.js';

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
		<Route path="/" component={AppContainer} xcxc="xcxc">
			<IndexRoute
				component={App}
				onEnter={App.onEnter}
				queries={queries}
			/>
			<Route
				path="signin"
				component={SignInModal}
				onEnter={SignInModal.onEnter}
			/>
			<Route
				path="search/(:query)"
				component={App}
				onEnter={App.onEnter}
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
