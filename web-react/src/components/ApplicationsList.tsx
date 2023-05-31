import React from "react";
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

const ApplicationsList: React.FC<ApplicationsListProps> = ({
  applications,
}) => {
  return (
    <div>
      <h1>Reproductions</h1>
      {/*
      <List>
        {applications.map((application) => (
          <ListItem button key={application.id}>
            <ListItemText primary={application.name} />
          </ListItem>
        ))}
      </List>
      */}
      
      {
      <Grid container spacing = {3}>
        {applications.map((application) => (
          <ApplicationBlock application = {application} />
        ))}
      </Grid>
      } 
    </div>
  );
};

export default ApplicationsList;
