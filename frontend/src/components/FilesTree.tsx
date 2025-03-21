import ArchiveIcon from '@mui/icons-material/Archive';
import DescriptionIcon from '@mui/icons-material/Description';
import FolderIcon from '@mui/icons-material/Folder';
import {
	Box,
	CircularProgress,
	Menu,
	MenuItem,
	styled,
	Typography,
} from '@mui/material';
import { RichTreeView, TreeItem2, TreeItem2Props } from '@mui/x-tree-view';
import { useMutation } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import React, { useEffect, useState } from 'react';
import axiosFetching from '../api/AxiosFetch';
import config from '../constants/Configurations.json';
import { Directory, TreeDataItem } from '../interfaces/FilesTree';

const getFolders = config.getFiles;
const createDirectory = config.createDirectory;
const deleteDirectory = config.deleteDirectory;
const createFile = config.createFile;
const deleteFile = config.deleteFile;

//У M-UI СВОЯ БИБЛИОТЕКА СТИЛЕЙ, В ЭТОМ КОМПОНЕНТЕ РЕШИЛ ИСПОЛЬЗОВАТЬ ЕЕ
const CustomTreeItem = styled(TreeItem2)(({ theme }) => ({
	// СТИЛИЗАЦИЯ КОНТЕЙНЕРА MuiTreeItem
	'& .MuiTreeItem-content': {
		padding: theme.spacing(0.5, 0),
	},
	'& .MuiTreeItem-label': {
		display: 'flex',
		alignItems: 'center',
		gap: theme.spacing(1),
	},
}));

// ПРЕОБРАЗОВАНИЕ ДАННЫХ В ФОРМАТ RichTreeView
const transformDataToTreeItems = (data: Directory[]): TreeDataItem[] => {
	const map = new Map<number, TreeDataItem>();

	//ИДЕМ В ДВА ЗАХОДА СНАЧАЛА ДОБАВЛЯЕМ ВСЕ ПАПКИ В MAP

	data.forEach(item => {
		map.set(item.id, {
			id: `dir-${item.id}`,
			label: item.name_folder,
			status: item.status,
			type: 'directory',
			children: [],
		});
	});

	//ВТОРАЯ ПРОХОДКА, ДОБАВЛЯЕМ ФАЙЛЫ. СОЗДАЕМ ИЕРАРХИЮ
	data.forEach(item => {
		const node = map.get(item.id)!;

		item.files.forEach(file => {
			node.children!.push({
				id: `file-${file.id}`,
				label: file.name_file,
				status: file.status,
				type: 'file',
			});
		});

		//ПРОВЕРКА НА РОДИТЕЛЬСКУЮ ПАПКУ. ЕСЛИ ЕСТЬ - ДОБАВЛЯЕМ В НЕЕ ДОЧЕРНИЕ ЭЛЕМЕНТЫ
		if (item.parent_path_id) {
			const parent = map.get(item.parent_path_id);
			if (parent) parent.children!.push(node);
		}
	});

	return data
		.filter(item => !item.parent_path_id)
		.map(item => map.get(item.id)!);
};

const FilesTree: React.FC = () => {
	const [contextMenu, setContextMenu] = useState<{
		mouseX: number;
		mouseY: number;
		itemId?: string;
		itemType?: 'directory' | 'file';
	} | null>(null);

	const {
		mutate,
		isPending,
		isError,
		error,
		data: apiResponse,
	} = useMutation({
		mutationFn: async () => {
			const response = await axiosFetching.post(getFolders, {
				is_archive: true,
			});
			return response.data;
		},
		onError: (error: AxiosError<{ message?: string }>) => {
			console.error(
				'Error fetching folders:',
				error.response?.data?.message || error.message
			);
		},
	});

	useEffect(() => {
		mutate();
	}, [mutate]);

	if (isPending) {
		return (
			<Box
				sx={{
					display: 'flex',
					justifyContent: 'center',
					alignItems: 'center',
					height: 200,
				}}
			>
				<CircularProgress />
			</Box>
		);
	}

	if (isError) {
		return (
			<Typography color='error' align='center'>
				Error loading data:{' '}
				{error instanceof AxiosError ? error.message : 'Unknown error'}
			</Typography>
		);
	}

	const treeItems = apiResponse
		? transformDataToTreeItems(apiResponse.data)
		: [];

	const handleContextMenu = (
		event: React.MouseEvent<HTMLDivElement>,
		itemId: string,
		itemType: 'directory' | 'file'
	) => {
		event.preventDefault();
		setContextMenu(
			contextMenu === null
				? {
						mouseX: event.clientX - 2,
						mouseY: event.clientY - 4,
						itemId,
						itemType,
				  }
				: null
		);
	};

	const handleCloseContextMenu = () => {
		setContextMenu(null);
	};

	const handleCreateFolder = () => {
		handleCloseContextMenu();
	};

	const handleDeleteFolder = () => {
		handleCloseContextMenu();
	};

	const handleCreateFile = () => {
		handleCloseContextMenu();
	};

	const handleDeleteFile = () => {
		handleCloseContextMenu();
	};

	return (
		<>
			<RichTreeView
				items={treeItems}
				defaultExpandedItems={['dir-1']}
				slots={{
					item: (props: TreeItem2Props) => {
						const { itemId, label, ...rest } = props;

						const findItem = (
							items: TreeDataItem[],
							id: string
						): TreeDataItem | undefined => {
							for (const item of items) {
								if (item.id === id) return item;
								if (item.children) {
									const found = findItem(item.children, id);
									if (found) return found;
								}
							}
							return undefined;
						};

						const itemData = findItem(treeItems, itemId!);

						if (!itemData) {
							return <TreeItem2 {...rest} itemId={itemId} label={label} />;
						}

						return (
							<CustomTreeItem
								{...rest}
								itemId={itemId}
								label={
									<Box
										display='flex'
										alignItems='center'
										gap={1}
										onContextMenu={event =>
											handleContextMenu(event, itemId!, itemData.type)
										}
									>
										{itemData.type === 'directory' ? (
											itemData.status === 'archive' ? (
												<ArchiveIcon color='error' />
											) : (
												<FolderIcon color='primary' />
											)
										) : (
											<DescriptionIcon color='secondary' />
										)}
										<span>
											{label} ({itemData.status})
										</span>
									</Box>
								}
							/>
						);
					},
				}}
				sx={{
					width: '100%',
					maxWidth: 400,
					bgcolor: 'background.paper',
					border: '1px solid #ccc',
					borderRadius: 1,
					padding: 1,
				}}
			/>

			<Menu
				open={contextMenu !== null}
				onClose={handleCloseContextMenu}
				anchorReference='anchorPosition'
				anchorPosition={
					contextMenu !== null
						? { top: contextMenu.mouseY, left: contextMenu.mouseX }
						: undefined
				}
			>
				{contextMenu?.itemType === 'directory' && (
					<>
						<MenuItem onClick={handleCreateFolder}>Создать папку</MenuItem>
						<MenuItem onClick={handleCreateFile}>Создать файл</MenuItem>
						<MenuItem onClick={handleDeleteFolder}>Удалить папку</MenuItem>
					</>
				)}
				{contextMenu?.itemType === 'file' && (
					<MenuItem onClick={handleDeleteFile}>Удалить файл</MenuItem>
				)}
			</Menu>
		</>
	);
};

export default FilesTree;
