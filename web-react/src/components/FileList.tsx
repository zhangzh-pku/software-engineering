import * as React from 'react';
import Box from '@mui/material/Box';
import List from '@mui/material/List';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import FileDownloadEntry from './FileDownloadEntry';

// Is this file used?
function generate(thelist: string[]) {
    console.log("wocaonima")
    console.log(process.cwd())
    return thelist.map((value) => <FileDownloadEntry name={value} taskid = "0"/>);
}

const mypath = ".."

export default function InteractiveList() {
    const [dense, setDense] = React.useState(false);
    const [secondary, setSecondary] = React.useState(false);

    return (
        <Box sx={{ flexGrow: 1, maxWidth: 752 }}>
            <Grid container spacing={2}>
                <Grid item xs={12} md={6}>
                    <Typography sx={{ mt: 4, mb: 2 }} variant="h6" component="div">
                        File Display
                    </Typography>
                    <List dense={dense}>
                        {generate(["file1", "fuck", "sleepy", "time"])}
                    </List>
                </Grid>
            </Grid>
        </Box>
    );
}