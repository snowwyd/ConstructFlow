import ConstructionIcon from '@mui/icons-material/Construction';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import PersonOutlineIcon from '@mui/icons-material/PersonOutline';
import VisibilityOffOutlinedIcon from '@mui/icons-material/VisibilityOffOutlined';
import VisibilityOutlinedIcon from '@mui/icons-material/VisibilityOutlined';
import {
    Alert,
    alpha,
    Box,
    Button,
    CircularProgress,
    IconButton,
    InputAdornment,
    Paper,
    TextField,
    Typography,
    useTheme,
} from '@mui/material';
import { useMutation } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import { useState } from 'react';
import { useNavigate } from 'react-router';
import { axiosFetching } from '../api/AxiosFetch';
import config from '../constants/Configurations.json';

const loginEndpoint = config.loginEndpoint;
const JWTresponse = config.checkJWT;

// Type for login credentials
interface LoginCredentials {
    login: string;
    password: string;
}

// Type for JWT validation response
interface JWTValidationResponse {
    id: number;
    login: string;
    role: string;
}

const Auth = () => {
    const theme = useTheme();
    const navigate = useNavigate();
    const [login, setLogin] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [showPassword, setShowPassword] = useState<boolean>(false);

    const { mutate, isPending } = useMutation({
        mutationFn: async () => {
            const response = await axiosFetching.post<{ token: string }>(
                loginEndpoint,
                { login, password } as LoginCredentials
            );
            return response.data;
        },
        onSuccess: async () => {
            try {
                const validateResponse = await axiosFetching.get<JWTValidationResponse>(
                    JWTresponse
                );

                if (validateResponse.data.id) {
                    setError(null);
                    navigate('/main');
                }
            } catch (error) {
                const axiosError = error as AxiosError<{ message?: string }>;
                setError(
                    axiosError.response?.data?.message || 'Ошибка валидации токена'
                );
            }
        },
        onError: (error: AxiosError<{ message?: string }>) => {
            setError(
                error.response?.data?.message ||
                    error.message ||
                    'Произошла ошибка при входе'
            );
        },
    });

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        if (!login || !password) {
            setError('Заполните все поля');
            return;
        }
        mutate();
    };

    const togglePasswordVisibility = () => {
        setShowPassword(!showPassword);
    };

    return (
		<Box
			sx={{
				minHeight: '90vh',
				display: 'flex',
				alignItems: 'center',
				justifyContent: 'center',
				backgroundImage: `radial-gradient(circle at 50% 14em, ${alpha(
					theme.palette.primary.light,
					0.15
				)} 0%, ${alpha(
					theme.palette.primary.main,
					0.05
				)} 60%, transparent 100%)`,
				padding: 2,
				flexDirection: 'column', // Изменяем направление на колонку
			}}
		>
			{/* Логотип ConstructHub */}
			<Box
				sx={{
					textAlign: 'center',
					mb: 4, // Отступ снизу для разделения логотипа и формы
				}}
			>
				<ConstructionIcon
					sx={{
						fontSize: 48,
						color: theme.palette.primary.main,
						mb: 1,
						transform: 'rotate(15deg)',
					}}
				/>
				<Typography
					variant='h4'
					component='h1'
					color={theme.palette.primary.main}
					fontWeight={700}
				>
					ConstructHub
				</Typography>
			</Box>

			{/* Форма авторизации */}
			<Paper
				elevation={6}
				sx={{
					maxWidth: 450,
					width: '100%',
					borderRadius: 4,
					overflow: 'hidden',
					boxShadow: `0 10px 40px ${alpha(theme.palette.primary.main, 0.1)}`,
				}}
			>
				{/* Header section */}
				<Box
					sx={{
						bgcolor: theme.palette.primary.main,
						py: 2, // Уменьшен паддинг
						px: 3,
						textAlign: 'center',
						position: 'relative',
						overflow: 'hidden',
					}}
				>
					<Typography variant='body1' color={alpha('#fff', 0.8)}>
						Вход в систему
					</Typography>
				</Box>

				{/* Form section */}
				<Box component='form' onSubmit={handleSubmit} sx={{ p: 4 }}>
					{error && (
						<Alert
							severity='error'
							sx={{
								mb: 3,
								borderRadius: 2,
								animation: 'fadeIn 0.4s ease-out',
								'@keyframes fadeIn': {
									'0%': { opacity: 0, transform: 'translateY(-10px)' },
									'100%': { opacity: 1, transform: 'translateY(0)' },
								},
							}}
						>
							{error}
						</Alert>
					)}

					<TextField
						fullWidth
						id='login'
						label='Логин'
						variant='outlined'
						margin='normal'
						autoComplete='username'
						value={login}
						onChange={(e) => setLogin(e.target.value)}
						InputProps={{
							startAdornment: (
								<InputAdornment position='start'>
									<PersonOutlineIcon color='action' />
								</InputAdornment>
							),
						}}
						sx={{
							mb: 3,
							'& .MuiOutlinedInput-root': {
								borderRadius: 2,
							},
						}}
					/>

					<TextField
						fullWidth
						id='password'
						label='Пароль'
						type={showPassword ? 'text' : 'password'}
						variant='outlined'
						margin='normal'
						autoComplete='current-password'
						value={password}
						onChange={(e) => setPassword(e.target.value)}
						InputProps={{
							startAdornment: (
								<InputAdornment position='start'>
									<LockOutlinedIcon color='action' />
								</InputAdornment>
							),
							endAdornment: (
								<InputAdornment position='end'>
									<IconButton
										aria-label='toggle password visibility'
										onClick={togglePasswordVisibility}
										edge='end'
									>
										{showPassword ? (
											<VisibilityOffOutlinedIcon />
										) : (
											<VisibilityOutlinedIcon />
										)}
									</IconButton>
								</InputAdornment>
							),
						}}
						sx={{
							mb: 4,
							'& .MuiOutlinedInput-root': {
								borderRadius: 2,
							},
						}}
					/>

					<Button
						type='submit'
						fullWidth
						variant='contained'
						size='large'
						disabled={isPending}
						sx={{
							py: 1.5,
							borderRadius: 2,
							textTransform: 'none',
							fontSize: '1rem',
							fontWeight: 600,
							boxShadow: `0 4px 12px ${alpha(theme.palette.primary.main, 0.3)}`,
							position: 'relative',
							overflow: 'hidden',
							transition: 'all 0.3s ease',
							'&:hover': {
								boxShadow: `0 6px 16px ${alpha(
									theme.palette.primary.main,
									0.4
								)}`,
								transform: 'translateY(-2px)',
							},
							'&:active': {
								transform: 'translateY(0)',
								boxShadow: `0 2px 8px ${alpha(
									theme.palette.primary.main,
									0.3
								)}`,
							},
							'&::before': {
								content: '""',
								position: 'absolute',
								top: 0,
								left: '-100%',
								width: '100%',
								height: '100%',
								background: `linear-gradient(90deg, transparent, ${alpha(
									'#fff',
									0.2
								)}, transparent)`,
								transition: 'all 0.6s ease',
							},
							'&:hover::before': {
								left: '100%',
							},
						}}
					>
						{isPending ? (
							<CircularProgress size={24} color='inherit' />
						) : (
							'Войти'
						)}
					</Button>

					{/* Надпись под кнопкой */}
					<Typography
						variant='body2'
						color='text.secondary'
						align='center'
						sx={{ mt: 2 }}
					>
						Если вы забыли пароль или логин, обратитесь к администратору
					</Typography>
				</Box>
			</Paper>
		</Box>
    );
};

export default Auth;	