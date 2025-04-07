import { useMutation, useQuery } from "@tanstack/react-query";
import config from '../constants/Configurations.json';
import axiosFetching from "../api/AxiosFetch";
import ErrorState from "./ErrorState";
import LoadingState from "./LoadingState";
import { useDispatch } from "react-redux";
import { useEffect, useState } from "react";
import { CheckCircleOutline, CancelOutlined } from '@mui/icons-material';
import { setPendingCount } from "../store/Slices/pendingApprovalsSlice";
import { Box, IconButton, Snackbar, Tooltip, Dialog, DialogTitle, DialogContent, TextField, DialogActions, Button } from "@mui/material";

const getApprovals = config.getApprovals;
const approveDocument = config.approveDocument;
const annotateDocument = config.annotateDocument;
const finalizeDocument = config.finalizeDocument;

interface ResponseProps {
    approval_id: number,
    file_id: number,
    file_name: string,
    status: string,
    workflow_order: number,
    workflow_user_count: number
}

export const ApprovalsPage = () => {
    const [snackbarOpen, setSnackbarOpen] = useState(false);
    const [snackbarMessage, setSnackbarMessage] = useState('');
    const [isAnnotationModalOpen, setIsAnnotationModalOpen] = useState(false); 
    const [annotationMessage, setAnnotationMessage] = useState('');
    const [selectedFileId, setSelectedFileId] = useState<number | null>(null); 

    const dispatch = useDispatch();
    const isAdmin = true;

    const {
        data: apiResponse,
        isLoading,
        isError,
        refetch,
    } = useQuery({
        queryKey: ['approvals'],
        queryFn: async () => {
            const response = await axiosFetching.get(getApprovals);
            return response.data;
        },
    });

    useEffect(() => {
        if (apiResponse) {
            const pendingCount = apiResponse.length;
            dispatch(setPendingCount(pendingCount));
        }
    }, [apiResponse, dispatch]);

    const approveDocumentMutation = useMutation({
        mutationFn: async (approvalId: number) => {
            const url = approveDocument.replace(':approval_id', String(approvalId));
            const response = await axiosFetching.put(url);
            return response.data;
        },
        onSuccess: () => {
            setSnackbarMessage('Документ был согласован!');
            setSnackbarOpen(true);
            refetch();
        },
        onError: () => {
            setSnackbarMessage('Ошибка при согласовании документа.');
            setSnackbarOpen(true);
        },
    });

    const annotateDocumentMutation = useMutation({
        mutationFn: async ({ approvalId, message }: { approvalId: number; message: string }) => {
            const url = annotateDocument.replace(':approval_id', String(approvalId));
            const response = await axiosFetching.put(url, { message }); 
            return response.data;
        },
        onSuccess: () => {
            setSnackbarMessage('Документ был отправлен на аннотацию!');
            setSnackbarOpen(true);
            refetch();
            setIsAnnotationModalOpen(false);
            setAnnotationMessage(''); 
        },
        onError: () => {
            setSnackbarMessage('Ошибка при аннотации документа.');
            setSnackbarOpen(true);
        },
    });

    const finalizeDocumentMutation = useMutation({
        mutationFn: async (approvalId: number) => {
            const url = finalizeDocument.replace(':approval_id', String(approvalId));
            const response = await axiosFetching.put(url);
            return response.data;
        },
        onSuccess: () => {
            setSnackbarMessage('Документ был финализирован!');
            setSnackbarOpen(true);
            refetch();
        },
        onError: () => {
            setSnackbarMessage('Ошибка при финализации документа.');
            setSnackbarOpen(true);
        },
    });

    const handleApproveOrFinalize = (document: ResponseProps) => {
        if (document.workflow_order === document.workflow_user_count) {
            finalizeDocumentMutation.mutate(document.approval_id);
        } else {
            approveDocumentMutation.mutate(document.approval_id);
        }
    };

    const handleAnnotateClick = (fileId: number) => {
        setSelectedFileId(fileId); 
        setIsAnnotationModalOpen(true); 
    };

    const handleAnnotationSubmit = () => {
        if (selectedFileId && annotationMessage.trim()) {
            annotateDocumentMutation.mutate({ approvalId: selectedFileId, message: annotationMessage });
        }
    };

    if (isLoading) {
        return <LoadingState />;
    }

    if (isError) {
        return <ErrorState />;
    }

    if (!apiResponse || !Array.isArray(apiResponse)) {
        return <p>Нет документов на согласовании</p>;
    }

    return (
        <div>
            <ul style={{ listStyle: 'none', padding: 0 }}>
                {apiResponse.map((document: ResponseProps) => (
                    <li
                        key={document.approval_id}
                        style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            marginBottom: '1rem',
                            padding: '0.5rem',
                            border: '1px solid #ccc',
                            borderRadius: '4px',
                        }}
                    >
                        <div>
                            <h3>{document.file_name}</h3>
                            <p>Статус: {document.status}</p>
                            <p>Рабочий номер: {document.workflow_order}</p>
                        </div>

                        {isAdmin && (
                            <Box display="flex" gap={1}>
                                {/* Кнопка "Approve/Finalize" */}
                                <Tooltip title={document.workflow_order === document.workflow_user_count ? "Finalize" : "Approve"}>
                                    <IconButton
                                        sx={{
                                            border: `2px solid ${
                                                document.workflow_order === document.workflow_user_count ? 'blue' : 'green'
                                            }`,
                                            color: `${
                                                document.workflow_order === document.workflow_user_count ? 'blue' : 'green'
                                            }`,
                                            backgroundColor: 'transparent',
                                            borderRadius: '15%',
                                            '&:hover': {
                                                backgroundColor: `${
                                                    document.workflow_order === document.workflow_user_count
                                                        ? 'rgba(0, 0, 255, 0.1)'
                                                        : 'rgba(0, 255, 0, 0.1)'
                                                }`,
                                            },
                                        }}
                                        onClick={() => handleApproveOrFinalize(document)}
                                    >
                                        <CheckCircleOutline />
                                    </IconButton>
                                </Tooltip>

                                {/* Кнопка "Annotate" */}
                                <Tooltip title="Annotate">
                                    <IconButton
                                        sx={{
                                            border: '2px solid red',
                                            color: 'red',
                                            backgroundColor: 'transparent',
                                            borderRadius: '15%',
                                            '&:hover': {
                                                backgroundColor: 'rgba(255, 0, 0, 0.1)',
                                            },
                                        }}
                                        onClick={() => handleAnnotateClick(document.approval_id)}
                                    >
                                        <CancelOutlined />
                                    </IconButton>
                                </Tooltip>
                            </Box>
                        )}
                    </li>
                ))}
            </ul>

            {/* Уведомление */}
            <Snackbar
                open={snackbarOpen}
                autoHideDuration={3000}
                onClose={() => setSnackbarOpen(false)}
                message={snackbarMessage}
            />

            {/* Модальное окно для аннотации */}
            <Dialog open={isAnnotationModalOpen} onClose={() => setIsAnnotationModalOpen(false)}>
                <DialogTitle>Аннотация Документа</DialogTitle>
                <DialogContent>
                    <TextField
                        autoFocus
                        margin="dense"
                        label="Message"
                        fullWidth
                        value={annotationMessage}
                        onChange={(e) => setAnnotationMessage(e.target.value)}
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setIsAnnotationModalOpen(false)} color="primary">
                        Отмена
                    </Button>
                    <Button onClick={handleAnnotationSubmit} color="primary" disabled={!annotationMessage.trim()}>
                        Подтвердить
                    </Button>
                </DialogActions>
            </Dialog>
        </div>
    );
};

export default ApprovalsPage;