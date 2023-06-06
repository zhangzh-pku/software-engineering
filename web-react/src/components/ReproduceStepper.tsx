import * as React from 'react';
import Box from '@mui/material/Box';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import TextField from "@mui/material/TextField";
import { useState, ChangeEvent, useEffect } from "react";
import Linklist from "./Linklist"
import InteractiveList from './FileList';

const steps = ['Upload your dataset', 'Run', 'Download the output'];
{/* Attention that here we must use addresses with http://! 
Otherwise it will go to localhost:3000/path instead of it directly!*/}
const links = ["http://www.baidu.com", "http://www.taobao.com"];

interface FileDisplayProps {
  path: string
}

const FileInput = () => {
  const [fileContent, setFileContent] = useState('');

  const handleFileRead = (e: any) => {
    const content = e.target.result;
    setFileContent(content);
  };

  const handleFileChosen = (file: any) => {
    const fileReader = new FileReader();
    fileReader.onloadend = handleFileRead;
    fileReader.readAsText(file[0]);
  };

  return (
    <div>
      <input type="file" onChange={(e) => handleFileChosen(e.target.files)} />
      <p>{fileContent}</p>
    </div>
  );
};



{/*This component is used to display log content dynamically 
  Now the path is /public/try.txt
*/ }
function FileDisplay({ path }: FileDisplayProps) {
  const [fileContents, setFileContents] = useState('');

  const fetchFileContents = async () => {
    try {
      const response = await fetch(path);
      const content = await response.text();
      setFileContents(content);
    } catch (error) {
      console.error('Error loading file:', error);
    }
  };

  const longPollingFetch = () => {
    fetchFileContents(); // 发起获取文件内容的请求

    setTimeout(longPollingFetch, 1000); // 1秒后再次进行长轮询
  };

  useEffect(() => {
    longPollingFetch(); // 开始进行长轮询
  }, [path]);

  return (
    <div>
      <pre>{fileContents}</pre>
    </div>
  );
}

export default function HorizontalLinearStepper() {
  const [activeStep, setActiveStep] = React.useState(0);
  const [skipped, setSkipped] = React.useState(new Set<number>());


  const isStepOptional = (step: number) => {
    return step === 5;
  };

  const isStepSkipped = (step: number) => {
    return skipped.has(step);
  };

  const handleNext = () => {
    let newSkipped = skipped;
    {console.log(process.cwd())}
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

  const handleSkip = () => {
    if (!isStepOptional(activeStep)) {
      // You probably want to guard against something like this,
      // it should never occur unless someone's actively trying to break something.
      throw new Error("You can't skip a step that isn't optional.");
    }

    setActiveStep((prevActiveStep) => prevActiveStep + 1);
    setSkipped((prevSkipped) => {
      const newSkipped = new Set(prevSkipped.values());
      newSkipped.add(activeStep);
      return newSkipped;
    });
  };

  const handleReset = () => {
    setActiveStep(0);
  };

  {/*For step 0 */ }
  const [zippath, setZippath] = useState("")

  {/*For step 1 */ }
  const [taskstatus, setTaskstatus] = React.useState("waiting"); {/* 0 for waiting, 1 for running, 2 for exited, 3 for dead*/ }

  const handleRun = () => {
    setTaskstatus("running");
  };


  const contentView = () => {
    switch (activeStep) {
      case 0:
        return (<Grid container spacing={3} justifyContent="center">
          {/* Upload Dataset */}
          <Grid item xs={12} md={8} lg={9}>
            <Paper
              sx={{
                p: 2,
                display: 'flex',
                flexDirection: 'column',
                height: 400,
              }}
            >
              <div>
                <h2>Upload Dataset Zip</h2>
                <TextField
                  label="Path of zip"
                  value={zippath}
                  onChange={(event) => setZippath(event.target.value)}
                  fullWidth
                  margin="normal"
                />
                <h3>Maybe you will be using following datasets:</h3>
                <Linklist links={links} />
              </div>
            </Paper>
          </Grid>
        </Grid>)

      case 1: return (<Grid container spacing={3} justifyContent="center">
        {/* Run the image */}
        <Grid item xs={12} md={8} lg={9}>
          <Paper
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              height: 200,
            }}
          >
            <div>
              <h2>Now the zip of dataset you've chosen is {zippath}</h2>
              <h2>Status of the task: {taskstatus} </h2>
              <Button
                variant="contained"
                color="secondary"
                disabled={taskstatus !== "waiting"}
                sx={{ mr: 1 }}
                onClick={handleRun}
              >
                Run
              </Button>
            </div>
          </Paper>
          <Paper
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              height: 200,
            }}
          >
            <FileDisplay path={"../try.txt"} />
          </Paper>
        </Grid>
      </Grid>)

      case 2: return(<Grid container spacing={3} justifyContent="center">
      {/* Run the image */}
      <Grid item xs={12} md={8} lg={9}>
        <Paper
          sx={{
            p: 2,
            display: 'flex',
            flexDirection: 'column',
            height: 400,
            maxHeight: "100%",
            overflow:"auto"
          }}
        >
          {/*<FileInput />*/}
          <InteractiveList/>
        </Paper>
      </Grid>
    </Grid>)
    }
  }

  return (
    <Box sx={{ width: '100%' }}>
      <Stepper activeStep={activeStep} alternativeLabel>
        {steps.map((label, index) => {
          const stepProps: { completed?: boolean } = {};
          const labelProps: {
            optional?: React.ReactNode;
          } = {};
          if (isStepOptional(index)) {
            labelProps.optional = (
              <Typography variant="caption">Optional</Typography>
            );
          }
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
          <Box sx={{ display: 'flex', flexDirection: 'row', pt: 2 }}>
            <Box sx={{ flex: '1 1 auto' }} />
            <Button onClick={handleReset}>Reset</Button>
          </Box>
        </React.Fragment>
      ) : (
        <React.Fragment>
          {contentView()}
          <Box sx={{ display: 'flex', flexDirection: 'row', pt: 2 }}>
            <Button
              color="inherit"
              disabled={activeStep === 0}
              onClick={handleBack}
              sx={{ mr: 1 }}
            >
              Back
            </Button>
            <Box sx={{ flex: '1 1 auto' }} />
            {isStepOptional(activeStep) && (
              <Button color="inherit" onClick={handleSkip} sx={{ mr: 1 }}>
                Skip
              </Button>
            )}
            {/*|| (activeStep === 1 && taskstatus !== "finish")
              Add this to the condition of disabled*/}
            <Button onClick={handleNext} disabled = {(activeStep === 0 && zippath === '') }>
              {activeStep === steps.length - 1 ? 'Finish' : 'Next'}
            </Button>
          </Box>
        </React.Fragment>
      )}
    </Box>
  );
}