import React, { useEffect, useState } from "react";
import axios from 'axios';
import { Application } from "../types";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import ApplicationBlock from "./ApplicationBlock"
import Button from "@mui/material/Button";
import Grid from '@mui/material/Grid';

interface ApplicationsListProps {
  applications: Application[];
}

const ApplicationsList: React.FC<ApplicationsListProps> = () => {
  const [applications, setApplications] = useState<Application[]>([]);

  useEffect(() => {
    const getApplications = async () => {
      try {
        const response = await axios.get('http://localhost:8080/reproductions');
        setApplications(response.data);
      } catch (error) {
        console.error(`Failed to fetch applications: ${error}`);
      }
    };

    getApplications();
  }, []);

  return (
    <div>
      <h1>Reproductions</h1>
      <List>
        {applications.map((application) => (
          <ListItem button key={application.id}>
            <ListItemText primary={application.name} />
          </ListItem>
        ))}
      </List>
      
      <Grid container spacing={3}>
        {applications.map((application) => (
          <ApplicationBlock application={application} />
        ))}
      </Grid>
    </div>
  );
};

export default ApplicationsList;
