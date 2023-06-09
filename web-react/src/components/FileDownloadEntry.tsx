import * as React from 'react';
import ListItem from '@mui/material/ListItem';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import ListItemText from '@mui/material/ListItemText';
import Avatar from '@mui/material/Avatar';
import IconButton from '@mui/material/IconButton';
import FolderIcon from '@mui/icons-material/Folder';
import DownloadIcon from '@mui/icons-material/Download';

interface FileDownloadEntryProps {
    name: string
    taskid : string
}


const handleClick = async (name: string, taskid : string) => {
    const response = await fetch('http://localhost:8080/file/' + taskid + '/' + name);
    const blob = await response.blob();

    // 注意这里要判断 HTTP 状态码，如果不是 200，说明可能出现错误。
    if (!response.ok) {
        console.error("Error downloading file:", response.status, response.statusText);
        return;
    }

    const url = window.URL.createObjectURL(new Blob([blob]));
    const link = document.createElement('a');
    link.href = url;

    // 这里使用实际的文件名作为下载文件的名字。
    link.setAttribute('download', name);
    document.body.appendChild(link);
    link.click();

    // 记得在完成后删除这个元素。
    document.body.removeChild(link);
};





export default function FileDownloadEntry(params: FileDownloadEntryProps) {
    return (
        <ListItem
            secondaryAction={
                <IconButton edge="end" aria-label="delete">
                    <DownloadIcon onClick={() => handleClick(params.name, params.taskid)} />
                </IconButton>
            }
        >
            <ListItemAvatar>
                <Avatar>
                    <FolderIcon />
                </Avatar>
            </ListItemAvatar>
            <ListItemText
                primary={params.name}
            />
        </ListItem>
    );
}
