import { useQuery } from "@tanstack/react-query";
import config from '../constants/Configurations.json';
import axiosFetching from "../api/AxiosFetch";
import ErrorState from "./ErrorState";
import LoadingState from "./LoadingState";
import { useDispatch } from "react-redux";
import { useEffect } from "react";
import { CheckCircleOutline, CancelOutlined } from '@mui/icons-material';
import { setPendingCount } from "../store/Slices/pendingApprovalsSlice";
import { Box, IconButton, Tooltip } from "@mui/material";

const getApprovals = config.getApprovals;

interface ResponseProps{
    id: number,
    file_id: number,
    file_name: string,
    status: string,
    workflow_order: number
}

export const ApprovalsPage = () => {
    
    const dispatch = useDispatch();
    const isAdmin = true;
    const {
        data: apiResponse,
		isLoading,
		isError,
	} = useQuery({
		queryKey: ['approvals'],
		queryFn: async () => {
			const response = await axiosFetching.get(getApprovals);
			return response.data;
		},
	});

    useEffect(() => {
        if (apiResponse) {
          const pendingCount = apiResponse.filter(
            (doc: ResponseProps) => doc.status === 'on approval'
          ).length;
          dispatch(setPendingCount(pendingCount)); 
        }
      }, [apiResponse, dispatch]);

    if (isLoading){
        return <LoadingState/>
    };
    if(isError){
        return <ErrorState/>
    };
    return (
        <div>
            <ul style={{ listStyle: 'none', padding: 0 }}>
                {apiResponse.map((document: ResponseProps) => (
                    <li
                        key={document.id}
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
                            <p>Status: {document.status}</p>
                            <p>Workflow Order: {document.workflow_order}</p>
                        </div>

                        {isAdmin && (
                            <Box display="flex" gap={1}>
                                <Tooltip title="Approve">
                                    <IconButton
                                        sx={{
                                            border: '2px solid green',
                                            color: 'green',
                                            backgroundColor: 'transparent',
                                            borderRadius: '50%',
                                            '&:hover': {
                                                backgroundColor: 'rgba(0, 255, 0, 0.1)',
                                            },
                                        }}
                                    >
                                        <CheckCircleOutline />
                                    </IconButton>
                                </Tooltip>

                                <Tooltip title="Reject">
                                    <IconButton
                                        sx={{
                                            border: '2px solid red',
                                            color: 'red',
                                            backgroundColor: 'transparent',
                                            borderRadius: '50%',
                                            '&:hover': {
                                                backgroundColor: 'rgba(255, 0, 0, 0.1)',
                                            },
                                        }}
                                    >
                                        <CancelOutlined />
                                    </IconButton>
                                </Tooltip>
                            </Box>
                        )}
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default ApprovalsPage;
