// This file creates a navigation service that can be used outside of React components

let navigateFunction: (path: string) => void;

/**
 * Set the navigation function for use throughout the application
 * @param navigate React Router's navigate function
 */
export const setNavigateFunction = (navigate: (path: string) => void) => {
	navigateFunction = navigate;
};

/**
 * Redirect to login page
 * Used when session expires or user is not authenticated
 */
export const redirectToLogin = () => {
	if (navigateFunction) {
		navigateFunction('/');
	} else {
		// Fallback if navigate function is not set
		window.location.href = '/';
	}
};

/**
 * Trigger the update of approvals count globally
 * This function dispatches a custom event that Header.tsx listens for
 */
export const updateApprovalsCount = () => {
	console.log('Triggering global approvals count update');
	window.dispatchEvent(new Event('update-approval-count'));
};

/**
 * Navigate to the approvals page
 */
export const navigateToApprovals = () => {
	if (navigateFunction) {
		navigateFunction('/approvals');
	} else {
		window.location.href = '/approvals';
	}
};

/**
 * Generic navigation helper
 * @param path The path to navigate to
 */
export const navigateTo = (path: string) => {
	if (navigateFunction) {
		navigateFunction(path);
	} else {
		window.location.href = path;
	}
};
