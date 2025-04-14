import { CheckCircleOutline } from '@mui/icons-material';
import CloseIcon from '@mui/icons-material/Close';
import CreateNewFolderOutlinedIcon from '@mui/icons-material/CreateNewFolderOutlined';
import DeleteOutlineOutlinedIcon from '@mui/icons-material/DeleteOutlineOutlined';
import NoteAddOutlinedIcon from '@mui/icons-material/NoteAddOutlined';
import {
	Alert,
	alpha,
	Box,
	Dialog,
	DialogActions,
	DialogContent,
	DialogTitle,
	IconButton,
	Menu,
	MenuItem,
	Snackbar,
	TextField,
	Typography,
	useTheme,
} from '@mui/material';
import Button from '@mui/material/Button';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import axios, { AxiosError } from 'axios';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {axiosFetching, axiosFetchingFiles} from '../api/AxiosFetch';
import { updateApprovalsCount } from '../api/NavigationService';
import config from '../constants/Configurations.json';
import { closeContextMenu } from '../store/Slices/contextMenuSlice';
import { setPendingCount } from '../store/Slices/pendingApprovalsSlice';
import { RootState } from '../store/store';

const createFolder = config.createDirectory;
const deleteFolder = config.deleteDirectory;
const deleteFile = config.deleteFile;
const sendFileToApprove = config.sendFileToApprove;
const getApprovals = config.getApprovals;

type CreateFolderPayload = {
	parent_path_id: number;
	name: string;
};

type DeleteFolderPayload = {
	directory_id: number;
};

type DeleteFilePayload = {
	file_id: number;
};

