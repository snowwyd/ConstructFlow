import axios from 'axios';
import config from '../constants/Configurations.json';

const serverURL = config.apiKey;

const axiosFetching = axios.create({
	baseURL: serverURL,
	withCredentials: true,
});

export default axiosFetching;
