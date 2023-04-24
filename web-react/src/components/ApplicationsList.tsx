import React from "react";
import { Application } from "../types";

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
      <ul>
        {applications.map((app) => (
          <li key={app.id}>
            <a href={`/applications/${app.id}`}>{app.name}</a>
          </li>
        ))}
      </ul>
      <button onClick={onCreateApplication}>Create Application</button>
    </div>
  );
};

export default ApplicationsList;
