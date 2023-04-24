import React from "react";
import { Application, Output } from "../types";

interface ApplicationDetailsProps {
  application: Application;
  onBack: () => void;
}

const ApplicationDetails: React.FC<ApplicationDetailsProps> = ({
  application,
  onBack,
}) => {
  // 模拟输出数据
  const outputs: Output[] = [
    { step: 1, content: "Output 1" },
    { step: 2, content: "Output 2" },
    { step: 3, content: "Output 3" },
  ];

  const downloadOutputs = () => {
    // 在这里实现下载功能，如将输出数据转换为文件并使用浏览器下载
  };

  return (
    <div>
      <button onClick={onBack}>Back to Applications</button>
      <h2>{application.name}</h2>
      <p>Script: {application.script}</p>
      <p>DOI: {application.doi}</p>
      <h3>Outputs</h3>
      <ul>
        {outputs.map((output) => (
          <li key={output.step}>
            Step {output.step}: {output.content}
          </li>
        ))}
      </ul>
      <button onClick={downloadOutputs}>Download Outputs</button>
    </div>
  );
};

export default ApplicationDetails;
