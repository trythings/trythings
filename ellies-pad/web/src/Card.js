import React from 'react';

import resetStyles from './resetStyles.js';
import theme from './theme.js';

export default class Card extends React.Component {
	static propTypes = {
		autoFocus: React.PropTypes.bool,
		children: React.PropTypes.node,
		flex: React.PropTypes.string,
	};

	state = {
		hasFocus: false,
	};

	onBlur = () => {
		this.setState({ hasFocus: false });
	};

	onFocus = () => {
		this.setState({ hasFocus: true });
	};

	ref = (node) => {
		if (node && this.props.autoFocus) {
			node.focus();
		}
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

	cardStyle() {
		let style = Card.styles.card;
		if (this.state.hasFocus) {
			style = {
				...style,
				...theme.elevation[8],
				marginLeft: -8,
				marginRight: -8,
			};
		} else {
			style = {
				...style,
				...theme.elevation[2],
			};
		}

		if (this.props.flex) {
			style = {
				...style,
				flex: this.props.flex,
			};
		}

		return style;
	}

	render() {
		return (
			<div
				onBlur={this.onBlur}
				onFocus={this.onFocus}
				tabIndex={-1}
				style={this.cardStyle()}
				ref={this.ref}
			>
				{this.props.children}
			</div>
		);
	}
}
