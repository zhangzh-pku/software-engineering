import * as React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import MenuIcon from "@mui/icons-material/Menu";
import ApplicationsList from "./ApplicationsList";
import CreateApplication from "./CreateApplication";
import ApplicationDetails from "./ApplicationDetails";
import { Application } from "../types";
import { useState } from "react";

export default function Dashboard() {
  const [applications, setApplications] = useState<Application[]>([]);
  const [selectedApplication, setSelectedApplication] =
    useState<Application | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  const handleCreateApplication = () => {
    setShowCreateForm(true);
  };

  const handleSubmitApplication = (script: string, doi: string) => {
    const newApplication: Application = {
      id: applications.length + 1,
      name: `Application ${applications.length + 1}`,
      script,
      doi,
    };
    setApplications([...applications, newApplication]);
    setShowCreateForm(false);
  };

  const handleSelectApplication = (id: number) => {
    const app = applications.find((application) => application.id === id);
    setSelectedApplication(app ?? null);
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <IconButton
            size="large"
            edge="start"
            color="inherit"
            aria-label="open drawer"
            sx={{ mr: 2 }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div">
            My Application
          </Typography>
        </Toolbar>
      </AppBar>
      <Box sx={{ p: 3 }}>
        {!selectedApplication && !showCreateForm && (
          <ApplicationsList
            applications={applications}
            onCreateApplication={handleCreateApplication}
          />
        )}
        {!selectedApplication && showCreateForm && (
          <CreateApplication onSubmit={handleSubmitApplication} />
        )}
        {selectedApplication && (
          <ApplicationDetails
            application={selectedApplication}
            onBack={() => setSelectedApplication(null)}
          />
        )}
      </Box>
    </Box>
  );
}
