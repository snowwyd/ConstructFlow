import axios from 'axios';
import config from '../constants/Configurations.json';
import { useNavigate } from 'react-router-dom';

const serverURL = config.serverUrl;

const axiosFetching = axios.create({
	baseURL: serverURL,
	withCredentials: true,
});

axiosFetching.interceptors.response.use(
	(response) => response,
	(error) => {
	  const navigate = useNavigate(); 
	  if (error.response && error.response.status === 401) {
		navigate('/'); 
	  }
	}
  );

export default axiosFetching;
