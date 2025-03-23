import {Button, Dialog, DialogActions, DialogContent, DialogTitle, Menu, MenuItem, TextField} from '@mui/material';
import {useDispatch, useSelector} from "react-redux";
import { closeContextMenu } from '../store/Slices/contexMenuSlice';
import config from "../constants/Configurations.json";
import { useRef, useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import axiosFetching from '../api/AxiosFetch';
import { AxiosError } from 'axios';


const createFolder =config.createDirectory;

const ContextMenu = () => {
    const dispatch = useDispatch();
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [newName, setNewName] = useState("");
    const treeRef = useRef<HTMLDivElement>(null);
    const {mouseX, mouseY, itemId, itemType} = useSelector((state: any) => state.contextMenu);

    const createFolderQuery = useMutation({
        mutationFn: async (data: {parent_path_id: number, name: string}) => {
            const response = await axiosFetching.post(createFolder, data);
            return response.data;
        },
        onSuccess: () => {
            console.log("Folder created");
            setIsDialogOpen(false);
            setNewName("");
        },
        onError: (error: AxiosError) => {
            console.error("Error: ", error);
        }
    });

    const handleCreateFolderSubmit = () => {
        if (!itemId || !newName.trim()) return;
        const parentPathId = parseInt(itemId.replace("dir-", ""), 10); 
        createFolderQuery.mutate({ parent_path_id: parentPathId, name: newName });
      };
    
    const handleCloseMenu = () => {
        dispatch(closeContextMenu());
    };

    const handleCreateFolder = () => {
        setIsDialogOpen(true);
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
    <>
        <Menu
		open={mouseX !== null && mouseY !== null}
		onClose={handleCloseMenu}
		anchorReference='anchorPosition'
		anchorPosition={mouseY !== null && mouseX !== null ? { top: mouseY, left: mouseX } : undefined}
		>
            {menuItems.length > 0 ? menuItems : null}
		</Menu>

        <Dialog open={isDialogOpen} onClose={() => {
            setIsDialogOpen(false);
            treeRef.current?.focus();
            }} >
        <DialogTitle>Создание новой папки</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Имя папки"
            fullWidth
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setIsDialogOpen(false)}>Отмена</Button>
          <Button onClick={handleCreateFolderSubmit} disabled={!newName.trim()}>
            Создать
          </Button>
        </DialogActions>
      </Dialog>
    </>
  )
}

export default ContextMenu