import React, { useState } from "react";
import Typography from '@mui/material/Typography'
import TextField from "@mui/material/TextField";
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import { ButtonGroup, Button } from "@mui/material";

export default function AccountInformation() {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [loggedIn, setLoggedIn] = useState(false);

	const handleLogin = () => {
		{/* Send to group 1 and verify */ }
		setLoggedIn(true);
	};

	const handleLogout = () => {
		setLoggedIn(false);
	};

	return (
		<Grid container spacing={3} justifyContent="center">
			<Grid item xs={12} md={8} lg={9}>
				<Paper sx={{ p: 2, display: 'flex', flexDirection: 'column', height: 400 }}>
					<div>
						<h1>Your Account</h1>
						<TextField label="UserName" value={username} onChange={(event) => setUsername(event.target.value)} fullWidth margin="normal" disabled={loggedIn} />
						<TextField label="Password" value={password} onChange={(event) => setPassword(event.target.value)} fullWidth margin="normal" disabled={loggedIn} />
						<Grid container justifyContent="center" sx={{ mt: 4 }}>
							<ButtonGroup variant="contained" size="large">
								{loggedIn ? (
									<>
										<Button color="success" disabled>
											Log in
										</Button>
										<Button color="error" onClick={handleLogout}>
											Log out
										</Button>
									</>
								) : (
									<>
										<Button color="success" onClick={handleLogin}>
											Log in
										</Button>
										<Button color="error" disabled>
											Log out
										</Button>
									</>
								)}
							</ButtonGroup>
						</Grid>
					</div>
				</Paper>
			</Grid>
		</Grid>
	);
}