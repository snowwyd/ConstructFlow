import { Box } from '@mui/material';
import ContextMenu from './ContextMenu';
import FilesTree from './FilesTree';

const MainPage = () => {
	return (
		<Box className='main-page-container' sx={{ p: 2 }}>
			<Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
				<FilesTree isArchive={false} />
				<FilesTree isArchive={true} />
				<ContextMenu />
			</Box>
		</Box>
	);
};

export default MainPage;
