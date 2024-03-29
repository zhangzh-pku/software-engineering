import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import DashboardIcon from '@mui/icons-material/Dashboard';
import UploadIcon from '@mui/icons-material/Upload';
import PeopleIcon from '@mui/icons-material/People';
import BarChartIcon from '@mui/icons-material/BarChart';
import LayersIcon from '@mui/icons-material/Layers';

interface listProps {
  handleClick: (x: string) => void
}

export default function MainListItems(param: listProps) {
  return (
    (
      <React.Fragment>
        <ListItemButton onClick={() => param.handleClick('Display')}>
          <ListItemIcon>
            <DashboardIcon />
          </ListItemIcon>
          <ListItemText primary="Display" />
        </ListItemButton>
        <ListItemButton onClick={() => param.handleClick('Create')}>
          <ListItemIcon>
            <UploadIcon />
          </ListItemIcon>
          <ListItemText primary="Create" />
        </ListItemButton>
        <ListItemButton onClick={() => param.handleClick('Account')}>
          <ListItemIcon>
            <PeopleIcon />
          </ListItemIcon>
          <ListItemText primary="Account" />
        </ListItemButton>
        <ListItemButton>
          <ListItemIcon>
            <BarChartIcon />
          </ListItemIcon>
          <ListItemText primary="Reports" />
        </ListItemButton>
        <ListItemButton>
          <ListItemIcon>
            <LayersIcon />
          </ListItemIcon>
          <ListItemText primary="Integrations" />
        </ListItemButton>
      </React.Fragment>
    )
  )
} 
