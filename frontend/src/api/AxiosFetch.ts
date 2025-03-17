//Сущность аксиос, настройка URL для получения даты.
// Важно. Определить ендпоинты получения JWT от сервера.
// Важно. Создать папку для хранения констант типа url сервер
import axios from 'axios';
import config from '../constants/Configurations.json';

const serverURL = config.apiKey;

const axiosFetching = axios.create({
	baseURL: serverURL,
	withCredentials: true,
});

export default axiosFetching;
