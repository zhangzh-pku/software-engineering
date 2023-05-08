import React, { useState } from "react";
import Grid from "@mui/material/Grid";
import Container from "@mui/material/Container";
import ApplicationCard from "./ApplicationCard";
import CreateApplicationCard from "./CreateApplicationCard";
import CreateApplicationDialog from "./CreateApplicationDialog";

interface Application {
  title: string;
}

const Dashboard: React.FC = () => {
  const [applications, setApplications] = useState<Application[]>([]);
  const [dialogOpen, setDialogOpen] = useState(false);

  const handleDialogOpen = () => {
    setDialogOpen(true);
  };

  const handleDialogClose = () => {
    setDialogOpen(false);
  };

  const handleCreateApplication = (doi: string, script: string) => {
    const newApplication: Application = {
      title: `${doi} - ${script}`,
    };

    setApplications([...applications, newApplication]);
  };

  return (
    <Container>
      <Grid container spacing={3}>
        {applications.map((application, index) => (
          <Grid item key={index}>
            <ApplicationCard application={application} onClick={() => {}} />
          </Grid>
        ))}
        <Grid item>
          <CreateApplicationCard onClick={handleDialogOpen} />
        </Grid>
      </Grid>
      <CreateApplicationDialog
        open={dialogOpen}
        onClose={handleDialogClose}
        onSubmit={handleCreateApplication}
      />
    </Container>
  );
};

export default Dashboard;
