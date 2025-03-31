import { alpha, CircularProgress, Paper, Typography, useTheme } from '@mui/material';
import React from 'react';

interface LoadingStateProps {
	height?: number | string;
	width?: number | string;
	message?: string;
	spinnerSize?: number;
}

const LoadingState: React.FC<LoadingStateProps> = ({
	height = 300,
	width = '100%',
	message = '',
	spinnerSize = 36,
}) => {
	const theme = useTheme();

	return (
		<Paper
			elevation={2}
			sx={{
				display: 'flex',
				flexDirection: 'column',
				justifyContent: 'center',
				alignItems: 'center',
				height,
				width,
				borderRadius: 3,
				backgroundColor: theme.palette.background.paper,
				boxShadow: `0 6px 20px ${alpha(theme.palette.primary.main, 0.12)}`,
				gap: 2,
			}}
		>
			<CircularProgress
				size={spinnerSize}
				color='primary'
				sx={{
					opacity: 0.8,
				}}
			/>
			{message && (
				<Typography
					variant='body2'
					color='text.secondary'
					sx={{ mt: 1 }}
					align='center'
				>
					{message}
				</Typography>
			)}
		</Paper>
	);
};

export default LoadingState;