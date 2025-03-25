import ArchiveIcon from '@mui/icons-material/Archive';
import DescriptionIcon from '@mui/icons-material/Description';
import FolderIcon from '@mui/icons-material/Folder';
import { Box, CircularProgress, styled, Typography } from '@mui/material';
import { RichTreeView, TreeItem2, TreeItem2Props } from '@mui/x-tree-view';
import { useMutation, useQuery } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import React, { useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';
import axiosFetching from '../api/AxiosFetch';
import config from '../constants/Configurations.json';
import { Directory, TreeDataItem } from '../interfaces/FilesTree';
import { openContextMenu } from '../store/Slices/contexMenuSlice';
import ContextMenu from './ContextMenu';

const getFolders = config.getFiles;
const createFile = config.createFile;

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

const HighlightedTreeItem = styled(CustomTreeItem)(({ theme }) => ({
	'& .MuiTreeItem-content': {
		backgroundColor: '#f0f0f0', // Серый фон
		borderRadius: theme.shape.borderRadius,
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

const FilesTree: React.FC<{ isArchive: boolean }> = ({ isArchive }) => {
	const dispatch = useDispatch();
	const [highlightedItemId, setHighlightedItemId] = useState<string | null>(
		null
	);

	const {
		data: apiResponse,
		isLoading,
		isError,
		error,
		refetch: refreshTree,
	} = useQuery({
		queryKey: ['directories', isArchive],
		queryFn: async () => {
			const response = await axiosFetching.post(getFolders, {
				is_archive: isArchive,
			});
			return response.data;
		},
	});

	const createFileMutation = useMutation({
		mutationFn: async (data: { directory_id: number; name: string }) => {
			const response = await axiosFetching.post(createFile, data);
			return response.data;
		},
		onSuccess: () => {
			refreshTree();
		},
		onError: (error: any) => {
			console.error('Error creating file:', error);
		},
	});

	useEffect(() => {
		refreshTree();
	}, []);

	const handleContextMenu = (
		event: React.MouseEvent<HTMLDivElement>,
		itemId: string,
		itemType: 'directory' | 'file'
	) => {
		event.preventDefault();
		dispatch(
			openContextMenu({
				mouseX: event.clientX - 2,
				mouseY: event.clientY - 4,
				itemId,
				itemType,
				treeType: isArchive ? 'archive' : 'work',
			})
		);
	};

	const handleDrop = (
		event: React.DragEvent<HTMLDivElement>,
		directoryId: number
	) => {
		event.preventDefault();

		const files = Array.from(event.dataTransfer.files);

		setHighlightedItemId(null);

		files.forEach(file => {
			const fileName = file.name;

			const existingFile = findFileInDirectory(
				treeItems,
				directoryId,
				fileName
			);

			if (existingFile) {
				if (confirm(`Файл "${fileName}" уже существует. Заменить?`)) {
					createFileMutation.mutate({
						directory_id: directoryId,
						name: fileName,
					});
				}
			} else {
				createFileMutation.mutate({
					directory_id: directoryId,
					name: fileName,
				});
			}
		});
	};

	const findFileInDirectory = (
		items: TreeDataItem[],
		directoryId: number,
		fileName: string
	): boolean => {
		const directory = items.find(item => item.id === `dir-${directoryId}`);
		if (!directory || !directory.children) return false;

		// Удаляем состояние (draft) и нормализуем имя нового файла
		const cleanFileName = fileName
			.replace(/\s*\(draft\)\s*$/, '')
			.trim()
			.toLowerCase();

		return directory.children.some(child => {
			// Удаляем состояние (draft) и нормализуем имя существующего файла
			const childName = child.label
				.replace(/\s*\(draft\)\s*$/, '')
				.trim()
				.toLowerCase();
			return childName === cleanFileName;
		});
	};

	if (isLoading) {
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

	return (
		<Box
			sx={{
				width: '100%',
				maxWidth: 400,
				bgcolor: 'background.paper',
				border: '1px solid #ccc',
				borderRadius: 1,
				padding: 2,
			}}
		>
			<Typography
				variant='h6'
				align='center'
				sx={{
					mb: 2,
					fontWeight: 'bold',
				}}
			>
				{isArchive ? 'Архивное дерево' : 'Рабочее дерево'}
			</Typography>

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

						const isHighlighted = highlightedItemId === itemId;

						const TreeComponent = isHighlighted
							? HighlightedTreeItem
							: CustomTreeItem;

						return (
							<TreeComponent
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
										onDragOver={event => event.preventDefault()}
										onDrop={event => {
											const directoryId = parseInt(
												itemId!.replace('dir-', ''),
												10
											);
											handleDrop(event, directoryId);
										}}
										onDragEnter={event => {
											event.preventDefault();
											if (itemData.type === 'directory') {
												setHighlightedItemId(itemId!);
											}
										}}
										onDragLeave={() => {
											setHighlightedItemId(null);
										}}
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
										<span>{label}</span>
									</Box>
								}
							/>
						);
					},
				}}
				sx={{
					width: '100%',
				}}
			/>

			<ContextMenu />
		</Box>
	);
};

export default FilesTree;
