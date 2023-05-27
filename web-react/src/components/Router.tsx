import React from 'react';
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import Dashboard from './Dashboard';
import ReproducePage from './ReproducePage';
import { AppProvider } from './AppContext';

const AppRouter: React.FC = () => {
  return (
    <Router>
      <AppProvider>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/:reproductionId" element = {<ReproducePage/>}/>
      </Routes>
      </AppProvider>
    </Router>
  );
};

export default AppRouter;