import React from "react";
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import Typography from "@mui/material/Typography";

interface Application {
  title: string;
}

interface ApplicationCardProps {
  application: Application;
  onClick: () => void;
}

const ApplicationCard: React.FC<ApplicationCardProps> = ({
  application,
  onClick,
}) => {
  return (
    <Card onClick={onClick}>
      <CardActionArea>
        <Typography variant="h5">{application.title}</Typography>
      </CardActionArea>
    </Card>
  );
};

export default ApplicationCard;
