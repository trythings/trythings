import _pick from 'lodash/pick';
import React from 'react';

import resetStyles from './resetStyles.js';
import theme from './theme.js';

export default class Card extends React.Component {
	static propTypes = {
		autoFocus: React.PropTypes.bool,
		children: React.PropTypes.node,
		style: React.PropTypes.shape({
			flex: React.PropTypes.string,
		}),
	};

	static styles = {
		card: {
			...resetStyles,

			alignItems: 'stretch',
			borderRadius: 2,
			backgroundColor: theme.colors.card,
			flexDirection: 'column',
			overflow: 'visible',
		},
	};

	state = {
		hasFocus: false,
	};

	onBlur = (event) => {
		if (event.relatedTarget && !event.currentTarget.contains(event.relatedTarget)) {
			this.setState({ hasFocus: false });
		}
	};

	onFocus = () => {
		this.setState({ hasFocus: true });
	};

	ref = (node) => {
		if (node && this.props.autoFocus) {
			node.focus();
		}
	};

	render() {
		let elevation = theme.elevation[2];
		if (this.state.hasFocus) {
			elevation = {
				...theme.elevation[8],
				marginLeft: -8,
				marginRight: -8,
			};
		}

		return (
			<div
				onBlur={this.onBlur}
				onFocus={this.onFocus}
				tabIndex={-1}
				style={{
					...Card.styles.card,
					..._pick(this.props.style, ['flex']),
					...elevation,
				}}
				ref={this.ref}
			>
				{this.props.children}
			</div>
		);
	}
}
