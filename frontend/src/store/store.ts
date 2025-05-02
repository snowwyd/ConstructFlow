import { configureStore } from '@reduxjs/toolkit';
import {
	default as approvalsReducer,
	default as contextMenuReducer,
} from './Slices/contextMenuSlice';

export const store = configureStore({
	reducer: {
		contextMenu: contextMenuReducer,
		approvals: approvalsReducer,
	},
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
