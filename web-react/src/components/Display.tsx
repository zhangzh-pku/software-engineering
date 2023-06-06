import * as React from 'react';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import { Application } from "../types";
import ApplicationsList from "./ApplicationsList";

interface displayProps {
  applications: Application[]
}

export default function DisplayContent(param: displayProps) {

  return (
    <Grid container spacing={3}>
      {/* ApplicationList */}
      <Grid item xs={12}>
        <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
          <ApplicationsList
            applications={param.applications}
          />
        </Paper>
      </Grid>
    </Grid>
  )
}
