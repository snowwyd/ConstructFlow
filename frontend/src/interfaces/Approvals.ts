/**
 * Interface for document approval response from API
 * Contains information about documents in the approval workflow
 */
export interface ApprovalResponse {
	// Unique identifier for the approval record
	approval_id: number;

	// Identifier for the file being approved
	file_id: number;

	// Name of the file being approved
	file_name: string;

	// Current status of the document in the approval process
	// Possible values: "on approval", "approved", "annotated", "finalized"
	status: string;

	// Current step in the approval workflow
	workflow_order: number;

	// Total number of steps in the approval workflow
	workflow_user_count: number;
}

/**
 * Interface for annotation request payload
 */
export interface AnnotationRequest {
	// Approval record identifier
	approvalId: number;

	// Annotation message/comment
	message: string;
}

/**
 * Interface for approval workflow status
 * Used for displaying status chips and workflow information
 */
export interface ApprovalStatus {
	// Display label for the status
	label: string;

	// Color for UI elements (chip, button, etc.)
	color:
		| 'default'
		| 'primary'
		| 'secondary'
		| 'error'
		| 'info'
		| 'success'
		| 'warning';

	// Icon identifier to display with the status
	icon: string;
}
