import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { ContextMenuStates } from '../../interfaces/ContextMenu';

const initialState: ContextMenuStates = {
	mouseX: null,
	mouseY: null,
	itemId: undefined,
	itemType: undefined,
	treeType: undefined,
};

// Добавляем новое действие для обновления позиции меню вместо закрытия и повторного открытия
const contextMenuSlice = createSlice({
	name: 'contextMenu',
	initialState,
	reducers: {
		openContextMenu: (state, action) => {
			// Обновляем все поля за один шаг
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
		// Новое действие для обновления позиции
		updateContextMenuPosition: (
			state,
			action: PayloadAction<{ mouseX: number; mouseY: number }>
		) => {
			state.mouseX = action.payload.mouseX;
			state.mouseY = action.payload.mouseY;
		},
	},
});

export const { openContextMenu, closeContextMenu, updateContextMenuPosition } =
	contextMenuSlice.actions;

export default contextMenuSlice.reducer;
