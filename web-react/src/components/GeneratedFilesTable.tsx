import React from "react";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";

interface File {
  name: string;
  size: number;
  type: string;
}

interface GeneratedFilesTableProps {
  files: File[];
}

const GeneratedFilesTable: React.FC<GeneratedFilesTableProps> = ({ files }) => {
  return (
    <TableContainer component={Paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Name</TableCell>
            <TableCell align="right">Size</TableCell>
            <TableCell align="right">Type</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {files.map((file, index) => (
            <TableRow key={index}>
              <TableCell component="th" scope="row">
                {file.name}
              </TableCell>
              <TableCell align="right">{file.size}</TableCell>
              <TableCell align="right">{file.type}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default GeneratedFilesTable;
