import React, { useState } from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";

interface CreateApplicationProps {
  onSubmit: (script: string, doi: string) => void;
}

const CreateApplication: React.FC<CreateApplicationProps> = ({ onSubmit }) => {
  const [dockerimage, setDockerimage] = useState("")
  const [script, setScript] = useState("");
  const [doi, setDoi] = useState("");

  const handleSubmit = () => {
    onSubmit(script, doi);
    setDockerimage("");
    setScript("");
    setDoi("");
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
