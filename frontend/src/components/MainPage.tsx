import { Box, Paper, Typography, alpha, useTheme } from '@mui/material';
import { useState } from 'react';
import { axiosFetchingFiles } from '../api/AxiosFetch';
import ContextMenu from './ContextMenu';
import FilesTree from './FilesTree';

const MainPage = () => {
	const [modelUrl, setModelUrl] = useState<string | null>(null);
	const theme = useTheme();
	// State to track selected file or folder for preview
	const [selectedItem, setSelectedItem] = useState<{
		id: string | null;
		type: 'file' | 'directory' | null;
		name: string | null;
	}>({
		id: null,
		type: null,
		name: null,
	});

	// This function will be passed to FilesTree components to handle item selection
	const handleItemSelect = async (
		id: string,
		type: 'file' | 'directory',
		name: string
	) => {
		if (type === 'file' && name.endsWith('.glb')) {
			try {
				const fileId = id.replace('file-', '');
				const response = await axiosFetchingFiles.get(`/files/${fileId}/download-direct`, {
					responseType: 'blob', // Важно указать, чтобы получить файл в виде Blob
				});
	
				// Создаем временный URL для загруженного файла
				const fileUrl = URL.createObjectURL(response.data);
				setModelUrl(fileUrl); // Сохраняем URL модели в состоянии
				setSelectedItem({ id, type, name });
			} catch (error) {
				console.error('Error downloading GLB file:', error);
			}
		} else {
			setSelectedItem({ id, type, name });
			setModelUrl(null); // Очищаем модель, если выбран не `.glb` файл
		}
	};

	return (
		<Box className='main-page-container' sx={{ p: 2 }}>
			<Box
				sx={{
					display: 'flex',
					gap: 2,
					height: 'calc(100vh - 120px)', // Adjust based on header height and padding
				}}
			>
				{/* Left side - Stacked file trees */}
				<Box
					sx={{
						display: 'flex',
						flexDirection: 'column',
						gap: 2,
						width: '350px',
						maxWidth: '350px',
						flexShrink: 0,
					}}
				>
					<FilesTree isArchive={false} onItemSelect={handleItemSelect} />
					<FilesTree isArchive={true} onItemSelect={handleItemSelect} />
				</Box>

				{/* Right side - Preview area */}
				<Paper
					elevation={2}
					sx={{
						flex: 1,
						borderRadius: 3,
						overflow: 'hidden',
						display: 'flex',
						flexDirection: 'column',
						backgroundColor: theme.palette.background.paper,
						boxShadow: `0 6px 20px ${alpha(theme.palette.primary.main, 0.12)}`,
						mt: '0px', // Align with the top of the file trees
					}}
				>
					<Box
						sx={{
							bgcolor: alpha(theme.palette.info.light, 0.1),
							borderBottom: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
							py: 2,
							px: 3,
							height: '57px', // Match the header height of FilesTree
							display: 'flex',
							alignItems: 'center',
						}}
					>
						<Typography variant='h6' fontWeight={600}>
							{selectedItem.name
								? `Предпросмотр элемента: ${
										selectedItem.type === 'file' ? '' : 'Папка '
								  }${selectedItem.name}`
								: 'Предпросмотр документа'}
						</Typography>
					</Box>

					<Box
						sx={{
							flex: 1,
							display: 'flex',
							alignItems: 'center',
							justifyContent: 'center',
							p: 4,
						}}
					>
						<Typography color='text.secondary' align='center'>
							Выберите документ или 3D модель для предпросмотра
						</Typography>
					</Box>
				</Paper>

				{/* Context menu component */}
				<ContextMenu />
			</Box>
		</Box>
	);
};

export default MainPage;
