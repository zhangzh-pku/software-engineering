import React, { useState } from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import axios from "axios";

interface CreateApplicationProps {
  onSubmit: (
    script: string,
    doi: string,
    reproduction: string,
    dockerimage: string,
    filePath: string
  ) => void;
}

const CreateApplication: React.FC<CreateApplicationProps> = ({ onSubmit }) => {
  const [dockerimage, setDockerimage] = useState("");
  const [script, setScript] = useState("");
  const [doi, setDoi] = useState("");
  const [reproduction, setAppName] = useState("");
  const [file, setFile] = useState<File | null>(null);


  const handleSubmit = async () => {
    // 先上传文件
    if (file) {
      const formData = new FormData();
      formData.append('file', file);

      try {
        const uploadResponse = await axios.post('http://localhost:8080/upload', formData);
        const filePath = uploadResponse.data.message;  // 假设返回的路径在data属性中

        // 使用路径提交其他数据
        onSubmit(script, doi, reproduction, dockerimage, filePath);

        console.log("docker image", dockerimage);
        console.log("run script", script);
        console.log("doi", doi);
        console.log("path", filePath)
        const createResponse = await axios.post("http://localhost:8080/task/create", {
          docker_image: dockerimage,
          run_script: script,
          dissertation_doi: doi,
          data_path: filePath // 假设服务器期望的参数名是file_path
        });

        console.log(createResponse);
      } catch (error) {
        console.log(error);
      }

      setDockerimage("");
      setScript("");
      setDoi("");
      setAppName("");
      setFile(null);
    }
  };
  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files ? event.target.files[0] : null;
    setFile(file);
  };

  return (
    <div>
      <h2>Create Reproduction</h2>
      <TextField
        label="Docker Image"
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
      <TextField
        type="file"
        onChange={handleFileChange}
        fullWidth
        margin="normal"
      />
      <Button
        variant="contained"
        color="primary"
        onClick={handleSubmit}
        disabled={
          dockerimage === "" ||
          doi === "" ||
          reproduction === "" ||
          script === ""
        }
      >
        Submit
      </Button>
    </div>
  );
};

export default CreateApplication;
