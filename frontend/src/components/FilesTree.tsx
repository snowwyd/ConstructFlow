import ArchiveIcon from '@mui/icons-material/Archive';
import DescriptionOutlinedIcon from '@mui/icons-material/DescriptionOutlined';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import FolderOutlinedIcon from '@mui/icons-material/FolderOutlined';
import {
	alpha,
	Box,
	CircularProgress,
	Paper,
	Typography,
	useTheme,
} from '@mui/material';
import { RichTreeView, TreeItem2, TreeItem2Props } from '@mui/x-tree-view';
import { useMutation, useQuery } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import React, { useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';
import axiosFetching from '../api/AxiosFetch';
import config from '../constants/Configurations.json';
import { Directory, TreeDataItem } from '../interfaces/FilesTree';
import {
	closeContextMenu,
	openContextMenu,
} from '../store/Slices/contextMenuSlice';

const getFolders = config.getFiles;
const createFile = config.createFile;

const FilesTree: React.FC<{ isArchive: boolean }> = ({ isArchive }) => {
	const theme = useTheme();
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
		event.stopPropagation();

		dispatch(closeContextMenu());
		dispatch(
			openContextMenu({
				mouseX: event.clientX,
				mouseY: event.clientY,
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
		event.stopPropagation();

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

		const cleanFileName = fileName
			.replace(/\s*\(draft\)\s*$/, '')
			.trim()
			.toLowerCase();

		return directory.children.some(child => {
			const childName = child.label
				.replace(/\s*\(draft\)\s*$/, '')
				.trim()
				.toLowerCase();
			return childName === cleanFileName;
		});
	};

	// Function to transform API data to tree format
	const transformDataToTreeItems = (data: Directory[]): TreeDataItem[] => {
		const map = new Map<number, TreeDataItem>();

		data.forEach(item => {
			map.set(item.id, {
				id: `dir-${item.id}`,
				label: item.name_folder,
				status: item.status,
				type: 'directory',
				children: [],
			});
		});

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

			if (item.parent_path_id) {
				const parent = map.get(item.parent_path_id);
				if (parent) parent.children!.push(node);
			}
		});

		return data
			.filter(item => !item.parent_path_id)
			.map(item => map.get(item.id)!);
	};

	if (isLoading) {
		return (
			<Paper
				elevation={2}
				sx={{
					display: 'flex',
					justifyContent: 'center',
					alignItems: 'center',
					height: 300,
					width: '100%',
					borderRadius: 3,
					backgroundColor: theme.palette.background.paper,
					boxShadow: `0 6px 20px ${alpha(theme.palette.primary.main, 0.12)}`,
				}}
			>
				<CircularProgress
					size={36}
					color='primary'
					sx={{
						opacity: 0.8,
					}}
				/>
			</Paper>
		);
	}

	if (isError) {
		return (
			<Paper
				elevation={2}
				sx={{
					display: 'flex',
					flexDirection: 'column',
					justifyContent: 'center',
					alignItems: 'center',
					height: 300,
					width: '100%',
					borderRadius: 3,
					backgroundColor: alpha(theme.palette.error.light, 0.05),
					borderLeft: `4px solid ${theme.palette.error.main}`,
					p: 3,
				}}
			>
				<Typography
					color='error'
					variant='subtitle1'
					fontWeight={600}
					align='center'
				>
					Ошибка загрузки данных
				</Typography>
				<Typography
					color='text.secondary'
					variant='body2'
					align='center'
					sx={{ mt: 1 }}
				>
					{error instanceof AxiosError ? error.message : 'Неизвестная ошибка'}
				</Typography>
			</Paper>
		);
	}

	const treeItems = apiResponse
		? transformDataToTreeItems(apiResponse.data)
		: [];

	return (
		<Paper
			elevation={2}
			sx={{
				backgroundColor: theme.palette.background.paper,
				borderRadius: 3,
				overflow: 'hidden',
				width: '100%',
				maxWidth: 400,
				height: 'fit-content',
				boxShadow: `0 6px 20px ${alpha(theme.palette.primary.main, 0.12)}`,
				position: 'relative',
				zIndex: 1,
				transition: 'box-shadow 0.3s ease, transform 0.2s ease',
				'&:hover': {
					boxShadow: `0 8px 30px ${alpha(theme.palette.primary.main, 0.18)}`,
				},
			}}
		>
			<Box
				sx={{
					bgcolor: isArchive
						? alpha(theme.palette.warning.light, 0.1)
						: alpha(theme.palette.primary.light, 0.1),
					borderBottom: `1px solid ${alpha(
						isArchive ? theme.palette.warning.main : theme.palette.primary.main,
						0.1
					)}`,
					py: 2,
					px: 3,
					display: 'flex',
					alignItems: 'center',
					gap: 1,
				}}
			>
				{isArchive ? (
					<ArchiveIcon
						sx={{
							color: theme.palette.warning.main,
							fontSize: 22,
						}}
					/>
				) : (
					<FolderOpenIcon
						sx={{
							color: theme.palette.primary.main,
							fontSize: 22,
						}}
					/>
				)}
				<Typography
					variant='h6'
					fontWeight={600}
					sx={{
						fontSize: '1.1rem',
						color: isArchive
							? theme.palette.warning.dark
							: theme.palette.primary.dark,
						letterSpacing: '-0.3px',
					}}
				>
					{isArchive ? 'Архивное хранилище' : 'Рабочее хранилище'}
				</Typography>
			</Box>

			<Box sx={{ p: 2, maxHeight: 600, overflowY: 'auto' }}>
				<RichTreeView
					items={treeItems}
					defaultExpandedItems={['dir-1']} // Возвращаем оригинальное свойство вместо управляемого состояния
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
							const isDirectory = itemData.type === 'directory';
							const isArchiveItem = itemData.status === 'archive';

							return (
								<TreeItem2
									{...rest}
									itemId={itemId}
									className='tree-item-component'
									onContextMenu={event => {
										event.stopPropagation();
										handleContextMenu(event, itemId!, itemData.type);
										return false;
									}}
									sx={{
										'& .MuiTreeItem-content': {
											padding: '4px 8px',
											margin: '2px 0',
											borderRadius: 1.5,
											transition: 'all 0.2s ease',
											backgroundColor: isHighlighted
												? alpha(theme.palette.primary.main, 0.12)
												: 'transparent',
											'&:hover': {
												backgroundColor: alpha(
													theme.palette.primary.main,
													0.06
												),
											},
											'&.Mui-focused, &.Mui-selected': {
												backgroundColor: alpha(
													theme.palette.primary.main,
													0.12
												),
											},
										},
									}}
									label={
										<Box
											display='flex'
											alignItems='center'
											gap={1.5}
											sx={{
												py: 0.5,
												width: '100%',
												height: '100%',
											}}
											onDragOver={
												isDirectory
													? event => event.preventDefault()
													: undefined
											}
											onDrop={
												isDirectory
													? event => {
															const directoryId = parseInt(
																itemId!.replace('dir-', ''),
																10
															);
															handleDrop(event, directoryId);
													  }
													: undefined
											}
											onDragEnter={
												isDirectory
													? event => {
															event.preventDefault();
															event.stopPropagation();
															setHighlightedItemId(itemId!);
													  }
													: undefined
											}
											onDragLeave={
												isDirectory
													? () => {
															setHighlightedItemId(null);
													  }
													: undefined
											}
										>
											{isDirectory ? (
												isArchiveItem ? (
													<ArchiveIcon
														sx={{
															color: theme.palette.warning.main,
															fontSize: 20,
														}}
													/>
												) : props.expansionState === 'expanded' ? (
													<FolderOpenIcon
														sx={{
															color: theme.palette.primary.main,
															fontSize: 20,
														}}
													/>
												) : (
													<FolderOutlinedIcon
														sx={{
															color: theme.palette.primary.main,
															fontSize: 20,
														}}
													/>
												)
											) : (
												<DescriptionOutlinedIcon
													sx={{
														color: theme.palette.grey[600],
														fontSize: 20,
													}}
												/>
											)}
											<Typography
												variant='body2'
												sx={{
													fontWeight: isDirectory ? 500 : 400,
													color: isArchiveItem
														? theme.palette.warning.dark
														: theme.palette.text.primary,
													overflow: 'hidden',
													textOverflow: 'ellipsis',
													whiteSpace: 'nowrap',
													flex: 1,
												}}
											>
												{label}
											</Typography>
										</Box>
									}
								/>
							);
						},
					}}
					sx={{
						width: '100%',
						'& .MuiTreeItem-root': {
							position: 'relative',
							'&::before': {
								content: '""',
								position: 'absolute',
								left: '10px',
								top: '24px',
								bottom: 0,
								width: '1px',
								bgcolor: alpha(theme.palette.primary.main, 0.15),
								display: 'none',
							},
							'&:has(.MuiTreeItem-group)::before': {
								display: 'block',
							},
						},
						'& .MuiTreeItem-group': {
							marginLeft: '16px',
							paddingLeft: '12px',
							borderLeft: `1px dashed ${alpha(
								theme.palette.primary.main,
								0.15
							)}`,
						},
						'& .MuiTreeItem-iconContainer': {
							width: '20px',
							display: 'inline-flex',
							alignItems: 'center',
							justifyContent: 'center',
							marginRight: '4px',
							'& svg': {
								fontSize: '1rem',
								color: theme.palette.action.active,
							},
						},
					}}
				/>
			</Box>
		</Paper>
	);
};

export default FilesTree;
