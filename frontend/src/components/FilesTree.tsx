import ArchiveIcon from '@mui/icons-material/Archive';
import DescriptionOutlinedIcon from '@mui/icons-material/DescriptionOutlined';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import FolderOutlinedIcon from '@mui/icons-material/FolderOutlined';
import { alpha, Box, Paper, Typography, useTheme } from '@mui/material';
import { RichTreeView, TreeItem2, TreeItem2Props } from '@mui/x-tree-view';
import { useMutation, useQuery } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import React, {
	useCallback,
	useEffect,
	useMemo,
	useRef,
	useState,
} from 'react';
import { useDispatch } from 'react-redux';
import {axiosFetchingFiles} from '../api/AxiosFetch';
import config from '../constants/Configurations.json';
import { Directory, TreeDataItem } from '../interfaces/FilesTree';
import {
	closeContextMenu,
	openContextMenu,
} from '../store/Slices/contextMenuSlice';
import ErrorState from './ErrorState';
import LoadingState from './LoadingState';

const getFolders = config.getFiles;
const createFile = config.createFile;

interface FileUploadData {
	directory_id: number;
	name: string;
}

// Расширяем интерфейс TreeItem2Props
interface ExtendedTreeItem2Props extends TreeItem2Props {
	expansionState?: 'expanded' | 'collapsed' | 'loading';
}

// Расширяем интерфейс FilesTree добавляя onItemSelect
interface FilesTreeProps {
	isArchive: boolean;
	onItemSelect?: (id: string, type: 'file' | 'directory', name: string) => void;
}

