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

// Placeholder component for Approval Editor
const ApprovalEditorPage = () => {
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
						Редактор согласования
					</Typography>
					<Typography variant='body1' color='text.secondary'>
						Управление и настройка процессов согласования документов
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
							label='Конструктор шаблонов согласования'
							icon={<AccountTreeIcon />}
							iconPosition='start'
						/>
						<Tab
							label='Назначение шаблонов согласования'
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
								Конструктор шаблонов согласования
							</Typography>
							<Divider sx={{ my: 2 }} />
							<Typography variant='body1' color='text.secondary'>
								Этот раздел находится в разработке. Здесь будет отображаться
								интерфейс для управления шаблонами согласования, создания и
								редактирования маршрутов согласования документов,
								последовательности этапов и назначения ответственных лиц.
							</Typography>
						</Box>
					)}
					{activeTab === 1 && (
						<Box>
							<Typography variant='h6' gutterBottom>
								Назначение шаблонов согласования
							</Typography>
							<Divider sx={{ my: 2 }} />
							<Typography variant='body1' color='text.secondary'>
								Этот раздел находится в разработке. Здесь будет отображаться
								интерфейс для назначения шаблонов согласования.
							</Typography>
						</Box>
					)}
				</Box>
			</Paper>
		</Box>
	);
};

export default ApprovalEditorPage;
