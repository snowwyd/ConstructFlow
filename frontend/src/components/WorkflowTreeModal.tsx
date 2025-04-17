import { FolderOutlined } from '@mui/icons-material';
import CloseIcon from '@mui/icons-material/Close';
import {
	alpha,
	Box,
	Button,
	CircularProgress,
	Dialog,
	DialogActions,
	DialogContent,
	DialogTitle,
	Divider,
	IconButton,
	List,
	ListItem,
	ListItemIcon,
	ListItemText,
	Typography,
	useTheme,
} from '@mui/material';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import React, { useState } from 'react';
import { axiosFetching, axiosFetchingFiles } from '../api/AxiosFetch';

interface WorkflowTreeModalProps {
	open: boolean;
	onClose: () => void;
	workflowId: number | null;
	workflowName: string;
}

interface Directory {
	directory_id: number;
	name_folder: string;
	parent_path_id?: number | null;
	status: string;
	current_workflow_assigned?: boolean;
}

/**
 * Модальное окно для просмотра и назначения шаблона согласования на директории
 */
const WorkflowTreeModal: React.FC<WorkflowTreeModalProps> = ({
	open,
	onClose,
	workflowId,
	workflowName,
}) => {
	const theme = useTheme();
	const queryClient = useQueryClient();
	const [selectedDirectories, setSelectedDirectories] = useState<number[]>([]);
	const [isSaving, setIsSaving] = useState(false);

	// Запрос для получения дерева директорий привязанных к workflow
	const {
		data: workflowTree,
		isLoading,
		refetch,
	} = useQuery({
		queryKey: ['admin', 'workflowTree', workflowId],
		queryFn: async () => {
			if (!workflowId) return { directories: [] };

			// Используем правильный путь к API из микросервиса файлов
			// GET "/admin/workflows/:workflow_id/tree"
			const response = await axiosFetchingFiles.get(
				`/admin/workflows/${workflowId}/tree`
			);
			return response.data;
		},
		enabled: open && !!workflowId,
	});

	// Мутация для назначения workflow на директории
	const assignWorkflowMutation = useMutation({
		mutationFn: async (directoryIds: number[]) => {
			if (!workflowId) throw new Error('Workflow ID is required');

			// Используем правильный эндпоинт ЯДРА
			// PUT "/admin/workflows/:workflow_id/assign"
			// с правильной структурой данных { directory_ids: [...] }
			const response = await axiosFetching.put(
				`/admin/workflows/${workflowId}/assign`,
				{ directory_ids: directoryIds }
			);
			return response.data;
		},
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: ['admin', 'workflowTree', workflowId],
			});
			queryClient.invalidateQueries({ queryKey: ['admin', 'assignmentRules'] });
			refetch();
			setSelectedDirectories([]);
		},
	});

	// Обработчик выбора директории
	const handleDirectorySelect = (directoryId: number) => {
		setSelectedDirectories(prev => {
			if (prev.includes(directoryId)) {
				return prev.filter(id => id !== directoryId);
			} else {
				return [...prev, directoryId];
			}
		});
	};

	// Обработчик назначения шаблона на выбранные директории
	const handleAssignWorkflow = async () => {
		if (selectedDirectories.length === 0 || !workflowId) return;

		setIsSaving(true);
		try {
			await assignWorkflowMutation.mutateAsync(selectedDirectories);
		} finally {
			setIsSaving(false);
		}
	};

	const handleClose = () => {
		setSelectedDirectories([]);
		onClose();
	};

	return (
		<Dialog
			open={open}
			onClose={handleClose}
			maxWidth='md'
			fullWidth
			PaperProps={{
				sx: {
					borderRadius: 3,
					boxShadow: `0 8px 32px ${alpha(theme.palette.primary.main, 0.15)}`,
				},
			}}
		>
			<DialogTitle sx={{ pb: 2 }}>
				<Box display='flex' justifyContent='space-between' alignItems='center'>
					<Box display='flex' alignItems='center' gap={1}>
						<FolderOutlined color='primary' />
						<Typography variant='h6' fontWeight={600}>
							Дерево директорий шаблона {workflowName}
						</Typography>
					</Box>
					<IconButton onClick={handleClose}>
						<CloseIcon />
					</IconButton>
				</Box>
			</DialogTitle>
			<Divider />
			<DialogContent sx={{ pt: 3, minHeight: 400 }}>
				{isLoading ? (
					<Box
						sx={{
							display: 'flex',
							flexDirection: 'column',
							justifyContent: 'center',
							alignItems: 'center',
							height: 300,
							gap: 2,
						}}
					>
						<CircularProgress />
						<Typography variant='body2' color='text.secondary'>
							Загрузка директорий...
						</Typography>
					</Box>
				) : workflowTree?.directories?.length > 0 ? (
					<>
						<Typography variant='body2' color='text.secondary' sx={{ mb: 2 }}>
							Выберите директории, для которых нужно назначить шаблон
							согласования. Текущие назначения отмечены иконкой.
						</Typography>
						<List>
							{workflowTree.directories.map((dir: Directory) => (
								<ListItem
									key={dir.directory_id}
									sx={{
										borderBottom: `1px solid ${alpha(
											theme.palette.divider,
											0.1
										)}`,
										py: 1,
										backgroundColor: selectedDirectories.includes(
											dir.directory_id
										)
											? alpha(theme.palette.primary.main, 0.08)
											: 'transparent',
									}}
									button
									onClick={() => handleDirectorySelect(dir.directory_id)}
								>
									<ListItemIcon>
										<FolderOutlined
											color={
												dir.current_workflow_assigned ? 'primary' : 'action'
											}
										/>
									</ListItemIcon>
									<ListItemText
										primary={dir.name_folder}
										secondary={`ID: ${dir.directory_id} ${
											dir.current_workflow_assigned
												? '• Шаблон уже назначен'
												: ''
										}`}
									/>
								</ListItem>
							))}
						</List>

						{selectedDirectories.length > 0 && (
							<Box
								sx={{
									mt: 2,
									p: 2,
									bgcolor: alpha(theme.palette.info.light, 0.1),
									borderRadius: 2,
								}}
							>
								<Typography variant='body2'>
									Выбрано директорий: {selectedDirectories.length}
								</Typography>
							</Box>
						)}
					</>
				) : (
					<Box sx={{ textAlign: 'center', py: 4 }}>
						<FolderOutlined
							sx={{
								fontSize: 60,
								color: alpha(theme.palette.text.secondary, 0.2),
								mb: 2,
							}}
						/>
						<Typography color='text.secondary'>
							Для этого шаблона согласования не найдено связанных директорий
						</Typography>
					</Box>
				)}
			</DialogContent>
			<Divider />
			<DialogActions sx={{ p: 2 }}>
				<Button
					onClick={handleClose}
					variant='outlined'
					sx={{ borderRadius: 2 }}
				>
					Закрыть
				</Button>

				{selectedDirectories.length > 0 && (
					<Button
						onClick={handleAssignWorkflow}
						variant='contained'
						sx={{ borderRadius: 2 }}
						disabled={isSaving}
						startIcon={isSaving ? <CircularProgress size={20} /> : null}
					>
						{isSaving ? 'Назначение...' : 'Назначить шаблон'}
					</Button>
				)}
			</DialogActions>
		</Dialog>
	);
};

export default WorkflowTreeModal;