const FilesTree: React.FC<FilesTreeProps> = ({ isArchive, onItemSelect }) => {
	const theme = useTheme();
	const dispatch = useDispatch();
	const treeViewRef = useRef<HTMLUListElement>(null);
	const [highlightedItemId, setHighlightedItemId] = useState<string | null>(
		null
	);
	const [selectedItemId, setSelectedItemId] = useState<string | null>(null);

	// Поиск элемента в дереве любой глубины вложенности
	const findTreeItem = useCallback(
		(items: TreeDataItem[], id: string): TreeDataItem | undefined => {
			for (const item of items) {
				if (item.id === id) {
					return item;
				}

				if (item.children && item.children.length > 0) {
					const foundInChildren = findTreeItem(item.children, id);
					if (foundInChildren) {
						return foundInChildren;
					}
				}
			}
			return undefined;
		},
		[]
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
			const response = await axiosFetchingFiles.post(getFolders, {
				is_archive: isArchive,
			});
			return response.data;
		},
	});

	// Преобразование данных API в формат дерева
	const transformDataToTreeItems = useCallback(
		(directories: Directory[]): TreeDataItem[] => {
			if (!directories || directories.length === 0) return [];

			// Создаем Map для хранения всех элементов дерева
			const treeItemsMap = new Map<number, TreeDataItem>();

			// Заполняем Map элементами
			directories.forEach(directory => {
				treeItemsMap.set(directory.id, {
					id: `dir-${directory.id}`,
					label: directory.name_folder,
					status: directory.status,
					type: 'directory',
					children: [],
				});
			});

			// Добавляем файлы и связываем директории с родителями
			directories.forEach(directory => {
				const directoryNode = treeItemsMap.get(directory.id)!;

				// Добавляем файлы как дочерние элементы
				directory.files.forEach(file => {
					directoryNode.children!.push({
						id: `file-${file.id}`,
						label: file.name_file,
						status: file.status,
						type: 'file',
					});
				});

				// Связываем с родителем, если он есть в нашей Map
				if (
					directory.parent_path_id &&
					treeItemsMap.has(directory.parent_path_id)
				) {
					const parentNode = treeItemsMap.get(directory.parent_path_id);
					if (parentNode) parentNode.children!.push(directoryNode);
				}
			});

			// Находим виртуальные корневые элементы
			const rootItems: TreeDataItem[] = [];

			directories.forEach(directory => {
				const directoryNode = treeItemsMap.get(directory.id)!;

				// Если у директории нет родителя или родитель недоступен пользователю,
				// считаем её корневой
				const hasAccessibleParent =
					directory.parent_path_id &&
					directories.some(
						parentDir => parentDir.id === directory.parent_path_id
					);

				if (!directory.parent_path_id || !hasAccessibleParent) {
					rootItems.push(directoryNode);
				}
			});

			return rootItems;
		},
		[]
	);

	// Мемоизируем дерево элементов
	const treeItems = useMemo(
		() => (apiResponse ? transformDataToTreeItems(apiResponse.data) : []),
		[apiResponse, transformDataToTreeItems]
	);

	// Проверка наличия файла в директории
	const findFileInDirectory = useCallback(
		(items: TreeDataItem[], directoryId: number, fileName: string): boolean => {
			const directory = findTreeItem(items, `dir-${directoryId}`);
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
		},
		[findTreeItem]
	);

	// Снятие фокуса с элементов дерева
	const clearTreeFocus = useCallback(() => {
		if (treeViewRef.current) {
			const focusedElements =
				treeViewRef.current.querySelectorAll('.Mui-focused');
			focusedElements.forEach(el => {
				el.classList.remove('Mui-focused');
			});
		}
	}, []);

	// Обработка выбора элемента
	const handleItemSelection = useCallback(
		(itemId: string, itemType: 'directory' | 'file', label: string) => {
			setSelectedItemId(itemId);
			clearTreeFocus();

			if (onItemSelect) {
				onItemSelect(itemId, itemType, label);
			}
		},
		[onItemSelect, clearTreeFocus]
	);

	// Обработка правого клика для контекстного меню
	const handleContextMenu = useCallback(
		(
			event: React.MouseEvent<Element>,
			itemId: string,
			itemType: 'directory' | 'file',
			label: string
		) => {
			event.preventDefault();
			event.stopPropagation();

			handleItemSelection(itemId, itemType, label);

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
		},
		[dispatch, isArchive, handleItemSelection]
	);

	// Обработка создания файла через mutation
	const createFileMutation = useMutation({
		mutationFn: async (data: FileUploadData) => {
			const response = await axiosFetchingFiles.post(createFile, data);
			return response.data;
		},
		onSuccess: () => {
			refreshTree();
		},
		onError: (error: AxiosError<{ message?: string }>) => {
			console.error('Error creating file:', error);
		},
	});

	// Обработка перетаскивания файлов
	const handleDrop = useCallback(async (event: React.DragEvent<HTMLDivElement>, directoryId: number) => {
		event.preventDefault();
		event.stopPropagation();
		const files = Array.from(event.dataTransfer.files);
		setHighlightedItemId(null);
	
		// Обрабатываем каждый файл
		files.forEach(async (file) => {
			// Проверяем, что файл имеет расширение .glb
			if (!file.name.endsWith('.glb')) {
				console.warn(`File "${file.name}" is not a .glb file and will be skipped.`);
				return;
			}
	
			// Создаем FormData для файла
			const formData = new FormData();
			formData.append('file', file); // Файл
			formData.append('directory_id', String(directoryId)); // ID директории
			formData.append('name', file.name); // Имя файла
	
			try {
				// Отправляем файл на сервер
				await axiosFetchingFiles.post('/files/upload', formData, {
					headers: {
						'Content-Type': 'multipart/form-data',
					},
				});
	
				// Обновляем дерево после успешной загрузки
				refreshTree();
			} catch (error) {
				console.error('Error uploading file:', error);
			}
		});
	}, [refreshTree]);

	// Глобальный обработчик кликов для снятия выделения
	useEffect(() => {
		const handleClickOutside = (event: MouseEvent) => {
			if (
				treeViewRef.current &&
				!treeViewRef.current.contains(event.target as Node)
			) {
				const contextMenu = document.querySelector('.MuiMenu-paper');
				if (!contextMenu || !contextMenu.contains(event.target as Node)) {
					setSelectedItemId(null);
				}
			}
		};

		document.addEventListener('mousedown', handleClickOutside);
		return () => {
			document.removeEventListener('mousedown', handleClickOutside);
		};
	}, []);

	// Загрузка дерева при монтировании компонента
	useEffect(() => {
		refreshTree();
	}, [refreshTree]);

	// Сброс выбора при изменении типа хранилища
	useEffect(() => {
		setSelectedItemId(null);
		clearTreeFocus();
	}, [isArchive, clearTreeFocus]);

	// Проверка наличия виртуальных корней
	const hasPartialTreeAccess = useMemo(() => {
		if (!apiResponse?.data) return false;

		return apiResponse.data.some(
			(directory: Directory) =>
				directory.parent_path_id &&
				!apiResponse.data.some(
					(parentDir: Directory) => parentDir.id === directory.parent_path_id
				)
		);
	}, [apiResponse]);

	// Отображение состояния загрузки
	if (isLoading) {
		return (
			<LoadingState
				message={`Загрузка ${
					isArchive ? 'архивного' : 'рабочего'
				} хранилища...`}
			/>
		);
	}

	// Отображение ошибки
	if (isError) {
		return <ErrorState error={error} onRetry={refreshTree} />;
	}

	// Основной рендер компонента
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
			{/* Заголовок хранилища */}
			<Box
			sx={{
				bgcolor: isArchive
					? alpha(theme.palette.primary.light, 0.1) // Синий фон для архивного дерева
					: alpha(theme.palette.primary.light, 0.1), // Оставляем как есть для рабочего дерева
				borderBottom: `1px solid ${alpha(
					isArchive
						? theme.palette.primary.main // Синяя граница для архивного дерева
						: theme.palette.primary.main,
					0.1
				)}`,
				py: 2,
				px: 3,
				display: 'flex',
				alignItems: 'center',
				gap: 1,
				height: '57px',
			}}
		>
			{isArchive ? (
				<ArchiveIcon
					sx={{
						color: theme.palette.primary.main, // Синий цвет иконки для архивного дерева
						fontSize: 22,
					}}
				/>
			) : (
				<FolderOpenIcon
					sx={{
						color: theme.palette.primary.main, // Оставляем синий цвет для рабочего дерева
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
						? theme.palette.primary.dark // Синий текст для архивного дерева
						: theme.palette.primary.dark, // Оставляем синий текст для рабочего дерева
					letterSpacing: '-0.3px',
				}}
			>
				{isArchive ? 'Архивное хранилище' : 'Рабочее хранилище'}
			</Typography>
		</Box>

			{/* Индикатор для случая, когда показывается только часть дерева */}
			{hasPartialTreeAccess && (
				<Box
					sx={{
						py: 1,
						px: 2,
						bgcolor: alpha(theme.palette.info.light, 0.1),
						borderBottom: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
					}}
				>
					<Typography
						variant='caption'
						color='text.secondary'
						sx={{
							display: 'block',
							fontSize: '0.75rem',
						}}
					>
						Отображаются только доступные вам папки и файлы
					</Typography>
				</Box>
			)}

			{/* Контейнер для дерева */}
			<Box sx={{ p: 2, maxHeight: 600, overflowY: 'auto' }}>
				{treeItems.length === 0 ? (
					<Typography color='text.secondary' sx={{ p: 2, textAlign: 'center' }}>
						Нет доступных файлов или папок
					</Typography>
				) : (
					<RichTreeView
						ref={treeViewRef}
						items={treeItems}
						defaultExpandedItems={treeItems
							.filter(item => item.type === 'directory')
							.map(item => item.id)}
						disableSelection={true}
						autoFocus={false}
						slots={{
							item: (props: ExtendedTreeItem2Props) => {
								const { itemId = '', label, expansionState, ...rest } = props;

								const itemData =
									treeItems.length > 0
										? findTreeItem(treeItems, itemId)
										: undefined;

								if (!itemData) {
									return (
										<TreeItem2 {...rest} itemId={itemId} label={label || ''} />
									);
								}

								const isHighlighted = highlightedItemId === itemId;
								const isSelected = selectedItemId === itemId;
								const isDirectory = itemData.type === 'directory';
								const isArchiveItem = itemData.status === 'archive';
								const isExpanded = expansionState === 'expanded';

								// Проверка является ли папка виртуальным корнем
								const isVirtualRoot = (() => {
									if (!itemData.type || itemData.type !== 'directory')
										return false;

									const directoryId = parseInt(itemId.replace('dir-', ''), 10);
									const directory = apiResponse?.data?.find(
										(dir: Directory) => dir.id === directoryId
									);

									return (
										directory?.parent_path_id &&
										!apiResponse?.data?.some(
											(parentDir: Directory) =>
												parentDir.id === directory.parent_path_id
										)
									);
								})();

								return (
									<TreeItem2
										{...rest}
										itemId={itemId}
										className={`tree-item-component ${
											isSelected ? 'custom-selected' : ''
										} ${isVirtualRoot ? 'virtual-root' : ''}`}
										tabIndex={-1}
										onContextMenu={event => {
											handleContextMenu(
												event,
												itemId,
												itemData.type,
												itemData.label
											);
											return false;
										}}
										onClick={e => {
											e.preventDefault();
											e.stopPropagation();
										}}
										label={
											<Box
												display='flex'
												alignItems='center'
												gap={1.5}
												onClick={e => {
													e.stopPropagation();
													handleItemSelection(
														itemId,
														itemData.type,
														itemData.label
													);
												}}
												sx={{
													py: 0.5,
													width: '100%',
													height: '100%',
													padding: '4px 8px',
													margin: '2px 0',
													borderRadius: 1.5,
													transition: 'all 0.2s ease',
													backgroundColor: isSelected
														? alpha(theme.palette.primary.main, 0.12)
														: isHighlighted
														? alpha(theme.palette.primary.main, 0.08)
														: 'transparent',
													'&:hover': {
														backgroundColor: alpha(
															theme.palette.primary.main,
															0.06
														),
													},
													// Стили для виртуальных корней
													...(isVirtualRoot && {
														borderLeft: `2px solid ${theme.palette.info.main}`,
														pl: 1.5,
													}),
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
																	itemId.replace('dir-', ''),
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
																setHighlightedItemId(itemId);
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
												{/* Иконка элемента */}
												{isDirectory ? (
													isArchiveItem ? (
														<ArchiveIcon
															sx={{
																color: theme.palette.warning.main,
																fontSize: 20,
															}}
														/>
													) : isExpanded ? (
														<FolderOpenIcon
															sx={{
																color: isVirtualRoot
																	? theme.palette.info.main
																	: theme.palette.primary.main,
																fontSize: 20,
															}}
														/>
													) : (
														<FolderOutlinedIcon
															sx={{
																color: isVirtualRoot
																	? theme.palette.info.main
																	: theme.palette.primary.main,
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
												{/* Название элемента */}
												<Typography
													variant='body2'
													sx={{
														fontWeight: isDirectory ? 500 : 400,
														color: isArchiveItem
															? theme.palette.warning.dark
															: isVirtualRoot
															? theme.palette.info.dark
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
								'&.Mui-focused': {
									outline: 'none !important',
									backgroundColor: 'transparent !important',
								},
								'&.Mui-selected': {
									backgroundColor: 'transparent !important',
									outline: 'none !important',
								},
								'&.custom-selected > .MuiTreeItem-content': {
									backgroundColor: 'transparent !important',
								},
								'&.virtual-root': {
									position: 'relative',
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
							'& .MuiTreeItem-content.Mui-focused': {
								backgroundColor: 'transparent !important',
							},
							'& .MuiTreeItem-content.Mui-selected': {
								backgroundColor: 'transparent !important',
							},
						}}
					/>
				)}
			</Box>
		</Paper>
	);
};

export default FilesTree;