export interface Job {
  id: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  type: string;
  data: string;
  result?: {
    processed_at: string;
    input_count: number;
    message: string;
  };
  error?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateJobRequest {
  type: string;
  data: string;
}

export interface JobListResponse {
  count: number;
  jobs: Job[];
}