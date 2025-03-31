import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
import { alpha, Box, Button, Paper, Typography, useTheme } from '@mui/material';
import { AxiosError } from 'axios';
import React from 'react';

interface ErrorStateProps {
	error?: Error | AxiosError | unknown;
	height?: number | string;
	width?: number | string;
	title?: string;
	message?: string;
	onRetry?: () => void;
}

const ErrorState: React.FC<ErrorStateProps> = ({
	error,
	height = 300,
	width = '100%',
	title = 'Ошибка загрузки данных',
	message,
	onRetry,
}) => {
	const theme = useTheme();

	// Extract error message from various error types
	const getErrorMessage = (): string => {
		if (message) return message;
		
		if (error instanceof AxiosError) {
			return error.response?.data?.message || error.message || 'Ошибка сервера';
		}
		
		if (error instanceof Error) {
			return error.message;
		}
		
		return 'Неизвестная ошибка';
	};

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
				backgroundColor: alpha(theme.palette.error.light, 0.05),
				borderLeft: `4px solid ${theme.palette.error.main}`,
				p: 3,
			}}
		>
			<ErrorOutlineIcon
				color='error'
				sx={{
					fontSize: 40,
					mb: 2,
					opacity: 0.8,
				}}
			/>
			
			<Typography
				color='error'
				variant='subtitle1'
				fontWeight={600}
				align='center'
			>
				{title}
			</Typography>
			
			<Typography
				color='text.secondary'
				variant='body2'
				align='center'
				sx={{ mt: 1, mb: onRetry ? 2 : 0 }}
			>
				{getErrorMessage()}
			</Typography>
			
			{onRetry && (
				<Box sx={{ mt: 2 }}>
					<Button
						variant='outlined'
						color='primary'
						onClick={onRetry}
						sx={{
							borderRadius: 2,
							textTransform: 'none',
							px: 3,
						}}
					>
						Повторить
					</Button>
				</Box>
			)}
		</Paper>
	);
};

export default ErrorState;