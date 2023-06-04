import React, { useState } from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import axios from "axios";

interface CreateApplicationProps {
  onSubmit: (
    script: string,
    doi: string,
    reproduction: string,
    dockerimage: string
  ) => void;
}

const CreateApplication: React.FC<CreateApplicationProps> = ({ onSubmit }) => {
  const [dockerimage, setDockerimage] = useState("");
  const [script, setScript] = useState("");
  const [doi, setDoi] = useState("");
  const [reproduction, setAppName] = useState("");

  const handleSubmit = () => {
    onSubmit(script, doi, reproduction, dockerimage);
    axios
      .post("http://localhost:8080/task/create", {
        dockerimage: "docker_image",
        script: "run_script",
        doi: "dissertation_doi",
      })
      .then((response) => {
        console.log(response);
      })
      .catch((error) => {
        console.log(error);
      });
    setDockerimage("");
    setScript("");
    setDoi("");
    setAppName("");
  };

  return (
    <div>
      <h2>Create Reproduction</h2>
      <TextField
        label="Dockerimage"
        value={dockerimage}
        onChange={(event) => setDockerimage(event.target.value)}
        fullWidth
        margin="normal"
      />
      <TextField
        label="Reproduction Name"
        value={reproduction}
        onChange={(event) => setAppName(event.target.value)}
        fullWidth
        margin="normal"
      />
      <TextField
        label="Script"
        value={script}
        onChange={(event) => setScript(event.target.value)}
        fullWidth
        margin="normal"
      />
      <TextField
        label="DOI of the Paper"
        value={doi}
        onChange={(event) => setDoi(event.target.value)}
        fullWidth
        margin="normal"
      />

      <Button variant="contained" color="primary" onClick={handleSubmit}>
        Submit
      </Button>
    </div>
  );
};

export default CreateApplication;
