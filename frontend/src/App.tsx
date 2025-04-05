import { Box, Paper, Typography, alpha } from '@mui/material';
import CssBaseline from '@mui/material/CssBaseline';
import { ThemeProvider } from '@mui/material/styles';
import { Provider } from 'react-redux';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import NavigationWrapper from './NavigationWrapper';
import ApprovalEditorPage from './components/ApprovalEditorPage';
import Auth from './components/Auth';
import MainPage from './components/MainPage';
import UsersPermissionsPage from './components/UsersPermissionsPage';
import { store } from './store/store';
import theme from './theme';

// Placeholder component for Approvals page
const ApprovalsPage = () => (
	<Box sx={{ p: 4 }}>
		<Paper
			elevation={2}
			sx={{
				p: 4,
				borderRadius: 3,
				textAlign: 'center',
				maxWidth: 800,
				mx: 'auto',
				bgcolor: theme => alpha(theme.palette.primary.light, 0.05),
			}}
		>
			<Typography variant='h4' gutterBottom>
				Система согласований
			</Typography>
			<Typography variant='body1' color='text.secondary'>
				Эта страница находится в разработке. Здесь будет отображаться список
				согласований, требующих вашего внимания или решения.
			</Typography>
		</Paper>
	</Box>
);

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
