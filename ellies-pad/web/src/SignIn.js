import _uniqueId from 'lodash/uniqueId';
import gapi from 'gapi';
import React from 'react';

export default class SignIn extends React.Component {
	static contextTypes = {
		router: React.PropTypes.shape({
			push: React.PropTypes.func.isRequired,
		}).isRequired,
	};

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

	componentDidMount() {
		this._isMounted = true;
		gapi.auth2.getAuthInstance().isSignedIn.listen(this.onSignIn);
	}

	componentWillUnmount() {
		this._isMounted = false;
	}

	onSignIn = (isSignedIn) => {
		if (this._isMounted && isSignedIn) {
			this.context.router.push('/');
		}
	};

	ref = (div) => {
		if (div) {
			gapi.signin2.render(div.dataset.id);
		}
	};

	render() {
		return <div id={this.id} data-id={this.id} ref={this.ref}></div>;
	}
}