const ContextMenu = () => {
	const theme = useTheme();
	const dispatch = useDispatch();
	const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);
	const [isDialogOpen, setIsDialogOpen] = useState(false);
	const [newName, setNewName] = useState('');
	const [snackbarOpen, setSnackbarOpen] = useState(false);
	const [snackbarMessage, setSnackbarMessage] = useState('');
	const [snackbarSeverity, setSnackbarSeverity] = useState<'success' | 'error'>(
		'success'
	);
	const treeRef = useRef<HTMLDivElement>(null);
	const { mouseX, mouseY, itemId, itemType, treeType } = useSelector(
		(state: RootState) => state.contextMenu
	);
	const queryClient = useQueryClient();

	// Handle closing the context menu
	const handleCloseMenu = useCallback(() => {
		dispatch(closeContextMenu());
	}, [dispatch]);

	// Обработчик нажатия Escape
	useEffect(() => {
		const handleEscape = (e: KeyboardEvent) => {
			if (e.key === 'Escape') {
				handleCloseMenu();
			}
		};

		document.addEventListener('keydown', handleEscape);

		return () => {
			document.removeEventListener('keydown', handleEscape);
		};
	}, [handleCloseMenu]); // Added handleCloseMenu to dependency array

	// Обработчик клика вне меню
	useEffect(() => {
		const handleClickOutside = (event: MouseEvent) => {
			// Проверяем, что меню открыто
			if (mouseX !== null && mouseY !== null) {
				// Находим элемент меню
				const menuElement = document.querySelector('.MuiMenu-paper');
				// Проверяем, что клик был не по меню
				if (menuElement && !menuElement.contains(event.target as Node)) {
					// Закрываем меню только если клик не по меню
					handleCloseMenu();
				}
			}
		};

		// Используем capture phase для гарантии перехвата события
		document.addEventListener('mousedown', handleClickOutside, true);

		return () => {
			document.removeEventListener('mousedown', handleClickOutside, true);
		};
	}, [mouseX, mouseY, handleCloseMenu]); // Added handleCloseMenu to dependency array

	const refreshCorrectTree = () => {
		if (treeType === 'work') {
			queryClient.invalidateQueries({ queryKey: ['directories', false] });
		} else if (treeType === 'archive') {
			queryClient.invalidateQueries({ queryKey: ['directories', true] });
		}
	};

	// Add this function to update approvals count without navigating to approvals page
	const updateApprovalsCountLocal = async () => {
		try {
			console.log('Updating approvals count');
			const response = await axiosFetching.get(getApprovals);
			if (response.data && Array.isArray(response.data)) {
				const pendingCount = response.data.length;
				// Use your Redux action to update the pending count
				dispatch(setPendingCount(pendingCount));
				console.log(`Updated pending count: ${pendingCount}`);
			}
		} catch (error) {
			console.error('Error updating approvals count:', error);
		}
	};

	const createFolderQuery = useMutation({
		mutationFn: async (data: CreateFolderPayload) => {
			const response = await axiosFetchingFiles.post(createFolder, data);
			return response.data;
		},
		onSuccess: () => {
			console.log('Folder created');
			refreshCorrectTree();
			setIsDialogOpen(false);
			setNewName('');
			setSnackbarMessage('Папка успешно создана');
			setSnackbarSeverity('success');
			setSnackbarOpen(true);
		},
		onError: (error: AxiosError) => {
			console.error('Error: ', error);
			setSnackbarMessage('Ошибка при создании папки');
			setSnackbarSeverity('error');
			setSnackbarOpen(true);
		},
	});

	const approveFileMutation = useMutation({
		mutationFn: async (fileId: number) => {
			// Log the request for debugging
			console.log(`Sending file ${fileId} for approval`);

			try {
				const url = sendFileToApprove.replace(':file_id', String(fileId));
				const response = await axiosFetching.put(url);
				return response.data;
			} catch (error) {
				// Enhanced error logging
				if (axios.isAxiosError(error)) {
					console.error(
						'API Error:',
						error.response?.status,
						error.response?.data
					);
					throw new Error(
						`Error ${error.response?.status}: ${JSON.stringify(
							error.response?.data
						)}`
					);
				}
				throw error;
			}
		},
		onSuccess: fileId => {
			// Add success notification
			setSnackbarMessage(`Файл успешно отправлен на согласование!`);
			setSnackbarSeverity('success');
			setSnackbarOpen(true);

			// Update both trees to reflect changes
			refreshCorrectTree();

			// Update approvals count globally
			updateApprovalsCount();
			updateApprovalsCountLocal();

			console.log(`File ${fileId} successfully sent for approval`);
		},
		onError: (error: Error) => {
			// User-friendly error message
			setSnackbarMessage(`Ошибка при отправке файла: ${error.message}`);
			setSnackbarSeverity('error');
			setSnackbarOpen(true);
			console.error('Error sending file for approval:', error);
		},
	});

	const deleteFolderMutation = useMutation({
		mutationFn: async (data: DeleteFolderPayload) => {
			const response = await axiosFetchingFiles.delete(deleteFolder, { data });
			return response.data;
		},
		onSuccess: () => {
			refreshCorrectTree();
			setSnackbarMessage('Папка успешно удалена');
			setSnackbarSeverity('success');
			setSnackbarOpen(true);
		},
		onError: (error: AxiosError) => {
			console.error('Error deleting folder:', error);
			setSnackbarMessage('Ошибка при удалении папки');
			setSnackbarSeverity('error');
			setSnackbarOpen(true);
		},
	});

	const deleteFileMutation = useMutation({
		mutationFn: async (data: DeleteFilePayload) => {
			const response = await axiosFetchingFiles.delete(deleteFile, { data });
			return response.data;
		},
		onSuccess: () => {
			refreshCorrectTree();
			setSnackbarMessage('Файл успешно удален');
			setSnackbarSeverity('success');
			setSnackbarOpen(true);
		},
		onError: (error: AxiosError) => {
			console.error('Error deleting folder:', error);
			setSnackbarMessage('Ошибка при удалении файла');
			setSnackbarSeverity('error');
			setSnackbarOpen(true);
		},
	});

	const handleCreateFolderSubmit = () => {
		if (!itemId || !newName.trim()) return;
		const parentPathId = parseInt(itemId.replace('dir-', ''), 10);
		createFolderQuery.mutate({ parent_path_id: parentPathId, name: newName });
	};

	const handleCreateFolder = () => {
		setIsDialogOpen(true);
		handleCloseMenu();
	};

	const handleSendForApproval = () => {
		setIsConfirmDialogOpen(true);
		handleCloseMenu();
	};

	const handleConfirmApproval = () => {
		if (!itemId) {
			console.error('No item selected for approval');
			setSnackbarMessage('Ошибка: не выбран файл для отправки');
			setSnackbarSeverity('error');
			setSnackbarOpen(true);
			setIsConfirmDialogOpen(false);
			return;
		}

		// Make sure we're dealing with a file
		if (!itemId.startsWith('file-')) {
			console.error('Selected item is not a file:', itemId);
			setSnackbarMessage('Ошибка: выбранный элемент не является файлом');
			setSnackbarSeverity('error');
			setSnackbarOpen(true);
			setIsConfirmDialogOpen(false);
			return;
		}

		try {
			const fileId = parseInt(itemId.replace('file-', ''), 10);
			if (isNaN(fileId)) {
				console.error('Invalid file ID:', itemId);
				setSnackbarMessage('Ошибка: неверный идентификатор файла');
				setSnackbarSeverity('error');
				setSnackbarOpen(true);
				setIsConfirmDialogOpen(false);
				return;
			}

			console.log('Preparing to send file for approval, ID:', fileId);
			approveFileMutation.mutate(fileId);
		} catch (error) {
			console.error('Error parsing file ID:', error);
			setSnackbarMessage('Ошибка при обработке идентификатора файла');
			setSnackbarSeverity('error');
			setSnackbarOpen(true);
		} finally {
			setIsConfirmDialogOpen(false);
		}
	};

	const handleCreateFile = () => {
		console.log('Create file for item:', itemId);
		handleCloseMenu();
	};

	const handleDeleteFolder = () => {
		if (!itemId) return;
		const directoryId = parseInt(itemId.replace('dir-', ''), 10);
		deleteFolderMutation.mutate({ directory_id: directoryId });
		handleCloseMenu();
	};

	const handleDeleteFile = () => {
		if (!itemId) return;
		const fileId = parseInt(itemId.replace('file-', ''), 10);
		deleteFileMutation.mutate({ file_id: fileId });
		handleCloseMenu();
	};

	// Проверка, является ли элемент архивным
	const isArchiveItem = treeType === 'archive';

	// Create menu items based on item type and archive status
	const menuItems = [];

	// Add create folder and file options for non-archive directories
	if (itemType === 'directory' && !isArchiveItem) {
		menuItems.push(
			<MenuItem
				key='create-folder'
				onClick={handleCreateFolder}
				sx={{
					py: 1,
					px: 2,
					'&:hover': {
						backgroundColor: alpha(theme.palette.primary.main, 0.08),
					},
				}}
			>
				<CreateNewFolderOutlinedIcon
					fontSize='small'
					sx={{
						mr: 1.5,
						color: theme.palette.primary.main,
					}}
				/>
				<Typography variant='body2'>Создать папку</Typography>
			</MenuItem>
		);

		menuItems.push(
			<MenuItem
				key='create-file'
				onClick={handleCreateFile}
				sx={{
					py: 1,
					px: 2,
					'&:hover': {
						backgroundColor: alpha(theme.palette.primary.main, 0.08),
					},
				}}
			>
				<NoteAddOutlinedIcon
					fontSize='small'
					sx={{
						mr: 1.5,
						color: theme.palette.primary.main,
					}}
				/>
				<Typography variant='body2'>Создать файл</Typography>
			</MenuItem>
		);
	}

	// Add delete option based on item type
	if (itemType === 'directory') {
		menuItems.push(
			<MenuItem
				key='delete-folder'
				onClick={handleDeleteFolder}
				sx={{
					py: 1,
					px: 2,
					'&:hover': {
						backgroundColor: alpha(theme.palette.error.main, 0.08),
					},
				}}
			>
				<DeleteOutlineOutlinedIcon
					fontSize='small'
					sx={{
						mr: 1.5,
						color: theme.palette.error.main,
					}}
				/>
				<Typography variant='body2' color={theme.palette.error.main}>
					Удалить папку
				</Typography>
			</MenuItem>
		);
	} else if (itemType === 'file') {
		menuItems.push(
			<MenuItem
				key='delete-file'
				onClick={handleDeleteFile}
				sx={{
					py: 1,
					px: 2,
					'&:hover': {
						backgroundColor: alpha(theme.palette.error.main, 0.08),
					},
				}}
			>
				<DeleteOutlineOutlinedIcon
					fontSize='small'
					sx={{
						mr: 1.5,
						color: theme.palette.error.main,
					}}
				/>
				<Typography variant='body2' color={theme.palette.error.main}>
					Удалить файл
				</Typography>
			</MenuItem>
		);
	}

	if (itemType === 'file') {
		menuItems.push(
			<MenuItem
				key='send-for-approval'
				onClick={handleSendForApproval}
				sx={{
					py: 1,
					px: 2,
					'&:hover': {
						backgroundColor: alpha(theme.palette.primary.main, 0.08),
					},
				}}
			>
				<CheckCircleOutline
					fontSize='small'
					sx={{
						mr: 1.5,
						color: theme.palette.primary.main,
					}}
				/>
				<Typography variant='body2'>Отправить на согласование</Typography>
			</MenuItem>
		);
	}

	// Fix for PaperProps deprecation - using slotProps instead
	return (
		<>
			<Menu
				open={mouseX !== null && mouseY !== null}
				onContextMenu={e => {
					e.preventDefault();
					e.stopPropagation();
				}}
				onClose={handleCloseMenu}
				anchorReference='anchorPosition'
				anchorPosition={
					mouseY !== null && mouseX !== null
						? { top: mouseY, left: mouseX }
						: undefined
				}
				disableRestoreFocus={true}
				disableAutoFocusItem={true}
				className='context-menu-component'
				autoFocus={false}
				slotProps={{
					paper: {
						elevation: 3,
						sx: {
							mt: 0.5,
							borderRadius: 2,
							minWidth: 180,
							padding: '4px 0',
							boxShadow: `0 4px 20px ${alpha(
								theme.palette.primary.main,
								0.15
							)}`,
							overflow: 'hidden',
						},
					},
				}}
			>
				{/* If no menu items, render an empty box to avoid empty Menu warning */}
				{menuItems.length > 0 ? menuItems : <Box sx={{ display: 'none' }} />}
			</Menu>

			<Dialog
				open={isDialogOpen}
				onClose={() => {
					setIsDialogOpen(false);
					treeRef.current?.focus();
				}}
				slotProps={{
					paper: {
						sx: {
							borderRadius: 3,
							boxShadow: `0 8px 32px ${alpha(
								theme.palette.primary.main,
								0.15
							)}`,
							maxWidth: 450,
							width: '100%',
						},
					},
				}}
			>
				<DialogTitle
					sx={{
						display: 'flex',
						alignItems: 'center',
						justifyContent: 'space-between',
						borderBottom: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
						pb: 2,
					}}
				>
					<Box display='flex' alignItems='center' gap={1}>
						<CreateNewFolderOutlinedIcon color='primary' />
						<Typography variant='h6' fontWeight={600}>
							Создание новой папки
						</Typography>
					</Box>
					<IconButton
						onClick={() => setIsDialogOpen(false)}
						size='small'
						sx={{ color: theme.palette.text.secondary }}
					>
						<CloseIcon fontSize='small' />
					</IconButton>
				</DialogTitle>
				<DialogContent sx={{ pt: 3, pb: 2 }}>
					<TextField
						autoFocus
						margin='dense'
						label='Имя папки'
						fullWidth
						value={newName}
						onChange={e => setNewName(e.target.value)}
						variant='outlined'
						sx={{
							'& .MuiOutlinedInput-root': {
								borderRadius: 2,
							},
						}}
					/>
				</DialogContent>
				<DialogActions
					sx={{
						px: 3,
						py: 2,
						borderTop: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
					}}
				>
					<Button
						onClick={() => setIsDialogOpen(false)}
						variant='outlined'
						sx={{
							borderRadius: 2,
							px: 3,
							textTransform: 'none',
						}}
					>
						Отмена
					</Button>
					<Button
						onClick={handleCreateFolderSubmit}
						disabled={!newName.trim()}
						variant='contained'
						sx={{
							borderRadius: 2,
							px: 3,
							textTransform: 'none',
							ml: 1,
							boxShadow: `0 4px 12px ${alpha(theme.palette.primary.main, 0.3)}`,
						}}
					>
						Создать
					</Button>
				</DialogActions>
			</Dialog>

			<Dialog
				open={isConfirmDialogOpen}
				onClose={() => setIsConfirmDialogOpen(false)}
				slotProps={{
					paper: {
						sx: {
							borderRadius: 3,
							boxShadow: `0 8px 32px ${alpha(
								theme.palette.primary.main,
								0.15
							)}`,
							maxWidth: 450,
							width: '100%',
						},
					},
				}}
			>
				<DialogTitle
					sx={{
						display: 'flex',
						alignItems: 'center',
						justifyContent: 'space-between',
						borderBottom: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
						pb: 2,
					}}
				>
					<Box display='flex' alignItems='center' gap={1}>
						<CheckCircleOutline color='primary' />
						<Typography variant='h6' fontWeight={600}>
							Подтверждение отправки
						</Typography>
					</Box>
					<IconButton
						onClick={() => setIsConfirmDialogOpen(false)}
						size='small'
						sx={{ color: theme.palette.text.secondary }}
					>
						<CloseIcon fontSize='small' />
					</IconButton>
				</DialogTitle>
				<DialogContent sx={{ pt: 3, pb: 2 }}>
					<Typography variant='body1'>
						Вы уверены, что хотите отправить файл на согласование?
					</Typography>
				</DialogContent>
				<DialogActions
					sx={{
						px: 3,
						py: 2,
						borderTop: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
					}}
				>
					<Button
						onClick={() => setIsConfirmDialogOpen(false)}
						variant='outlined'
						sx={{
							borderRadius: 2,
							px: 3,
							textTransform: 'none',
						}}
					>
						Отмена
					</Button>
					<Button
						onClick={handleConfirmApproval}
						variant='contained'
						sx={{
							borderRadius: 2,
							px: 3,
							textTransform: 'none',
							ml: 1,
							boxShadow: `0 4px 12px ${alpha(theme.palette.primary.main, 0.3)}`,
						}}
					>
						Принять
					</Button>
				</DialogActions>
			</Dialog>

			{/* Snackbar for notifications */}
			<Snackbar
				open={snackbarOpen}
				autoHideDuration={4000}
				onClose={() => setSnackbarOpen(false)}
				anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
			>
				<Alert
					onClose={() => setSnackbarOpen(false)}
					severity={snackbarSeverity}
					variant='filled'
					sx={{
						width: '100%',
						borderRadius: 2,
					}}
				>
					{snackbarMessage}
				</Alert>
			</Snackbar>
		</>
	);
};

export default ContextMenu;
