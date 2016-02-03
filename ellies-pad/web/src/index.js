import React from 'react';
import ReactDOM from 'react-dom';
import Relay from 'react-relay';

import App from './App.js';

class AppRoute extends Relay.Route {
	static routeName = 'AppRoute';

	static queries = {
		viewer: () => Relay.QL`
			query {
				viewer
			}
		`,
	};
}

const element = (
	<Relay.RootContainer
		Component={App}
		route={new AppRoute()}
	/>
);

Relay.injectNetworkLayer(new Relay.DefaultNetworkLayer('/graphql', {
	credentials: 'same-origin',
}));

ReactDOM.render(element, document.getElementById('App'));
