import AccountTreeIcon from '@mui/icons-material/AccountTree';
import ArticleOutlinedIcon from '@mui/icons-material/ArticleOutlined';
import {
	Box,
	Divider,
	Paper,
	Tab,
	Tabs,
	Typography,
	alpha,
	useTheme,
} from '@mui/material';
import { useState } from 'react';

// Placeholder component for UsersPermissionsPage
const UsersPermissionsPage = () => {
	const theme = useTheme();
	const [activeTab, setActiveTab] = useState(0);

	const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
		setActiveTab(newValue);
	};

	return (
		<Box sx={{ p: 4 }}>
			<Paper
				elevation={2}
				sx={{
					borderRadius: 3,
					overflow: 'hidden',
					maxWidth: 1200,
					mx: 'auto',
					bgcolor: theme.palette.background.paper,
					boxShadow: `0 6px 20px ${alpha(theme.palette.primary.main, 0.12)}`,
				}}
			>
				{/* Header section */}
				<Box
					sx={{
						bgcolor: alpha(theme.palette.secondary.light, 0.1),
						py: 2.5,
						px: 4,
						borderBottom: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
					}}
				>
					<Typography variant='h5' fontWeight={600} gutterBottom>
						Управление правами и пользователями
					</Typography>
					<Typography variant='body1' color='text.secondary'>
						Эта страница находится в разработке. Здесь будет отображаться список
						пользователей и функции для управления их учетными записями и
						правами.
					</Typography>
				</Box>

				{/* Tabs navigation */}
				<Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
					<Tabs
						value={activeTab}
						onChange={handleTabChange}
						aria-label='approval editor tabs'
						sx={{
							'& .MuiTab-root': {
								textTransform: 'none',
								fontSize: '0.95rem',
								fontWeight: 500,
								py: 2,
								px: 3,
							},
						}}
					>
						<Tab
							label='Назначение ролей пользователям'
							icon={<AccountTreeIcon />}
							iconPosition='start'
						/>
						<Tab
							label='Назначение пользователей на проекты'
							icon={<ArticleOutlinedIcon />}
							iconPosition='start'
						/>
					</Tabs>
				</Box>

				{/* Tab panels */}
				<Box sx={{ p: 4 }}>
					{activeTab === 0 && (
						<Box>
							<Typography variant='h6' gutterBottom>
								Назначение ролей пользователям
							</Typography>
							<Divider sx={{ my: 2 }} />
							<Typography variant='body1' color='text.secondary'>
								Этот раздел находится в разработке. Здесь будет отображаться
								интерфейс назначения ролей пользователям.
							</Typography>
						</Box>
					)}
					{activeTab === 1 && (
						<Box>
							<Typography variant='h6' gutterBottom>
								Назначение пользователей на проекты
							</Typography>
							<Divider sx={{ my: 2 }} />
							<Typography variant='body1' color='text.secondary'>
								Этот раздел находится в разработке. Здесь будет отображаться
								интерфейс для назначения пользователей на проекты.
							</Typography>
						</Box>
					)}
				</Box>
			</Paper>
		</Box>
	);
};

export default UsersPermissionsPage;
