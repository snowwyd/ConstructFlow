import {Menu, MenuItem} from '@mui/material';
import {useDispatch, useSelector} from "react-redux";
import { closeContextMenu } from '../store/Slices/contexMenuSlice';
import config from "../constants/Configurations.json";


const createFolder =config.createDirectory;

const ContextMenu = () => {
    const dispatch = useDispatch();
    const {mouseX, mouseY, itemId, itemType} = useSelector((state: any) => state.contextMenu)
    
    const handleCloseMenu = () => {
        dispatch(closeContextMenu());
    };

    const handleCreateFolder = () => {

        handleCloseMenu();
      };
    
      const handleCreateFile = () => {
        console.log("Create file for item:", itemId);
        handleCloseMenu();
      };
    
      const handleDeleteFolder = () => {
        console.log("Delete item:", itemId, "Type:", itemType);
        handleCloseMenu();
      };

      const handleDeleteFile = () => {
        console.log("Delete item:", itemId, "Type:", itemType);
        handleCloseMenu();
      };

      const menuItems = [];

      if (itemType === "directory") {
        menuItems.push(
          <MenuItem key="create-folder" onClick={handleCreateFolder}>Создать папку</MenuItem>,
          <MenuItem key="create-file" onClick={handleCreateFile}>Создать файл</MenuItem>,
          <MenuItem key="delete-folder" onClick={handleDeleteFolder}>Удалить папку</MenuItem>
        );
      }
    
      if (itemType === "file") {
        menuItems.push(
          <MenuItem key="delete-file" onClick={handleDeleteFile}>Удалить файл</MenuItem>
        );
      }

  return (
    <div>
        <Menu
		open={mouseX !== null && mouseY !== null}
		onClose={handleCloseMenu}
		anchorReference='anchorPosition'
		anchorPosition={mouseY !== null && mouseX !== null ? { top: mouseY, left: mouseX } : undefined}
		>
            {menuItems.length > 0 ? menuItems : null}
		</Menu>
    </div>
  )
}

export default ContextMenu