import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { ContexMenuStates } from '../../interfaces/ContextMenu';

const initialState: ContexMenuStates = {
	mouseX: null,
	mouseY: null,
	itemId: undefined,
	itemType: undefined,
    treeType: undefined,
};

const contextMenuSlice = createSlice({
	name: 'contextMenu',
	initialState,
	reducers: {
		openContextMenu: (
			state,
			action: PayloadAction<{
				mouseX: number;
				mouseY: number;
				itemId: string;
				itemType: 'directory' | 'file';
                treeType: 'work' | 'archive';
			}>
		) => {
			state.mouseX = action.payload.mouseX;
			state.mouseY = action.payload.mouseY;
			state.itemId = action.payload.itemId;
			state.itemType = action.payload.itemType;
            state.treeType = action.payload.treeType;
		},
		closeContextMenu: state => {
			state.mouseX = null;
			state.mouseY = null;
		},
	},
});

export const { openContextMenu, closeContextMenu } = contextMenuSlice.actions;
export default contextMenuSlice.reducer;
