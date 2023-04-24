import React, { useState } from "react";

interface CreateApplicationProps {
  onSubmit: (script: string, doi: string) => void;
}

const CreateApplication: React.FC<CreateApplicationProps> = ({ onSubmit }) => {
  const [script, setScript] = useState("");
  const [doi, setDoi] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(script, doi);
  };

  return (
    <div>
      <h2>Create Application</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Script:</label>
          <textarea
            value={script}
            onChange={(e) => setScript(e.target.value)}
          />
        </div>
        <div>
          <label>DOI:</label>
          <input
            type="text"
            value={doi}
            onChange={(e) => setDoi(e.target.value)}
          />
        </div>
        <button type="submit">Submit</button>
      </form>
    </div>
  );
};

export default CreateApplication;
