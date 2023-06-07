import * as React from "react";
import Box from "@mui/material/Box";
import Stepper from "@mui/material/Stepper";
import Step from "@mui/material/Step";
import StepLabel from "@mui/material/StepLabel";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import Grid from "@mui/material/Grid";
import Paper from "@mui/material/Paper";
import TextField from "@mui/material/TextField";
import { useState, ChangeEvent, useEffect } from "react";
import { useLocation } from "react-router-dom";
import axios from "axios";
import { Application } from "../types";
import Circular from "./CircularProgress";
import FileDownloadEntry from "./FileDownloadEntry";
const steps = ["Upload your dataset", "Run", "Download the output"];

export default function HorizontalLinearStepper() {
  const [activeStep, setActiveStep] = React.useState(0);
  const [skipped, setSkipped] = React.useState(new Set<number>());
  const location = useLocation();
  const application = location.state?.application;
  // 创建状态来存储任务ID和任务状态
  const [taskId, setTaskId] = useState<string | null>(null);
  const [taskStatus, setTaskStatus] = useState<string | null>("pending");
  const [fileList, setFileList] = useState<string[]>([]);

  const isStepSkipped = (step: number) => {
    return skipped.has(step);
  };

  const handleNext = () => {
    if (activeStep === 0) {
      console.log("application.id", application.id);
      const createResponse = axios
        .post("http://localhost:8080/run", {
          doid: application.id,
        })
        .then((response) => {
          setTaskId(response.data.taskID);
        });
    }

    let newSkipped = skipped;
    if (isStepSkipped(activeStep)) {
      newSkipped = new Set(newSkipped.values());
      newSkipped.delete(activeStep);
    }

    setActiveStep((prevActiveStep) => prevActiveStep + 1);
    setSkipped(newSkipped);
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleReset = () => {
    setActiveStep(0);
  };

  useEffect(() => {
    // 检查任务ID是否已经存在，如果不存在，则不需要检查任务状态
    if (taskId) {
      // 发送请求来获取任务状态
      const interval = setInterval(() => {
        axios
          .get(`http://localhost:8080/task_status/${taskId}`)
          .then((response) => {
            setTaskStatus(response.data);
          })
          .catch((error) => {
            console.error("Error fetching task status:", error);
          });
      }, 5000); // 每5秒钟获取一次状态

      // 清除定时器
      return () => clearInterval(interval);
    }
    if (taskStatus === "success") {
      axios
        .get(`http://localhost:8080/files/${taskId}`)
        .then((response) => {
          setFileList(response.data.files);
        })
        .catch((error) => {
          console.error("Error fetching file list:", error);
        });
    }
  }, [taskId, taskStatus]);

  const contentView = () => {
    switch (activeStep) {
      case 0:
        return (
          <Grid container spacing={3} justifyContent="center">
            <Grid item xs={12} md={8} lg={9}>
              <Paper
                sx={{
                  p: 2,
                  display: "flex",
                  flexDirection: "column",
                  height: 400,
                  alignItems: "center",
                  justifyContent: "center",
                  textAlign: "center",
                }}
              >
                <div>
                  <h3>Free dataset is saved in our system!</h3>
                </div>
              </Paper>
            </Grid>
          </Grid>
        );
      case 1:
        return (
          <Grid container spacing={3} justifyContent="center">
            <Grid item xs={12} md={8} lg={9}>
              <Paper
                sx={{
                  p: 2,
                  display: "flex",
                  flexDirection: "column",
                  height: 400,
                  alignItems: "center",
                  justifyContent: "center",
                  textAlign: "center",
                }}
              >
                <div>
                  <h3>Task Status: {taskStatus}</h3>
                  <Circular taskstatus={taskStatus}/>
                </div>
              </Paper>
            </Grid>
          </Grid>
        );
      case 2:
        return (
          <Grid container spacing={3} justifyContent="center">
            <Grid item xs={12} md={8} lg={9}>
              <Paper
                sx={{
                  p: 2,
                  display: "flex",
                  flexDirection: "column",
                  height: 400,
                  alignItems: "center",
                  justifyContent: "center",
                  textAlign: "center",
                }}
              >
                <div>
                  <h3>File List:</h3>
                  <ul>
                    {/*Modify here to map to the component FileDownloadEntry*/}
                    {taskId !== null && fileList.map((fileName) => (
                      <FileDownloadEntry name = {fileName} taskid={taskId} />
                    ))}
                  </ul>
                </div>
              </Paper>
            </Grid>
          </Grid>
        );
    }
  };

  return (
    <Box sx={{ width: "100%" }}>
      <Stepper activeStep={activeStep} alternativeLabel>
        {steps.map((label, index) => {
          const stepProps: { completed?: boolean } = {};
          const labelProps: {
            optional?: React.ReactNode;
          } = {};

          if (isStepSkipped(index)) {
            stepProps.completed = false;
          }
          return (
            <Step key={label} {...stepProps}>
              <StepLabel {...labelProps}>{label}</StepLabel>
            </Step>
          );
        })}
      </Stepper>
      {activeStep === steps.length ? (
        <React.Fragment>
          <Typography sx={{ mt: 2, mb: 1 }}>
            All steps completed - you&apos;re finished
          </Typography>
          <Box sx={{ display: "flex", flexDirection: "row", pt: 2 }}>
            <Box sx={{ flex: "1 1 auto" }} />
            <Button onClick={handleReset}>Reset</Button>
          </Box>
        </React.Fragment>
      ) : (
        <React.Fragment>
          {contentView()}
          <Box sx={{ display: "flex", flexDirection: "row", pt: 2 }}>
            <Button
              color="inherit"
              disabled={activeStep === 0}
              onClick={handleBack}
              sx={{ mr: 1 }}
            >
              Back
            </Button>
            <Box sx={{ flex: "1 1 auto" }} />
            <Button onClick={handleNext} disabled={activeStep === 1 && taskStatus !== "success"}>
              {activeStep === steps.length - 1 ? "Finish" : "Next"}
            </Button>
          </Box>
        </React.Fragment>
      )}
    </Box>
  );
}
