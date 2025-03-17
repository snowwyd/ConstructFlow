import { BrowserRouter, Routes, Route } from 'react-router';
import Auth from './components/Auth';
import MainPage from './components/MainPage';

const App = () => {
	return (
		<BrowserRouter>
		<Routes>
			<Route path='/' element={<Auth/>}/>
			<Route path='/main' element={<MainPage/>}/>
		</Routes>
		</BrowserRouter>
		
	)
};

export default App;
