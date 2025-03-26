// MainPage.tsx
import ContextMenu from './ContextMenu';
import FilesTree from './FilesTree';

const MainPage = () => {
	return (
		<div className='main-page-container'>
			<FilesTree isArchive={false} />
			<FilesTree isArchive={true} />
			<ContextMenu />
		</div>
	);
};

export default MainPage;
