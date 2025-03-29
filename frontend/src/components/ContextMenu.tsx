import CloseIcon from '@mui/icons-material/Close';
import CreateNewFolderOutlinedIcon from '@mui/icons-material/CreateNewFolderOutlined';
import DeleteOutlineOutlinedIcon from '@mui/icons-material/DeleteOutlineOutlined';
import NoteAddOutlinedIcon from '@mui/icons-material/NoteAddOutlined';
import {
	alpha,
	Box,
	Dialog,
	DialogActions,
	DialogContent,
	DialogTitle,
	IconButton,
	Menu,
	MenuItem,
	TextField,
	Typography,
	useTheme,
} from '@mui/material';
import Button from '@mui/material/Button';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import { useEffect, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import axiosFetching from '../api/AxiosFetch';
import config from '../constants/Configurations.json';
import { closeContextMenu } from '../store/Slices/contextMenuSlice';
import { RootState } from '../store/store';

const createFolder = config.createDirectory;
const deleteFolder = config.deleteDirectory;
const deleteFile = config.deleteFile;

const ContextMenu = () => {
	const theme = useTheme();
	const dispatch = useDispatch();
	const [isDialogOpen, setIsDialogOpen] = useState(false);
	const [newName, setNewName] = useState('');
	const treeRef = useRef<HTMLDivElement>(null);
	const { mouseX, mouseY, itemId, itemType, treeType } = useSelector(
		(state: RootState) => state.contextMenu
	);
	const queryClient = useQueryClient();

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
	}, []);

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
	}, [mouseX, mouseY]);

	const refreshCorrectTree = () => {
		if (treeType === 'work') {
			queryClient.invalidateQueries({ queryKey: ['directories', false] });
		} else if (treeType === 'archive') {
			queryClient.invalidateQueries({ queryKey: ['directories', true] });
		}
	};

	const createFolderQuery = useMutation({
		mutationFn: async (data: { parent_path_id: number; name: string }) => {
			const response = await axiosFetching.post(createFolder, data);
			return response.data;
		},
		onSuccess: () => {
			console.log('Folder created');
			refreshCorrectTree();
			setIsDialogOpen(false);
			setNewName('');
		},
		onError: (error: AxiosError) => {
			console.error('Error: ', error);
		},
	});

	const deleteFolderMutation = useMutation({
		mutationFn: async (data: { directory_id: number }) => {
			const response = await axiosFetching.delete(deleteFolder, { data });
			return response.data;
		},
		onSuccess: () => {
			refreshCorrectTree();
		},
		onError: (error: any) => {
			console.error('Error deleting folder:', error);
		},
	});

	const deleteFileMutation = useMutation({
		mutationFn: async (data: { file_id: number }) => {
			const response = await axiosFetching.delete(deleteFile, { data });
			return response.data;
		},
		onSuccess: () => {
			refreshCorrectTree();
		},
		onError: (error: any) => {
			console.error('Error deleting folder:', error);
		},
	});

	const handleCreateFolderSubmit = () => {
		if (!itemId || !newName.trim()) return;
		const parentPathId = parseInt(itemId.replace('dir-', ''), 10);
		createFolderQuery.mutate({ parent_path_id: parentPathId, name: newName });
	};

	const handleCloseMenu = () => {
		dispatch(closeContextMenu());
	};

	const handleCreateFolder = () => {
		setIsDialogOpen(true);
		handleCloseMenu();
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
				PaperProps={{
					elevation: 3,
					sx: {
						mt: 0.5,
						borderRadius: 2,
						minWidth: 180,
						padding: '4px 0',
						boxShadow: `0 4px 20px ${alpha(theme.palette.primary.main, 0.15)}`,
						overflow: 'hidden',
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
				PaperProps={{
					sx: {
						borderRadius: 3,
						boxShadow: `0 8px 32px ${alpha(theme.palette.primary.main, 0.15)}`,
						maxWidth: 450,
						width: '100%',
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
		</>
	);
};

export default ContextMenu;
