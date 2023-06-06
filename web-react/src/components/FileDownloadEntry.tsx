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
}


const handleClick = async (name: string) => {
    const response = await fetch('http://localhost:8080/' + name);
    const blob = await response.blob();
    const url = window.URL.createObjectURL(new Blob([blob]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', 'file.pdf');
    document.body.appendChild(link);
    link.click();
};




export default function FileDownloadEntry(params: FileDownloadEntryProps) {
    return (
        <ListItem
            secondaryAction={
                <IconButton edge="end" aria-label="delete">
                    <DownloadIcon onClick={() => handleClick(params.name)} />
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
