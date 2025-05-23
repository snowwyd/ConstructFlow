import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import BuildCircleIcon from '@mui/icons-material/BuildCircle';
import ConstructionIcon from '@mui/icons-material/Construction';
import FolderOutlinedIcon from '@mui/icons-material/FolderOutlined';
import LogoutIcon from '@mui/icons-material/Logout';
import MenuIcon from '@mui/icons-material/Menu';
import PeopleOutlineIcon from '@mui/icons-material/PeopleOutline';
import {
	alpha,
	AppBar,
	Avatar,
	Badge,
	Box,
	Button,
	Divider,
	IconButton,
	Menu,
	MenuItem,
	Toolbar,
	Typography,
	useTheme,
} from '@mui/material';
import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {axiosFetching} from '../api/AxiosFetch';
import { redirectToLogin } from '../api/NavigationService';
import config from '../constants/Configurations.json';
import { ApprovalResponse } from '../interfaces/Approvals';

// Interface for UserInfo
interface UserInfo {
	id: number;
	login: string;
	role: string;
}

const Header = () => {
	const theme = useTheme();
	const location = useLocation();
	const navigate = useNavigate();
	const isLoginPage = location.pathname === '/';
	const [isLoggedIn, setIsLoggedIn] = useState(false);
	const [userInfo, setUserInfo] = useState<UserInfo | null>(null);
	const [mobileMenuAnchor, setMobileMenuAnchor] = useState<null | HTMLElement>(
		null
	);
	const [userMenuAnchor, setUserMenuAnchor] = useState<null | HTMLElement>(
		null
	);
	const [pendingApprovals, setPendingApprovals] = useState(0);

	const mobileMenuOpen = Boolean(mobileMenuAnchor);
	const userMenuOpen = Boolean(userMenuAnchor);

	// Check if user is admin
	const isAdmin = userInfo?.role === 'admin';

	// Check auth status and load user data
	useEffect(() => {
		const checkAuthStatus = async () => {
			if (isLoginPage) {
				setIsLoggedIn(false);
				return;
			}

			try {
				const response = await axiosFetching.get(config.checkJWT);
				if (response.data && response.data.id) {
					setIsLoggedIn(true);
					setUserInfo(response.data);

					// Fetch pending approvals if user is logged in
					fetchPendingApprovals();
				} else {
					setIsLoggedIn(false);
				}
			} catch (error) {
				console.error('Auth check error:', error);
				setIsLoggedIn(false);
			}
		};

		checkAuthStatus();
	}, [location.pathname, isLoginPage]);

	/**
	 * Fetches pending approvals that require user attention
	 * Updates the global state with pending approvals count
	 */
	const fetchPendingApprovals = async () => {
		try {
			console.log('Fetching pending approvals');
			const response = await axiosFetching.get(config.getApprovals);

			// Make sure we have valid response data
			if (response.data && Array.isArray(response.data)) {
				// Use the proper type instead of 'any'
				const pendingCount = response.data.filter(
					(approval: ApprovalResponse) => approval.status === 'on approval'
				).length;

				console.log(`Found ${pendingCount} pending approvals`);
				setPendingApprovals(pendingCount);
			} else {
				console.warn('Unexpected response format:', response.data);
				setPendingApprovals(0);
			}
		} catch (error) {
			console.error('Error fetching approvals:', error);
			// Default to 0 in case of error
			setPendingApprovals(0);
		}
	};

	// Add event listener for global approval count updates
	useEffect(() => {
		// Define the event handler
		const handleUpdateApprovalCount = () => {
			console.log('Approval count update requested');
			fetchPendingApprovals();
		};

		// Add event listener
		window.addEventListener('update-approval-count', handleUpdateApprovalCount);

		// Clean up
		return () => {
			window.removeEventListener(
				'update-approval-count',
				handleUpdateApprovalCount
			);
		};
	}, []);

	// Menu handlers
	const handleMobileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
		setMobileMenuAnchor(event.currentTarget);
	};

	const handleMobileMenuClose = () => {
		setMobileMenuAnchor(null);
	};

	const handleUserMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
		setUserMenuAnchor(event.currentTarget);
	};

	const handleUserMenuClose = () => {
		setUserMenuAnchor(null);
	};

	// Logout handler - client-side approach
	const handleLogout = () => {
		try {
			// Clear cookies
			document.cookie.split(';').forEach(c => {
				document.cookie = c
					.replace(/^ +/, '')
					.replace(/=.*/, `=;expires=${new Date().toUTCString()};path=/`);
			});

			localStorage.removeItem('token');
			handleUserMenuClose();
			redirectToLogin();
		} catch (error) {
			console.error('Logout error:', error);
		}
	};

	// Helper function to check if a menu item is active
	const isActive = (path: string) => {
		return location.pathname === path;
	};

	// Navigation to route
	const handleNavigation = (path: string) => {
		navigate(path);
		handleMobileMenuClose();
	};

	// Base navigation items with icons
	const navItems = [
		{
			label: 'Файлы',
			icon: <FolderOutlinedIcon sx={{ mr: 1 }} />,
			path: '/main',
			active: isActive('/main'),
		},
		{
			label: 'Согласования',
			icon: (
				<Badge badgeContent={pendingApprovals} color='error' sx={{ mr: 1 }}>
					<AssignmentTurnedInIcon />
				</Badge>
			),
			path: '/approvals',
			active: isActive('/approvals'),
		},
	];

	// Add Admin section if user has admin role
	if (isAdmin) {
		// Add Approval Editor tab for admins
		navItems.push({
			label: 'Редактор согласования',
			icon: <BuildCircleIcon sx={{ mr: 1 }} />,
			path: '/approval-editor',
			active: isActive('/approval-editor'),
		});

		// Keep Users management as the last tab
		navItems.push({
			label: 'Права и пользователи',
			icon: <PeopleOutlineIcon sx={{ mr: 1 }} />,
			path: '/users',
			active: isActive('/users'),
		});
	}

	return (
		<AppBar
			position='static'
			elevation={0}
			sx={{
				backgroundColor:
					theme.palette.mode === 'light'
						? '#ffffff'
						: alpha(theme.palette.background.paper, 0.9),
				borderBottom: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
				color: theme.palette.text.primary,
			}}
		>
			<Toolbar sx={{ minHeight: 64 }}>
				{/* Logo and brand */}
				<Box sx={{ display: 'flex', alignItems: 'center' }}>
					<ConstructionIcon
						color='primary'
						sx={{
							mr: 1.5,
							fontSize: 28,
							transform: 'rotate(15deg)',
						}}
					/>
					<Typography
						variant='h6'
						noWrap
						component='div'
						sx={{
							fontWeight: 700,
							letterSpacing: '-0.5px',
							background: `linear-gradient(90deg, ${theme.palette.primary.main} 0%, ${theme.palette.secondary.main} 100%)`,
							WebkitBackgroundClip: 'text',
							WebkitTextFillColor: 'transparent',
						}}
					>
						ConstructFlow
					</Typography>
				</Box>

				{/* Main content - changes based on login state */}
				<Box sx={{ flexGrow: 1 }}>
					{isLoggedIn && (
						<Box
							sx={{
								display: { xs: 'none', md: 'flex' },
								justifyContent: 'center',
								gap: 2,
							}}
						>
							{navItems.map(item => (
								<Button
									key={item.label}
									color={item.active ? 'primary' : 'inherit'}
									startIcon={item.icon}
									onClick={() => handleNavigation(item.path)}
									sx={{
										py: 1,
										px: 2,
										borderRadius: 2,
										position: 'relative',
										fontWeight: item.active ? 600 : 400,
										backgroundColor: item.active
											? alpha(theme.palette.primary.main, 0.08)
											: 'transparent',
										'&:hover': {
											backgroundColor: item.active
												? alpha(theme.palette.primary.main, 0.12)
												: alpha(theme.palette.primary.main, 0.08),
										},
										'&::after': item.active
											? {
													content: '""',
													position: 'absolute',
													bottom: 5,
													left: '30%',
													width: '40%',
													height: 3,
													borderRadius: 3,
													backgroundColor: theme.palette.primary.main,
											  }
											: {},
									}}
								>
									{item.label}
								</Button>
							))}
						</Box>
					)}
				</Box>

				{/* User section */}
				{isLoggedIn ? (
					<Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
						{/* User profile button */}
						<Button
							onClick={handleUserMenuOpen}
							color='inherit'
							sx={{
								borderRadius: 6,
								textTransform: 'none',
								px: { xs: 1, sm: 2 },
								'&:hover': {
									backgroundColor: alpha(theme.palette.primary.main, 0.08),
								},
							}}
							startIcon={
								<Avatar
									sx={{
										width: 32,
										height: 32,
										bgcolor: theme.palette.primary.main,
									}}
								>
									{userInfo?.login?.charAt(0).toUpperCase() || 'U'}
								</Avatar>
							}
						>
							<Box
								sx={{
									display: { xs: 'none', sm: 'block' },
									textAlign: 'left',
								}}
							>
								<Typography variant='body2' fontWeight={600} noWrap>
									{userInfo?.login || 'User'}
								</Typography>
								<Typography
									variant='caption'
									color='text.secondary'
									sx={{ display: 'block' }}
								>
									{userInfo?.role || 'Role'}
								</Typography>
							</Box>
						</Button>

						{/* Mobile menu button */}
						<IconButton
							size='large'
							edge='end'
							color='inherit'
							aria-label='menu'
							onClick={handleMobileMenuOpen}
							sx={{ display: { xs: 'block', md: 'none' } }}
						>
							<MenuIcon />
						</IconButton>

						{/* User profile menu */}
						<Menu
							anchorEl={userMenuAnchor}
							open={userMenuOpen}
							onClose={handleUserMenuClose}
							slotProps={{
								paper: {
									elevation: 3,
									sx: {
										mt: 1.5,
										overflow: 'visible',
										filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.08))',
										minWidth: 200,
										'&:before': {
											content: '""',
											display: 'block',
											position: 'absolute',
											top: 0,
											right: 14,
											width: 10,
											height: 10,
											bgcolor: 'background.paper',
											transform: 'translateY(-50%) rotate(45deg)',
											zIndex: 0,
										},
									},
								},
							}}
							transformOrigin={{ horizontal: 'right', vertical: 'top' }}
							anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
						>
							<Box sx={{ px: 2, py: 1.5 }}>
								<Typography variant='subtitle1' fontWeight={600}>
									{userInfo?.login || 'User'}
								</Typography>
								<Typography variant='body2' color='text.secondary'>
									{userInfo?.role || 'Role'}
								</Typography>
							</Box>
							<Divider />
							<MenuItem onClick={handleUserMenuClose} sx={{ py: 1.5 }}>
								<AccountCircleIcon fontSize='small' sx={{ mr: 2 }} />
								Мой профиль
							</MenuItem>
							<MenuItem onClick={handleLogout} sx={{ py: 1.5 }}>
								<LogoutIcon fontSize='small' sx={{ mr: 2 }} />
								Выйти
							</MenuItem>
						</Menu>

						{/* Mobile navigation menu */}
						<Menu
							anchorEl={mobileMenuAnchor}
							open={mobileMenuOpen}
							onClose={handleMobileMenuClose}
							slotProps={{
								paper: {
									elevation: 3,
									sx: {
										mt: 1.5,
										minWidth: 200,
									},
								},
							}}
							transformOrigin={{ horizontal: 'right', vertical: 'top' }}
							anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
						>
							{navItems.map(item => (
								<MenuItem
									key={item.label}
									onClick={() => handleNavigation(item.path)}
									sx={{
										py: 1.5,
										backgroundColor: item.active
											? alpha(theme.palette.primary.main, 0.08)
											: 'transparent',
										'&:hover': {
											backgroundColor: item.active
												? alpha(theme.palette.primary.main, 0.12)
												: alpha(theme.palette.primary.light, 0.08),
										},
										borderLeft: item.active
											? `3px solid ${theme.palette.primary.main}`
											: '3px solid transparent',
									}}
								>
									{item.icon}
									<Typography
										variant='body1'
										color={item.active ? 'primary' : 'inherit'}
										fontWeight={item.active ? 600 : 400}
									>
										{item.label}
									</Typography>
								</MenuItem>
							))}
						</Menu>
					</Box>
				) : null}
			</Toolbar>
		</AppBar>
	);
};

export default Header;
