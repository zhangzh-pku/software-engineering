import React from "react";
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import Fab from "@mui/material/Fab";
import AddIcon from "@mui/icons-material/Add";

interface CreateApplicationCardProps {
  onClick: () => void;
}

const CreateApplicationCard: React.FC<CreateApplicationCardProps> = ({
  onClick,
}) => {
  return (
    <Card onClick={onClick}>
      <CardActionArea>
        <Fab color="primary">
          <AddIcon />
        </Fab>
      </CardActionArea>
    </Card>
  );
};

export default CreateApplicationCard;
