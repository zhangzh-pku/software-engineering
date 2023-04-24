import React, { useState } from "react";
import "./App.css";
import { Application, Output } from "./types";
import ApplicationsList from "./components/ApplicationsList";
import CreateApplication from "./components/CreateApplication";
import ApplicationDetails from "./components/ApplicationDetails";

function App() {
  const [applications, setApplications] = useState<Application[]>([]);
  const [selectedApplication, setSelectedApplication] =
    useState<Application | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  const handleCreateApplication = () => {
    setShowCreateForm(true);
  };

  const handleSubmitApplication = (script: string, doi: string) => {
    const newApplication: Application = {
      id: applications.length + 1,
      name: `Application ${applications.length + 1}`,
      script,
      doi,
    };
    setApplications([...applications, newApplication]);
    setShowCreateForm(false);
  };

  const handleSelectApplication = (id: number) => {
    const app = applications.find((application) => application.id === id);
    setSelectedApplication(app ?? null);
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>My Application</h1>
      </header>
      <main>
        {!selectedApplication && !showCreateForm && (
          <ApplicationsList
            applications={applications}
            onCreateApplication={handleCreateApplication}
          />
        )}
        {!selectedApplication && showCreateForm && (
          <CreateApplication onSubmit={handleSubmitApplication} />
        )}
        {selectedApplication && (
          <ApplicationDetails
            application={selectedApplication}
            onBack={() => setSelectedApplication(null)}
          />
        )}
      </main>
    </div>
  );
}

export default App;
