import React from "react";
import { Grid } from "@mui/material";
import { Application } from "../types";
import { useNavigate } from "react-router-dom";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import CardMedia from "@mui/material/CardMedia";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";

interface ApplicationBlockProps {
  application: Application;
}

function hash(str: string, length: number) {
  let sum = 0;
  for (let i = 0; i < str.length; i++) {
    sum += str.charCodeAt(i);
  }
  return sum % length;
}

const ApplicationBlock: React.FC<ApplicationBlockProps> = ({ application }) => {
  const navigate = useNavigate();

  const handleClickOnBlock = () => {
    console.log("I am clicked");
    navigate(`./${application.doi}`, { state: { application } });
  };

  console.log("application.id");
  console.log(application.id);
  const imagename: string = hash(application.id, 7).toString();
  console.log("imagename:");
  console.log(imagename);

  return (
    <Grid item xs={12} md={4} lg={3}>
      <Card sx={{ maxWidth: 345 }}>
        <CardMedia
          sx={{ height: 140 }}
          image={`./DecoImage/${imagename}.jpg`}
          title="green iguana"
        />
        <CardContent>
          <Typography gutterBottom variant="h5" component="div">
            {application.name}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Script:{application.script}
            <br />
            DOI:{application.doi}
          </Typography>
        </CardContent>
        <CardActions>
          <Button size="small" onClick={handleClickOnBlock}>
            Run
          </Button>
          <Button size="small">Learn More</Button>
        </CardActions>
      </Card>
    </Grid>
  );
};

export default ApplicationBlock;
