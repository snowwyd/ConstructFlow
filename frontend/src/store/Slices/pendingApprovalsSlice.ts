import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface ApprovalsState {
	pendingCount: number;
}

const initialState: ApprovalsState = {
	pendingCount: 0,
};

const approvalsSlice = createSlice({
	name: 'approvals',
	initialState,
	reducers: {
		setPendingCount: (state, action: PayloadAction<number>) => {
			state.pendingCount = action.payload;
		},
	},
});

export const { setPendingCount } = approvalsSlice.actions;
export default approvalsSlice.reducer;
