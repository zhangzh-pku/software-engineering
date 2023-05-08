import React, { useState } from "react";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogActions from "@mui/material/DialogActions";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";

interface CreateApplicationDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (doi: string, script: string) => void;
}

const CreateApplicationDialog: React.FC<CreateApplicationDialogProps> = ({
  open,
  onClose,
  onSubmit,
}) => {
  const [doi, setDoi] = useState("");
  const [script, setScript] = useState("");

  const handleDoiChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setDoi(event.target.value);
  };

  const handleScriptChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setScript(event.target.value);
  };

  const handleSubmit = () => {
    onSubmit(doi, script);
    setDoi("");
    setScript("");
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>Create Application</DialogTitle>
      <DialogContent>
        <TextField
          label="DOI"
          value={doi}
          onChange={handleDoiChange}
          fullWidth
        />
        <TextField
          label="Script"
          multiline
          value={script}
          onChange={handleScriptChange}
          fullWidth
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button onClick={handleSubmit}>Create</Button>
      </DialogActions>
    </Dialog>
  );
};

export default CreateApplicationDialog;
