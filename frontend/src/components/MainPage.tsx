import { Box, Paper, Typography, alpha, useTheme } from '@mui/material';
import { useState } from 'react';
import axios from 'axios';
import ContextMenu from './ContextMenu';
import FilesTree from './FilesTree';
import { Canvas } from '@react-three/fiber';
import { OrbitControls, useGLTF } from '@react-three/drei';
import { CheckCircleOutline } from '@mui/icons-material';

const MainPage = () => {
    const theme = useTheme();
	const [, setIsDownloading] = useState(false);
    const [selectedItem, setSelectedItem] = useState<{
        id: string | null;
        type: 'file' | 'directory' | null;
        name: string | null;
    }>({
        id: null,
        type: null,
        name: null,
    });
    const [modelUrl, setModelUrl] = useState<string | null>(null);

	const handleItemSelect = async (
		id: string,
		type: 'file' | 'directory',
		name: string
	) => {
		if (type === 'file' && name.endsWith('.glb')) {
			try {
				const fileId = id.replace('file-', '');
				const response = await axios.get(`/files/${fileId}/download-direct`, {
					responseType: 'blob', // Важно указать, чтобы получить файл в виде Blob
				});
	
				// Логируем тип ответа
				console.log('Response:', response);
	
				// Проверяем, является ли ответ HTML-страницей
				const text = await new Response(response.data).text();
				if (text.startsWith('<!DOCTYPE html>')) {
					console.error('Server returned HTML instead of GLB file');
					return;
				}
	
				const fileUrl = URL.createObjectURL(response.data);
				setModelUrl(fileUrl);
				setSelectedItem({ id, type, name });
			} catch (error) {
				console.error('Error downloading GLB file:', error);
			}
		} else {
			setSelectedItem({ id, type, name });
			setModelUrl(null);
		}
	};

	const handleDownload = async (fileId: string) => {
		setIsDownloading(true);
		try {
			const response = await axios.get(`/files/${fileId}/download-direct`, {
				responseType: 'blob',
			});
			const url = window.URL.createObjectURL(new Blob([response.data]));
			const link = document.createElement('a');
			link.href = url;
			link.setAttribute('download', selectedItem?.name || 'file.glb');
			document.body.appendChild(link);
			link.click();
			link.remove();
		} catch (error) {
			console.error('Error downloading file:', error);
		} finally {
			setIsDownloading(false);
		}
	};

    const ModelViewer = ({ url }: { url: string }) => {
        const gltf = useGLTF(url);

        return (
            <Canvas
                camera={{ position: [0, 0, 5], fov: 50 }}
                style={{ width: '100%', height: '100%' }}
            >
                <ambientLight intensity={0.5} />
                <spotLight position={[10, 10, 10]} angle={0.15} penumbra={1} />
                <primitive object={gltf.scene} scale={0.5} />
                <OrbitControls />
            </Canvas>
        );
    };

    return (
        <Box className='main-page-container' sx={{ p: 2 }}>
            <Box
                sx={{
                    display: 'flex',
                    gap: 2,
                    height: 'calc(100vh - 120px)',
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
                        mt: '0px',
                    }}
                >
				<Box
					sx={{
						bgcolor: alpha(theme.palette.info.light, 0.1),
						borderBottom: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
						py: 2,
						px: 3,
						height: '57px',
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
					{selectedItem.type === 'file' && selectedItem.id && (
						<CheckCircleOutline
							fontSize='small'
							sx={{
								ml: 2,
								color: theme.palette.success.main,
								cursor: 'pointer',
								'&:hover': {
									opacity: 0.8,
								},
							}}
							onClick={() => {if (selectedItem.id) { handleDownload(selectedItem.id.replace('file-', ''))}}}
						/>
					)}
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
                        {modelUrl ? (
                            <ModelViewer url={modelUrl} />
                        ) : (
                            <Typography color='text.secondary' align='center'>
                                Выберите документ или 3D модель для предпросмотра
                            </Typography>
                        )}
                    </Box>
                </Paper>

                {/* Context menu component */}
                <ContextMenu />
            </Box>
        </Box>
    );
};

export default MainPage;