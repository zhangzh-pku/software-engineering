import React, { createContext, useState, PropsWithChildren } from 'react';
import { Application } from '../types';

interface AppContextProps {
    uploadedApplications: Application[];
    addUploadedApplication: (application: Application) => void;
  }
  
  export const AppContext = createContext<AppContextProps>({
    uploadedApplications: [],
    addUploadedApplication: () => {},
  });
  
  export const AppProvider: React.FC<PropsWithChildren<{}>> = ({ children }) => {
    const [uploadedApplications, setUploadedApplications] = useState<Application[]>([]);
  
    const addUploadedApplication = (application: Application) => {
      setUploadedApplications([...uploadedApplications, application]);
    };
  
    return (
      <AppContext.Provider value={{ uploadedApplications, addUploadedApplication }}>
        {children}
      </AppContext.Provider>
    );
  };