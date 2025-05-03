import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    mode: 'dark', // #ТемнаяТема
    primary: {
      main: '#2196f3', // Синий
      light: '#4dabf5',
      dark: '#1976d2',
    },
    secondary: {
      main: '#ffffff', // Белый 
    },
    background: {
      default: '#000000', // Чёрный 
      paper: '#121212',   // Темно-серый 
    },
    text: {
      primary: '#ffffff',     // Белый текст
      secondary: '#aaaaaa',   // Серый текст 
    },
    success: {
      main: '#4caf50', //  Зелёный для успеха
    },
    warning: {
      main: '#ff9800', // Оранжевый для предупреждений
    },
    error: {
      main: '#f44336', // Красный для ошибок
    },
    info: {
      main: '#2196f3', // Инфо 
    },
  },
  typography: {
    fontFamily: [
      '-apple-system',
      'BlinkMacSystemFont',
      '"Segoe UI"',
      'Roboto',
      '"Helvetica Neue"',
      'Arial',
      'sans-serif',
      '"Apple Color Emoji"',
      '"Segoe UI Emoji"',
      '"Segoe UI Symbol"',
    ].join(','),
  },
  components: {
    MuiAppBar: {
      styleOverrides: {
        root: {
          boxShadow:
            '0 2px 6px rgba(0, 0, 0, 0.3), 0 1px 2px rgba(255, 255, 255, 0.05)',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 8,
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          backgroundColor: '#1e1e1e',
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
        },
      },
    },
    MuiTextField: {
      defaultProps: {
        variant: 'outlined',
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          backgroundColor: '#2a2a2a',
          borderRadius: 8,
        },
        notchedOutline: {
          borderColor: 'rgba(255, 255, 255, 0.2)',
        },
      },
    },
    MuiInputLabel: {
      styleOverrides: {
        root: {
          color: 'rgba(255, 255, 255, 0.7)',
        },
      },
    },
    MuiSelect: {
      styleOverrides: {
        select: {
          '&:focus': {
            backgroundColor: '#2a2a2a',
          },
        },
      },
    },
  },
});

export default theme;


// Light theme
// import { createTheme } from '@mui/material/styles';

// const theme = createTheme({
// 	palette: {
// 		primary: {
// 			main: '#1976d2',
// 			light: '#42a5f5',
// 			dark: '#1565c0',
// 		},
// 		secondary: {
// 			main: '#9c27b0',
// 			light: '#ba68c8',
// 			dark: '#7b1fa2',
// 		},
// 	},
// 	typography: {
// 		fontFamily: [
// 			'-apple-system',
// 			'BlinkMacSystemFont',
// 			'"Segoe UI"',
// 			'Roboto',
// 			'"Helvetica Neue"',
// 			'Arial',
// 			'sans-serif',
// 			'"Apple Color Emoji"',
// 			'"Segoe UI Emoji"',
// 			'"Segoe UI Symbol"',
// 		].join(','),
// 	},
// 	components: {
// 		MuiAppBar: {
// 			styleOverrides: {
// 				root: {
// 					boxShadow: '0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24)',
// 				},
// 			},
// 		},
// 	},
// });

// export default theme;
