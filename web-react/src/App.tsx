import React from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { Box, Drawer, AppBar, Toolbar, Typography } from "@mui/material";
import Dashboard from "./components/Dashboard";

const drawerWidth = 240;

const App: React.FC = () => {
  return (
    <Router>
      <Box sx={{ display: "flex" }}>
        <AppBar
          position="fixed"
          sx={{
            width: `calc(100% - ${drawerWidth}px)`,
            marginLeft: `${drawerWidth}px`,
          }}
        >
          <Toolbar>
            <Typography variant="h6" noWrap component="div">
              Dashboard
            </Typography>
          </Toolbar>
        </AppBar>
        <Drawer
          variant="permanent"
          sx={{
            width: drawerWidth,
            flexShrink: 0,
            "& .MuiDrawer-paper": {
              width: drawerWidth,
              boxSizing: "border-box",
            },
          }}
        >
          <Toolbar />
          {/* You can add a list of menu items here */}
        </Drawer>
        <Box
          component="main"
          sx={{
            flexGrow: 1,
            p: 3,
            width: `calc(100% - ${drawerWidth}px)`,
            marginLeft: `${drawerWidth}px`,
          }}
        >
          <Toolbar />
          <Switch>
            <Route exact path="/" component={Dashboard} />
            {/* You can add more routes here */}
          </Switch>
        </Box>
      </Box>
    </Router>
  );
};

export default App;
