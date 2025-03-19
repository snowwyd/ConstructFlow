import React from "react";
import { RichTreeView, TreeItem2, TreeItem2Props } from "@mui/x-tree-view";
import FolderIcon from "@mui/icons-material/Folder";
import DescriptionIcon from "@mui/icons-material/Description";
import ArchiveIcon from "@mui/icons-material/Archive";
import { Box, styled } from "@mui/material";
import { Directory } from "../interfaces/FilesTree";
import { TreeDataItem } from "../interfaces/FilesTree";


//ТЕСТОВЫЕ ДАННЫЕ УДАЛИТЬ
const apiResponse = {
  data: [
    {
      id: 1,
      name_folder: "ROOT",
      status: "archive",
      files: [
        {
          id: 1,
          name_file: "Archived1.txt",
          status: "archive",
          directory_id: 1,
        },
      ],
    },
    {
      id: 2,
      name_folder: "Archived Directory",
      status: "archive",
      parent_path_id: 1,
      files: [
        {
          id: 2,
          name_file: "Archived2.txt",
          status: "archive",
          directory_id: 2,
        },
      ],
    },
  ],
};
const CustomTreeItem = styled(TreeItem2)(({ theme }) => ({
  "& .MuiTreeItem-content": {
    padding: theme.spacing(0.5, 0),
  },
  "& .MuiTreeItem-label": {
    display: "flex",
    alignItems: "center",
    gap: theme.spacing(1),
  },
}));

const transformDataToTreeItems = (data: Directory[]): TreeDataItem[] => {
  const map = new Map<number, TreeDataItem>();

  data.forEach((item) => {
    map.set(item.id, {
      id: `dir-${item.id}`,
      label: item.name_folder,
      status: item.status,
      type: "directory",
      children: [],
    });
  });


  data.forEach((item) => {
    const node = map.get(item.id)!;

    item.files.forEach((file) => {
      node.children!.push({
        id: `file-${file.id}`,
        label: file.name_file,
        status: file.status,
        type: "file",
      });
    });

    if (item.parent_path_id) {
      const parent = map.get(item.parent_path_id);
      if (parent) parent.children!.push(node);
    }
  });

  return data
    .filter((item) => !item.parent_path_id)
    .map((item) => map.get(item.id)!);
};

const FilesTree: React.FC = () => {
  const treeItems = transformDataToTreeItems(apiResponse.data);

  return (
    <RichTreeView
      items={treeItems}
      defaultExpandedItems={["dir-1"]}
      slots={{
        item: (props: TreeItem2Props) => {
          const { itemId, label, ...rest } = props;

          const findItem = (items: TreeDataItem[], id: string): TreeDataItem | undefined => {
            for (const item of items) {
              if (item.id === id) return item;
              if (item.children) {
                const found = findItem(item.children, id);
                if (found) return found;
              }
            }
            return undefined;
          };

          const itemData = findItem(treeItems, itemId!);

          if (!itemData) {
            return <TreeItem2 {...rest} itemId={itemId} label={label} />;
          }

          return (
            <CustomTreeItem
              {...rest}
              itemId={itemId}
              label={
                <Box display="flex" alignItems="center" gap={1}>
                  {itemData.type === "directory" ? (
                    itemData.status === "archive" ? (
                      <ArchiveIcon color="error" />
                    ) : (
                      <FolderIcon color="primary" />
                    )
                  ) : (
                    <DescriptionIcon color="secondary" />
                  )}
                  <span>
                    {label} ({itemData.status})
                  </span>
                </Box>
              }
            />
          );
        },
      }}
      sx={{
        width: "100%",
        maxWidth: 400,
        bgcolor: "background.paper",
        border: "1px solid #ccc",
        borderRadius: 1,
        padding: 1,
      }}
    />
  );
};

export default FilesTree;