import CssBaseline from '@mui/material/CssBaseline';
import { ThemeProvider } from '@mui/material/styles';
import { Provider } from 'react-redux';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import NavigationWrapper from './NavigationWrapper';
import ApprovalEditorPage from './components/ApprovalEditorPage';
import ApprovalsPage from './components/ApprovalsPage';
import Auth from './components/Auth';
import MainPage from './components/MainPage';
import UsersPermissionsPage from './components/UsersPermissionsPage';
import { store } from './store/store';
import theme from './constants/theme';

const App = () => {
	return (
		<Provider store={store}>
			<ThemeProvider theme={theme}>
				<CssBaseline />
				<BrowserRouter>
					<NavigationWrapper>
						<Routes>
							<Route path='/' element={<Auth />} />
							<Route path='/main' element={<MainPage />} />
							<Route path='/approvals' element={<ApprovalsPage />} />
							<Route path='/users' element={<UsersPermissionsPage />} />
							<Route path='/approval-editor' element={<ApprovalEditorPage />} />
						</Routes>
					</NavigationWrapper>
				</BrowserRouter>
			</ThemeProvider>
		</Provider>
	);
};

export default App;
