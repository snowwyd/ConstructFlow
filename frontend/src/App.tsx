import { Provider } from 'react-redux';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import NavigationWrapper from './NavigationWrapper';
import Auth from './components/Auth';
import MainPage from './components/MainPage';
import { store } from './store/store';

const App = () => {
	return (
		<Provider store={store}>
			<BrowserRouter>
				<NavigationWrapper>
					<Routes>
						<Route path='/' element={<Auth />} />
						<Route path='/main' element={<MainPage />} />
					</Routes>
				</NavigationWrapper>
			</BrowserRouter>
		</Provider>
	);
};

export default App;
