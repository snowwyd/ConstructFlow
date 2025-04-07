import { configureStore } from '@reduxjs/toolkit';
import contextMenuReducer from './Slices/contextMenuSlice';
import approvalsReducer from './Slices/contextMenuSlice';

export const store = configureStore({
	reducer: {
		contextMenu: contextMenuReducer,
		approvals: approvalsReducer,
	},
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
