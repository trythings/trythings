import React from 'react';
import ReactDOM from 'react-dom';
import Relay from 'react-relay';

import App from './App.js';

class AppRoute extends Relay.Route {
	static routeName = 'AppRoute';
	
	static queries = {
		tasks: () => Relay.QL`query { tasks }`,
	};
}

const element = (
	<Relay.RootContainer
  Component={App}
  route={new AppRoute()}
	/>
);

ReactDOM.render(element, document.getElementById('App'));
