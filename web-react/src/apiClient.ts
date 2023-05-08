import axios, { AxiosInstance } from 'axios';

const apiClient: AxiosInstance = axios.create({
  baseURL: 'http://localhost:8080', // 替换为您的后端地址
});

export default apiClient;
