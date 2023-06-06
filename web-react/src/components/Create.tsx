import * as React from 'react';
import { styled, createTheme, ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import MuiDrawer from '@mui/material/Drawer';
import Box from '@mui/material/Box';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Badge from '@mui/material/Badge';
import Container from '@mui/material/Container';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import Link from '@mui/material/Link';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import NotificationsIcon from '@mui/icons-material/Notifications';
import AccountInformation from './AccountInformation';
import CreateApplication from "./CreateApplication";
import { Application } from "../types";
import { useState } from "react";
import ApplicationsList from "./ApplicationsList";

interface createProps {
  applications: Application[]
  setApplications: any
  changeView: any
}

export default function CreateContent(param: createProps) {
  {/*const [applications, setApplications] = useState<Application[]>([]);*/ }

  const handleSubmitApplication = (script: string, doi: string) => {
    const newApplication: Application = {
      id: param.applications.length + 1,
      name: `Reproduction ${param.applications.length + 1}`,
      script,
      doi,
    };
    param.setApplications(newApplication);
    param.changeView("Display")
  };

  return (
    <Grid container spacing={3} justifyContent="center">
      {/* Create Application */}
      <Grid item xs={12} md={8} lg={9}>
        <Paper
          sx={{
            p: 2,
            display: 'flex',
            flexDirection: 'column',
            height: 480,
          }}
        >
          <CreateApplication onSubmit={handleSubmitApplication} />
        </Paper>
      </Grid>
    </Grid>
  )
}