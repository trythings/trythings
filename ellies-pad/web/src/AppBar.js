import _pick from 'lodash/pick';
import React from 'react';

import resetStyles from './resetStyles.js';
import theme from './theme.js';

// TODO(annied): Add currently-working on tracker.
class AppBar extends React.Component {
	static propTypes = {
		children: React.PropTypes.node,
		style: React.PropTypes.shape({
			backgroundColor: React.PropTypes.string,
			color: React.PropTypes.string,
		}),
	};

	static styles = {
		appBar: {
			...resetStyles,
			...theme.elevation[4],

			alignItems: 'center',
			backgroundColor: theme.colors.primary.default,
			height: 56,
			justifyContent: 'space-between',
			minHeight: 56,
			paddingLeft: 16,
			paddingRight: 16,
		},
		title: {
			...resetStyles,
			...theme.text.light.primary,

			fontSize: 20,
		},
		children: {
			...resetStyles,
			flex: '1 0 auto',
			paddingLeft: 24,
		},
	};

	render() {
		return (
			<div
				style={{
					...AppBar.styles.appBar,
					..._pick(this.props.style, ['backgroundColor']),
				}}
			>
				<span
					style={{
						...AppBar.styles.title,
						..._pick(this.props.style, ['color']),
					}}
				>
					Ellie's Pad
				</span>

				<div style={AppBar.styles.children}>
					{this.props.children}
				</div>
			</div>
		);
	}
}

export default AppBar;
