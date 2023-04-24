import React from "react";
import { Application } from "../types";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import Button from "@mui/material/Button";

interface ApplicationsListProps {
  applications: Application[];
  onCreateApplication: () => void;
}

const ApplicationsList: React.FC<ApplicationsListProps> = ({
  applications,
  onCreateApplication,
}) => {
  return (
    <div>
      <h2>Applications</h2>
      <List>
        {applications.map((application) => (
          <ListItem button key={application.id}>
            <ListItemText primary={application.name} />
          </ListItem>
        ))}
      </List>
      <Button variant="contained" color="primary" onClick={onCreateApplication}>
        Create Application
      </Button>
    </div>
  );
};

export default ApplicationsList;
