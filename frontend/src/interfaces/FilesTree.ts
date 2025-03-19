export interface File {
    id: number;
    name_file: string;
    status: string;
    directory_id: number;
  }
  
export interface Directory {
    id: number;
    name_folder: string;
    status: string;
    parent_path_id?: number | null;
    files: File[];
}
  
export interface TreeDataItem {
    id: string;
    label: string;
    status: string;
    type: "directory" | "file";
    children?: TreeDataItem[];
}