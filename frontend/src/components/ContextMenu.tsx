import {
	Button,
	Dialog,
	DialogActions,
	DialogContent,
	DialogTitle,
	Menu,
	MenuItem,
	TextField,
} from '@mui/material';
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
		const handleEscape = e => {
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
		const handleClickOutside = event => {
			// Проверяем, что меню открыто
			if (mouseX !== null && mouseY !== null) {
				// Находим элемент меню
				const menuElement = document.querySelector('.MuiMenu-paper');
				// Проверяем, что клик был не по меню
				if (menuElement && !menuElement.contains(event.target)) {
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

	const menuItems = [];

	if (itemType === 'directory') {
		menuItems.push(
			<MenuItem key='create-folder' onClick={handleCreateFolder}>
				Создать папку
			</MenuItem>,
			<MenuItem key='create-file' onClick={handleCreateFile}>
				Создать файл
			</MenuItem>,
			<MenuItem key='delete-folder' onClick={handleDeleteFolder}>
				Удалить папку
			</MenuItem>
		);
	}

	if (itemType === 'file') {
		menuItems.push(
			<MenuItem key='delete-file' onClick={handleDeleteFile}>
				Удалить файл
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
			>
				{menuItems.length > 0 ? menuItems : null}
			</Menu>

			<Dialog
				open={isDialogOpen}
				onClose={() => {
					setIsDialogOpen(false);
					treeRef.current?.focus();
				}}
			>
				<DialogTitle>Создание новой папки</DialogTitle>
				<DialogContent>
					<TextField
						autoFocus
						margin='dense'
						label='Имя папки'
						fullWidth
						value={newName}
						onChange={e => setNewName(e.target.value)}
					/>
				</DialogContent>
				<DialogActions>
					<Button onClick={() => setIsDialogOpen(false)}>Отмена</Button>
					<Button onClick={handleCreateFolderSubmit} disabled={!newName.trim()}>
						Создать
					</Button>
				</DialogActions>
			</Dialog>
		</>
	);
};

export default ContextMenu;
