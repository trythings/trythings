import _uniqueId from 'lodash/uniqueId';
import gapi from 'gapi';
import React from 'react';

export default class SignInModal extends React.Component {
	constructor(...args) {
		super(...args);
		this.id = _uniqueId('SignInModal');
		this.state = {
			auth2: gapi.auth2,
		};
	}

	ref = (div) => {
		console.log(gapi.auth2.getAuthInstance().isSignedIn.get()); // xcxc
		gapi.signin2.render(div.dataset.id);
	};

	render() {
		return <div id={this.id} data-id={this.id} ref={this.ref}></div>;
	}
}
