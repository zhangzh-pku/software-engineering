import * as React from 'react';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import CreateApplication from "./CreateApplication";
import { Application } from "../types";

interface createProps {
  applications: Application[]
  setApplications: any
  changeView: any
}

export default function CreateContent(param: createProps) {

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