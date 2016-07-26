import _uniqueId from 'lodash/uniqueId';
import gapi from 'gapi';
import React from 'react';

export default class SignIn extends React.Component {
	// Routing.
	static onEnter = (nextState, replace) => {
		if (gapi.auth2.getAuthInstance().isSignedIn.get()) {
			replace('/');
		}
	};

	constructor(...args) {
		super(...args);
		this.id = _uniqueId('SignIn');
	}

	ref = (div) => {
		gapi.signin2.render(div.dataset.id);
	};

	render() {
		return <div id={this.id} data-id={this.id} ref={this.ref}></div>;
	}
}
