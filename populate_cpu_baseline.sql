-- Add baseline data (1-2 hours ago) for KL divergence comparison
-- This provides historical data to compare against

-- Baseline: Normal CPU usage patterns (1-2 hours ago)
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 90 MINUTE, 'cpu_normal_001', 'INFO: CPU usage at 45% on node worker-01'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 88 MINUTE, 'cpu_normal_001', 'INFO: CPU usage at 42% on node worker-02'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 86 MINUTE, 'cpu_normal_001', 'INFO: CPU usage at 48% on node worker-03'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 84 MINUTE, 'cpu_normal_001', 'INFO: CPU usage at 43% on node worker-04'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 82 MINUTE, 'cpu_normal_001', 'INFO: CPU usage at 46% on node worker-01');

-- Some throttling (less frequent than current)
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 95 MINUTE, 'cpu_throttle_002', 'CPU throttling detected: process apache (PID 4321) throttled for 80ms'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 85 MINUTE, 'cpu_throttle_002', 'CPU throttling detected: process mysql (PID 8765) throttled for 60ms');

-- Normal context switches (lower rate)
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 92 MINUTE, 'cpu_context_003', 'High context switch rate detected: 25000 switches/sec on core 0'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 87 MINUTE, 'cpu_context_003', 'High context switch rate detected: 28000 switches/sec on core 1');

-- Add template examples for baseline data
INSERT INTO template_examples (org, dashboard, panel_name, metric_name, template_id, message, timestamp) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_normal_001', 'INFO: CPU usage at 45% on node worker-01', now() - INTERVAL 90 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_normal_001', 'INFO: CPU usage at 42% on node worker-02', now() - INTERVAL 88 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_normal_001', 'INFO: CPU usage at 48% on node worker-03', now() - INTERVAL 86 MINUTE);
