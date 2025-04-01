import { useQuery } from "@tanstack/react-query";
import config from '../constants/Configurations.json';
import axiosFetching from "../api/AxiosFetch";
import ErrorState from "./ErrorState";
import LoadingState from "./LoadingState";
import { useDispatch } from "react-redux";
import { useEffect } from "react";
import { setPendingCount } from "../store/Slices/pendingApprovalsSlice";

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
    <ul>
    {apiResponse.map((document: ResponseProps) => (
      <li key={document.id} style={{ marginBottom: '1rem', padding: '0.5rem', border: '1px solid #ccc' }}>
        <h3>{document.file_name}</h3>
        <p>Status: {document.status}</p>
        <p>Workflow Order: {document.workflow_order}</p>
      </li>
    ))}
  </ul>
  </div>
  )
};

export default ApprovalsPage;
