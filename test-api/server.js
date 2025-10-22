const express = require('express');
const cors = require('cors');

const app = express();
const PORT = 3001;

// Middleware
app.use(cors());
app.use(express.json());

// Mock log data
const mockLogs = {
  'cpu_usage_percent': [
    {
      timestamp: '2025-01-13T10:30:00.000Z',
      level: 'ERROR',
      message: 'High CPU usage detected: 95%',
      source: 'system-monitor'
    },
    {
      timestamp: '2025-01-13T10:30:15.000Z',
      level: 'WARN',
      message: 'Thread pool exhausted, creating new threads',
      source: 'thread-manager'
    },
    {
      timestamp: '2025-01-13T10:30:30.000Z',
      level: 'INFO',
      message: 'CPU usage normalized to 45%',
      source: 'system-monitor'
    }
  ],
  'memory_usage_bytes': [
    {
      timestamp: '2025-01-13T10:30:00.000Z',
      level: 'ERROR',
      message: 'Memory allocation failed: Out of memory',
      source: 'memory-manager'
    },
    {
      timestamp: '2025-01-13T10:30:05.000Z',
      level: 'WARN',
      message: 'Falling back to alternative memory pool',
      source: 'memory-manager'
    }
  ],
  'response_time_ms': [
    {
      timestamp: '2025-01-13T10:30:00.000Z',
      level: 'WARN',
      message: 'Response time exceeded threshold: 5000ms',
      source: 'api-gateway'
    },
    {
      timestamp: '2025-01-13T10:30:10.000Z',
      level: 'INFO',
      message: 'Response time improved to 200ms',
      source: 'api-gateway'
    }
  ]
};

// API endpoint
app.post('/query_logs', (req, res) => {
  console.log('Received query:', req.body);
  
  const { metric_name, start_time, end_time } = req.body;
  
  // Get logs for the metric
  const logs = mockLogs[metric_name] || [];
  
  // Filter logs by time range
  const startTime = new Date(start_time);
  const endTime = new Date(end_time);
  
  const filteredLogs = logs.filter(log => {
    const logTime = new Date(log.timestamp);
    return logTime >= startTime && logTime <= endTime;
  });
  
  // Group logs and calculate relative change
  const logGroups = [
    {
      representative_logs: filteredLogs.map(log => 
        `${log.timestamp} [${log.level}] ${log.message}`
      ),
      relative_change: Math.random() * 100 - 50 // Random change between -50% and +50%
    }
  ];
  
  // If no logs found, return a sample
  if (logGroups[0].representative_logs.length === 0) {
    logGroups[0] = {
      representative_logs: [
        `${start_time} [INFO] No logs found for metric: ${metric_name}`,
        `${end_time} [INFO] Time range: ${start_time} to ${end_time}`
      ],
      relative_change: 0
    };
  }
  
  res.json({
    log_groups: logGroups
  });
});

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'healthy', timestamp: new Date().toISOString() });
});

// Start server
app.listen(PORT, '0.0.0.0', () => {
  console.log(`Mock API server running on port ${PORT}`);
  console.log(`Health check: http://localhost:${PORT}/health`);
  console.log(`Query logs: POST http://localhost:${PORT}/query_logs`);
});
