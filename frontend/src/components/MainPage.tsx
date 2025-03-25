import FilesTree from './FilesTree';

const MainPage = () => {
	return (
		<div>
			<FilesTree isArchive={false} />
			<FilesTree isArchive={true} />
		</div>
	);
};

export default MainPage;
