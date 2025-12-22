import { QueryClient, QueryClientProvider, useQuery, useMutation } from '@tanstack/react-query';
import { useState } from 'react';
import { jobsApi } from './services/api';
import { CreateJobRequest } from './types/job';
import './index.css';
import './App.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchInterval: 2000,
    },
  },
});

function JobDashboard() {
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [jobType, setJobType] = useState('data-processing');
  const [jobData, setJobData] = useState('');

  const { data: health } = useQuery({
    queryKey: ['health'],
    queryFn: async () => {
      const response = await jobsApi.health();
      return response.data;
    },
  });

  const { data: jobs, isLoading } = useQuery({
    queryKey: ['jobs', statusFilter],
    queryFn: async () => {
      const response = await jobsApi.listJobs(statusFilter || undefined);
      return response.data;
    },
  });

  const createJobMutation = useMutation({
    mutationFn: (data: CreateJobRequest) => jobsApi.createJob(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['jobs'] });
      setJobData('');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createJobMutation.mutate({
      type: jobType,
      data: jobData,
    });
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return '#10b981';
      case 'processing': return '#3b82f6';
      case 'failed': return '#ef4444';
      default: return '#6b7280';
    }
  };

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#f3f4f6' }}>
      <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
        {/* Header */}
        <div style={{ marginBottom: '2rem' }}>
          <h1 className="gradient-text" style={{ fontSize: '3rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
            Job Management Dashboard
          </h1>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <div style={{
              width: '8px',
              height: '8px',
              borderRadius: '50%',
              backgroundColor: health?.status === 'healthy' ? '#10b981' : '#ef4444',
            }}></div>
            <span style={{ fontSize: '0.875rem', color: '#6b7280' }}>
              API Status: <span style={{ fontWeight: '500', color: health?.status === 'healthy' ? '#10b981' : '#ef4444' }}>
                {health?.status || 'checking...'}
              </span>
            </span>
          </div>
        </div>

        {/* Create Job Form */}
        <div className="card-shadow" style={{
          backgroundColor: 'white',
          borderRadius: '12px',
          padding: '2rem',
          marginBottom: '2rem',
        }}>
          <h2 style={{ fontSize: '1.5rem', fontWeight: '600', marginBottom: '1.5rem', color: '#1f2937' }}>
            Create New Job
          </h2>
          <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            <div>
              <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' }}>
                Job Type
              </label>
              <select
                value={jobType}
                onChange={(e) => setJobType(e.target.value)}
                style={{
                  width: '100%',
                  padding: '0.75rem',
                  border: '2px solid #e5e7eb',
                  borderRadius: '8px',
                  fontSize: '1rem',
                  outline: 'none',
                }}
              >
                <option value="data-processing">Data Processing</option>
                <option value="health-report">Health Report</option>
                <option value="cleanup">Cleanup</option>
                <option value="analytics-aggregation">Analytics Aggregation</option>
              </select>
            </div>
            <div>
              <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' }}>
                Job Data
              </label>
              <textarea
                value={jobData}
                onChange={(e) => setJobData(e.target.value)}
                required
                style={{
                  width: '100%',
                  padding: '0.75rem',
                  border: '2px solid #e5e7eb',
                  borderRadius: '8px',
                  fontSize: '1rem',
                  outline: 'none',
                  resize: 'none',
                }}
                rows={3}
                placeholder="Enter job data..."
              />
            </div>
            <button
              type="submit"
              disabled={createJobMutation.isPending}
              className="btn-gradient"
              style={{
                alignSelf: 'flex-start',
                padding: '0.75rem 2rem',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                fontSize: '1rem',
                fontWeight: '500',
                cursor: createJobMutation.isPending ? 'not-allowed' : 'pointer',
                opacity: createJobMutation.isPending ? 0.7 : 1,
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem',
              }}
            >
              {createJobMutation.isPending && <div className="spinner"></div>}
              {createJobMutation.isPending ? 'Creating...' : 'Create Job'}
            </button>
          </form>
        </div>

        {/* Jobs List */}
        <div className="card-shadow" style={{
          backgroundColor: 'white',
          borderRadius: '12px',
          overflow: 'hidden',
        }}>
          <div style={{ padding: '2rem', borderBottom: '1px solid #e5e7eb' }}>
            <h2 style={{ fontSize: '1.5rem', fontWeight: '600', marginBottom: '1rem', color: '#1f2937' }}>
              Jobs
            </h2>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <label style={{ fontSize: '0.875rem', fontWeight: '500', color: '#374151' }}>
                Filter by status:
              </label>
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                style={{
                  padding: '0.5rem 1rem',
                  border: '2px solid #e5e7eb',
                  borderRadius: '6px',
                  fontSize: '0.875rem',
                  outline: 'none',
                }}
              >
                <option value="">All</option>
                <option value="pending">Pending</option>
                <option value="processing">Processing</option>
                <option value="completed">Completed</option>
                <option value="failed">Failed</option>
              </select>
            </div>
          </div>
          
          {isLoading ? (
            <div style={{ padding: '2rem' }}>
              {[1, 2, 3].map((i) => (
                <div key={i} style={{
                  height: '80px',
                  backgroundColor: '#f3f4f6',
                  borderRadius: '8px',
                  marginBottom: '1rem',
                  animation: 'pulse 1.5s ease-in-out infinite',
                }}></div>
              ))}
            </div>
          ) : (
            <div>
              {jobs?.jobs.length === 0 ? (
                <div style={{ padding: '4rem', textAlign: 'center' }}>
                  <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>📋</div>
                  <p style={{ color: '#6b7280', fontWeight: '500' }}>No jobs found</p>
                  <p style={{ color: '#9ca3af', fontSize: '0.875rem', marginTop: '0.5rem' }}>
                    Create a new job to get started
                  </p>
                </div>
              ) : (
                jobs?.jobs.map((job) => (
                  <div key={job.id} className="job-card" style={{
                    padding: '1.5rem 2rem',
                    borderBottom: '1px solid #e5e7eb',
                  }}>
                    <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between' }}>
                      <div style={{ flex: 1 }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '0.5rem' }}>
                          <span style={{ fontSize: '1.125rem', fontWeight: '600', color: '#1f2937' }}>
                            {job.type}
                          </span>
                          <span className="status-badge" style={{
                            padding: '0.25rem 0.75rem',
                            fontSize: '0.75rem',
                            fontWeight: '600',
                            color: 'white',
                            backgroundColor: getStatusColor(job.status),
                            borderRadius: '9999px',
                            textTransform: 'uppercase',
                          }}>
                            {job.status}
                          </span>
                        </div>
                        <p style={{ color: '#4b5563', marginBottom: '0.75rem' }}>{job.data}</p>
                        {job.result && (
                          <div style={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: '0.5rem',
                            padding: '0.5rem 0.75rem',
                            backgroundColor: '#d1fae5',
                            color: '#065f46',
                            borderRadius: '6px',
                            fontSize: '0.875rem',
                            marginBottom: '0.5rem',
                          }}>
                            <span>✓</span>
                            <span>{job.result.message} ({job.result.input_count} items)</span>
                          </div>
                        )}
                        {job.error && (
                          <div style={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: '0.5rem',
                            padding: '0.5rem 0.75rem',
                            backgroundColor: '#fee2e2',
                            color: '#991b1b',
                            borderRadius: '6px',
                            fontSize: '0.875rem',
                            marginBottom: '0.5rem',
                          }}>
                            <span>✗</span>
                            <span>{job.error}</span>
                          </div>
                        )}
                        <p style={{ fontSize: '0.75rem', color: '#9ca3af', display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                          <span>🕒</span>
                          {new Date(job.created_at).toLocaleString()}
                        </p>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <JobDashboard />
    </QueryClientProvider>
  );
}

export default App;