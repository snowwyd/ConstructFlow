export interface ContexMenuStates {
	mouseX: number | null;
	mouseY: number | null;
	itemId?: string;
	itemType?: 'directory' | 'file';
	treeType?: 'work' | 'archive';
}
