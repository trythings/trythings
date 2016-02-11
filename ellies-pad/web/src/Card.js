import React from 'react';

import resetStyles from './resetStyles.js';
import theme from './theme.js';

export default class Card extends React.Component {
	static propTypes = {
		children: React.PropTypes.node,
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

	static styles = {
		card: {
			...resetStyles,

			borderRadius: 2,
			backgroundColor: theme.colors.card,
			flexDirection: 'column',
			overflow: 'visible',
		},
	};

	cardStyle() {
		if (this.state.hasFocus) {
			return {
				...Card.styles.card,
				...theme.elevation[8],
				marginLeft: -8,
				marginRight: -8,
			};
		}
		return {
			...Card.styles.card,
			...theme.elevation[2],
		};
	}

	render() {
		return (
			<div
				onBlur={this.onBlur}
				onFocus={this.onFocus}
				tabIndex={-1}
				style={this.cardStyle()}
			>
				{this.props.children}
			</div>
		);
	}
}
