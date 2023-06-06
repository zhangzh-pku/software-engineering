import React from "react";
import { Application, Output } from "../types";
import Button from "@mui/material/Button";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";

//Is this file used?

interface ApplicationDetailsProps {
  application: Application;
  onBack: () => void;
}

const ApplicationDetails: React.FC<ApplicationDetailsProps> = ({
  application,
  onBack,
}) => {
  const outputs: Output[] = [
    { step: 1, content: "Output 1" },
    { step: 2, content: "Output 2" },
    { step: 3, content: "Output 3" },
  ];

  const downloadOutputs = () => {
    // 在这里实现下载功能，如将输出数据转换为文件并使用浏览器下载
  };

  return (
    <Container>
      <Button variant="outlined" onClick={onBack} sx={{ my: 2 }}>
        Back to Applications
      </Button>
      <Typography variant="h4" component="h2">
        {application.name}
      </Typography>
      <Typography variant="body1">Script: {application.script}</Typography>
      <Typography variant="body1">DOI: {application.doi}</Typography>
      <Typography variant="h6" component="h3" sx={{ mt: 2 }}>
        Outputs
      </Typography>
      <List>
        {outputs.map((output) => (
          <ListItem key={output.step}>
            <ListItemText primary={`Step ${output.step}: ${output.content}`} />
          </ListItem>
        ))}
      </List>
      <Button
        variant="contained"
        color="primary"
        onClick={downloadOutputs}
        sx={{ mt: 2 }}
      >
        Download Outputs
      </Button>
    </Container>
  );
};

export default ApplicationDetails;
