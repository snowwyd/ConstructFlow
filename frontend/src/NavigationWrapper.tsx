import { ReactNode, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { setNavigateFunction } from './api/NavigationService';

interface NavigationWrapperProps {
	children: ReactNode;
}

const NavigationWrapper = ({ children }: NavigationWrapperProps) => {
	const navigate = useNavigate();

	// Set up the navigation service at the root level
	useEffect(() => {
		setNavigateFunction(navigate);
	}, [navigate]);

	return <>{children}</>;
};

export default NavigationWrapper;
