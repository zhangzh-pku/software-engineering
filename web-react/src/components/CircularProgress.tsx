import * as React from 'react';
import Box from '@mui/material/Box';
import CircularProgress from '@mui/material/CircularProgress';
import { green, red } from '@mui/material/colors';
import Fab from '@mui/material/Fab';
import CheckIcon from '@mui/icons-material/Check';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import PriorityHighIcon from '@mui/icons-material/PriorityHigh';

interface CircularProps{
    taskstatus : string | null
}

export default function Circular(params : CircularProps) {

  const buttonSx = {
    ...(params.taskstatus === "success" && {
      bgcolor: green[500],
      '&:hover': {
        bgcolor: green[700],
      },
    } || params.taskstatus === "failed" && { bgcolor: red[500],
        '&:hover': {
          bgcolor: red[700],
        },}),
  };

  return (
    <Box sx={{ display: 'flex', alignItems: 'center' }}>
      <Box sx={{ m: 1, position: 'relative' }}>
        <Fab
          aria-label="save"
          color="primary"
          sx={buttonSx}
        >
          {params.taskstatus === "success" ? <CheckIcon /> : (params.taskstatus === "failed" ? <PriorityHighIcon/> : <PlayArrowIcon />)}
        </Fab>
        {params.taskstatus === "running" && (
          <CircularProgress
            size={68}
            sx={{
              color: green[500],
              position: 'absolute',
              top: -6,
              left: -6,
              zIndex: 1,
            }}
          />
        )}
      </Box>
    </Box>
  );
}