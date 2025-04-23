import axios from 'axios';
import config from '../constants/Configurations.json';
import { redirectToLogin } from './NavigationService';

const serverURL = config.serverUrl;
const serverURLFiles = config.serverUrlFiles;

export const axiosFetching = axios.create({
	baseURL: serverURL,
	withCredentials: true,
});

export const axiosFetchingFiles = axios.create({
    baseURL: serverURLFiles,
    withCredentials: true,
});

axiosFetching.interceptors.response.use(
	response => response,
	error => {
		if (error.response && error.response.status === 401) {
			redirectToLogin();
		}
		return Promise.reject(error); 
	}
);

axiosFetchingFiles.interceptors.response.use(
    response => response,
    error => {
        if (error.response && error.response.status === 401) {
            redirectToLogin();
        }
        return Promise.reject(error);
    }
);
