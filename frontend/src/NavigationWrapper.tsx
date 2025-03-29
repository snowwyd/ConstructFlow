import { ReactNode, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { setNavigateFunction } from './api/NavigationService';
import Header from './components/Header';

interface NavigationWrapperProps {
	children: ReactNode;
}

const NavigationWrapper = ({ children }: NavigationWrapperProps) => {
	const navigate = useNavigate();

	// Set up the navigation service at the root level
	useEffect(() => {
		setNavigateFunction(navigate);
	}, [navigate]);

	return (
		<>
			<Header />
			{children}
		</>
	);
};

export default NavigationWrapper;
