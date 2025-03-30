// This file creates a navigation service that can be used outside of React components

let navigateFunction: (path: string) => void;

export const setNavigateFunction = (navigate: (path: string) => void) => {
	navigateFunction = navigate;
};

export const redirectToLogin = () => {
	if (navigateFunction) {
		navigateFunction('/');
	} else {
		// Fallback if navigate function is not set
		window.location.href = '/';
	}
};
