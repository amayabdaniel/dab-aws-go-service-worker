import axios from 'axios';
import { Job, CreateJobRequest, JobListResponse } from '../types/job';

const API_URL = import.meta.env.VITE_API_URL || '/api';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const jobsApi = {
  health: () => api.get('/health'),
  
  createJob: (data: CreateJobRequest) => 
    api.post<Job>('/jobs', data),
  
  getJob: (id: string) => 
    api.get<Job>(`/jobs/${id}`),
  
  listJobs: (status?: string) => 
    api.get<JobListResponse>('/jobs', { params: { status } }),
};