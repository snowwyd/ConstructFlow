import ArchiveIcon from '@mui/icons-material/Archive';
import DescriptionOutlinedIcon from '@mui/icons-material/DescriptionOutlined';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import FolderOutlinedIcon from '@mui/icons-material/FolderOutlined';
import {
	alpha,
	Box,
	Paper,
	Typography,
	useTheme,
} from '@mui/material';
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
import axiosFetching from '../api/AxiosFetch';
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

	// Рекурсивная функция поиска элемента в дереве любой глубины вложенности
	const findTreeItem = useCallback(
		(items: TreeDataItem[], id: string): TreeDataItem | undefined => {
			// Перебор всех элементов на текущем уровне
			for (const item of items) {
				// Проверка текущего элемента
				if (item.id === id) {
					return item;
				}

				// Рекурсивный поиск в дочерних элементах
				if (item.children && item.children.length > 0) {
					const foundInChildren = findTreeItem(item.children, id);
					if (foundInChildren) {
						return foundInChildren;
					}
				}
			}

			// Если ничего не найдено
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
			const response = await axiosFetching.post(getFolders, {
				is_archive: isArchive,
			});
			return response.data;
		},
	});

	// Function to transform API data to tree format
	const transformDataToTreeItems = useCallback(
		(data: Directory[]): TreeDataItem[] => {
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
		},
		[]
	);

	// Мемоизируем treeItems, чтобы не вычислять их при каждом рендере
	const treeItems = useMemo(
		() => (apiResponse ? transformDataToTreeItems(apiResponse.data) : []),
		[apiResponse, transformDataToTreeItems]
	);

	// Обработка поиска файла в директории (для Drag & Drop)
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

	// Функция для снятия фокуса с элементов дерева
	const clearTreeFocus = useCallback(() => {
		// Находим все элементы с классом .Mui-focused и убираем этот класс
		if (treeViewRef.current) {
			const focusedElements =
				treeViewRef.current.querySelectorAll('.Mui-focused');
			focusedElements.forEach(el => {
				el.classList.remove('Mui-focused');
			});
		}
	}, []);

	// Главная функция для обработки выбора элемента (универсальная)
	const handleItemSelection = useCallback(
		(itemId: string, itemType: 'directory' | 'file', label: string) => {
			// Устанавливаем выбранный элемент
			setSelectedItemId(itemId);

			// Снимаем фокус со всех элементов
			clearTreeFocus();

			// Уведомляем родительский компонент о выборе
			if (onItemSelect) {
				onItemSelect(itemId, itemType, label);
			}
		},
		[onItemSelect, clearTreeFocus]
	);

	// Обработчик правого клика (контекстное меню)
	const handleContextMenu = useCallback(
		(
			event: React.MouseEvent<Element>,
			itemId: string,
			itemType: 'directory' | 'file',
			label: string
		) => {
			event.preventDefault();
			event.stopPropagation();

			// Сначала выбираем элемент с помощью общей функции
			handleItemSelection(itemId, itemType, label);

			// Затем открываем контекстное меню
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
			const response = await axiosFetching.post(createFile, data);
			return response.data;
		},
		onSuccess: () => {
			refreshTree();
		},
		onError: (error: AxiosError<{ message?: string }>) => {
			console.error('Error creating file:', error);
		},
	});

	// Обработчик перетаскивания файлов (Drag & Drop)
	const handleDrop = useCallback(
		(event: React.DragEvent<HTMLDivElement>, directoryId: number) => {
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
		},
		[createFileMutation, findFileInDirectory, treeItems]
	);

	// Глобальный обработчик кликов для снятия выделения
	useEffect(() => {
		const handleClickOutside = (event: MouseEvent) => {
			// Проверяем, был ли клик вне tree view
			if (
				treeViewRef.current &&
				!treeViewRef.current.contains(event.target as Node)
			) {
				// Если клик был вне дерева и не на контекстном меню
				const contextMenu = document.querySelector('.MuiMenu-paper');
				if (!contextMenu || !contextMenu.contains(event.target as Node)) {
					// Снимаем выделение
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

	// Отображение состояния загрузки
	if (isLoading) {
		return (
			<LoadingState 
				message={`Загрузка ${isArchive ? 'архивного' : 'рабочего'} хранилища...`}
			/>
		);
	}

	// Отображение ошибки
	if (isError) {
		return (
			<ErrorState
				error={error}
				onRetry={refreshTree}
			/>
		);
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
					height: '57px', // Фиксированная высота для соответствия заголовку предпросмотра
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

			{/* Контейнер для дерева */}
			<Box sx={{ p: 2, maxHeight: 600, overflowY: 'auto' }}>
				<RichTreeView
					ref={treeViewRef}
					items={treeItems}
					defaultExpandedItems={['dir-1']}
					// Отключаем стандартную навигацию по клавиатуре и автофокус
					disableSelection={true}
					autoFocus={false}
					slots={{
						item: (props: ExtendedTreeItem2Props) => {
							const { itemId = '', label, expansionState, ...rest } = props;

							// Поиск данных элемента
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

							return (
								<TreeItem2
									{...rest}
									itemId={itemId}
									className={`tree-item-component ${
										isSelected ? 'custom-selected' : ''
									}`}
									// Отключаем стандартную обработку клавиш навигации
									tabIndex={-1}
									// Обработчик правого клика для контекстного меню
									onContextMenu={event => {
										handleContextMenu(
											event,
											itemId,
											itemData.type,
											itemData.label
										);
										return false;
									}}
									// Убираем стандартный onClick
									onClick={e => {
										e.preventDefault();
										e.stopPropagation();
									}}
									label={
										<Box
											display='flex'
											alignItems='center'
											gap={1.5}
											// Ключевое изменение: обработчик клика для всех элементов
											onClick={e => {
												e.stopPropagation(); // Остановка всплытия для предотвращения двойной обработки
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
											}}
											// Обработчики перетаскивания файлов (только для директорий)
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
											{/* Название элемента */}
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
						// Стили для визуального оформления дерева
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
							// Важно: переопределяем стили для фокуса
							'&.Mui-focused': {
								outline: 'none !important',
								backgroundColor: 'transparent !important',
							},
							// Важно: стили для .Mui-selected тоже переопределяем
							'&.Mui-selected': {
								backgroundColor: 'transparent !important',
								outline: 'none !important',
							},
							// Элементы с нашим собственным классом выбора
							'&.custom-selected > .MuiTreeItem-content': {
								backgroundColor: 'transparent !important',
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
						// Важно: убираем синий фокус с элементов
						'& .MuiTreeItem-content.Mui-focused': {
							backgroundColor: 'transparent !important',
						},
						'& .MuiTreeItem-content.Mui-selected': {
							backgroundColor: 'transparent !important',
						},
					}}
				/>
			</Box>
		</Paper>
	);
};

export default FilesTree;
