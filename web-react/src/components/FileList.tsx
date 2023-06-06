import * as React from 'react';
import { styled } from '@mui/material/styles';
import Box from '@mui/material/Box';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Avatar from '@mui/material/Avatar';
import IconButton from '@mui/material/IconButton';
import FormGroup from '@mui/material/FormGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import FolderIcon from '@mui/icons-material/Folder';
import DownloadIcon from '@mui/icons-material/Download';
import * as BrowserFS from "browserfs"
import FileDownloadEntry from './FileDownloadEntry';

BrowserFS.install(window);

const fs = BrowserFS.BFSRequire("fs")

var path = require('path-browserify')

function getFilesInPath(dirPath: string): { [key: string]: string } {
    console.log(dirPath)
    console.log("This is current " + __dirname)
    const files = fs.readdirSync(dirPath);

    return files.reduce((acc: any, file: any) => {
        const filePath = path.join(dirPath, file);
        const fileStat = fs.statSync(filePath);

        if (fileStat.isFile()) {
            acc[file] = filePath;
        }

        {
            console.log("wocaonima")
            console.log(acc)
        }
        return acc;
    }, {});
}



function generate(thelist: string[]) {
    console.log("wocaonima")
    console.log(process.cwd())
    return thelist.map((value) => <FileDownloadEntry name = {value}/>);
}

const mypath = ".."

export default function InteractiveList() {
    const [dense, setDense] = React.useState(false);
    const [secondary, setSecondary] = React.useState(false);

    BrowserFS.initialize(new BrowserFS.FileSystem.LocalStorage());

    return (
        <Box sx={{ flexGrow: 1, maxWidth: 752 }}>
            <Grid container spacing={2}>
                <Grid item xs={12} md={6}>
                    <Typography sx={{ mt: 4, mb: 2 }} variant="h6" component="div">
                        File Display
                    </Typography>
                    <List dense={dense}>
                        {generate(["file1","fuck","sleepy","time"])}
                    </List>
                </Grid>
            </Grid>
        </Box>
    );
}