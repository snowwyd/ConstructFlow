// TextViewer.tsx
import { Box, Paper, alpha, useTheme } from '@mui/material';
import React from 'react';

interface TextViewerProps {
	text: string;
}

const TextViewer: React.FC<TextViewerProps> = ({ text }) => {
	const theme = useTheme();

	return (
		<Paper
			elevation={1}
			sx={{
				width: '100%',
				height: '100%',
				overflow: 'auto',
				p: 3,
				bgcolor: alpha(theme.palette.background.default, 0.7),
				border: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
				borderRadius: 2,
			}}
		>
			<Box
				component='pre'
				sx={{
					margin: 0,
					fontFamily: 'monospace',
					fontSize: '0.9rem',
					whiteSpace: 'pre-wrap',
					wordBreak: 'break-word',
					color: theme.palette.text.primary,
					lineHeight: 1.5,
				}}
			>
				{text}
			</Box>
		</Paper>
	);
};

export default TextViewer;
