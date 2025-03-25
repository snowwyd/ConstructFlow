import { configureStore } from '@reduxjs/toolkit';
import contextMenuReducer from './Slices/contexMenuSlice';

export const store = configureStore({
	reducer: {
		contextMenu: contextMenuReducer,
	},
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
