//Сущность аксиос, настройка URL для получения даты. 
// Важно. Определить ендпоинты получения JWT от сервера. 
// Важно. Создать папку для хранения констант типа url сервер
import axios from "axios";

const axiosFetching = axios.create(
    {
    baseURL: "https://danya-sdelat-servak-pomenyate.com", 
    withCredentials: true
    }
); 

export default axiosFetching;