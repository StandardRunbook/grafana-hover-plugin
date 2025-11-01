-- Clear existing data for CPU metrics
DELETE FROM logs WHERE org = '1' AND dashboard = 'CPU Usage';
DELETE FROM template_examples WHERE org = '1' AND dashboard = 'CPU Usage';

-- Insert realistic CPU-related logs with various patterns
-- Template 1: High CPU usage warnings
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 5 MINUTE, 'cpu_high_001', 'WARNING: CPU usage at 92% on node worker-01'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 4 MINUTE, 'cpu_high_001', 'WARNING: CPU usage at 94% on node worker-02'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 3 MINUTE, 'cpu_high_001', 'WARNING: CPU usage at 89% on node worker-03'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 2 MINUTE, 'cpu_high_001', 'WARNING: CPU usage at 96% on node worker-01'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 1 MINUTE, 'cpu_high_001', 'WARNING: CPU usage at 91% on node worker-04');

-- Template 2: CPU throttling events
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 10 MINUTE, 'cpu_throttle_002', 'CPU throttling detected: process nginx (PID 1234) throttled for 250ms'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 8 MINUTE, 'cpu_throttle_002', 'CPU throttling detected: process postgres (PID 5678) throttled for 180ms'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 6 MINUTE, 'cpu_throttle_002', 'CPU throttling detected: process redis (PID 9012) throttled for 320ms');

-- Template 3: Context switch storms
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 15 MINUTE, 'cpu_context_003', 'High context switch rate detected: 45000 switches/sec on core 0'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 13 MINUTE, 'cpu_context_003', 'High context switch rate detected: 52000 switches/sec on core 2'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 11 MINUTE, 'cpu_context_003', 'High context switch rate detected: 48000 switches/sec on core 1'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 9 MINUTE, 'cpu_context_003', 'High context switch rate detected: 61000 switches/sec on core 3');

-- Template 4: Load average spikes
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 20 MINUTE, 'cpu_load_004', 'Load average spike: 1min=8.2, 5min=6.1, 15min=4.3 on host api-server-01'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 18 MINUTE, 'cpu_load_004', 'Load average spike: 1min=9.5, 5min=7.2, 15min=5.1 on host api-server-02'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 16 MINUTE, 'cpu_load_004', 'Load average spike: 1min=7.8, 5min=5.9, 15min=4.7 on host api-server-03');

-- Template 5: Process consuming high CPU
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 25 MINUTE, 'cpu_process_005', 'Process java (PID 2345) consuming 85% CPU for 30 seconds'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 23 MINUTE, 'cpu_process_005', 'Process python (PID 6789) consuming 92% CPU for 45 seconds'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 21 MINUTE, 'cpu_process_005', 'Process node (PID 3456) consuming 78% CPU for 25 seconds'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 19 MINUTE, 'cpu_process_005', 'Process go-app (PID 7890) consuming 88% CPU for 40 seconds');

-- Template 6: CPU steal time (for VMs)
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 30 MINUTE, 'cpu_steal_006', 'High CPU steal time detected: 15% on VM instance i-abc123'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 28 MINUTE, 'cpu_steal_006', 'High CPU steal time detected: 22% on VM instance i-def456'),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', now() - INTERVAL 26 MINUTE, 'cpu_steal_006', 'High CPU steal time detected: 18% on VM instance i-ghi789');

-- Insert template examples (representative logs for each template)
INSERT INTO template_examples (org, dashboard, panel_name, metric_name, template_id, message, timestamp) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_high_001', 'WARNING: CPU usage at 92% on node worker-01', now() - INTERVAL 5 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_high_001', 'WARNING: CPU usage at 94% on node worker-02', now() - INTERVAL 4 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_high_001', 'WARNING: CPU usage at 89% on node worker-03', now() - INTERVAL 3 MINUTE);

INSERT INTO template_examples (org, dashboard, panel_name, metric_name, template_id, message, timestamp) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_throttle_002', 'CPU throttling detected: process nginx (PID 1234) throttled for 250ms', now() - INTERVAL 10 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_throttle_002', 'CPU throttling detected: process postgres (PID 5678) throttled for 180ms', now() - INTERVAL 8 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_throttle_002', 'CPU throttling detected: process redis (PID 9012) throttled for 320ms', now() - INTERVAL 6 MINUTE);

INSERT INTO template_examples (org, dashboard, panel_name, metric_name, template_id, message, timestamp) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_context_003', 'High context switch rate detected: 45000 switches/sec on core 0', now() - INTERVAL 15 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_context_003', 'High context switch rate detected: 52000 switches/sec on core 2', now() - INTERVAL 13 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_context_003', 'High context switch rate detected: 48000 switches/sec on core 1', now() - INTERVAL 11 MINUTE);

INSERT INTO template_examples (org, dashboard, panel_name, metric_name, template_id, message, timestamp) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_load_004', 'Load average spike: 1min=8.2, 5min=6.1, 15min=4.3 on host api-server-01', now() - INTERVAL 20 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_load_004', 'Load average spike: 1min=9.5, 5min=7.2, 15min=5.1 on host api-server-02', now() - INTERVAL 18 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_load_004', 'Load average spike: 1min=7.8, 5min=5.9, 15min=4.7 on host api-server-03', now() - INTERVAL 16 MINUTE);

INSERT INTO template_examples (org, dashboard, panel_name, metric_name, template_id, message, timestamp) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_process_005', 'Process java (PID 2345) consuming 85% CPU for 30 seconds', now() - INTERVAL 25 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_process_005', 'Process python (PID 6789) consuming 92% CPU for 45 seconds', now() - INTERVAL 23 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_process_005', 'Process node (PID 3456) consuming 78% CPU for 25 seconds', now() - INTERVAL 21 MINUTE);

INSERT INTO template_examples (org, dashboard, panel_name, metric_name, template_id, message, timestamp) VALUES
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_steal_006', 'High CPU steal time detected: 15% on VM instance i-abc123', now() - INTERVAL 30 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_steal_006', 'High CPU steal time detected: 22% on VM instance i-def456', now() - INTERVAL 28 MINUTE),
('1', 'CPU Usage', 'CPU Usage', 'cpu_usage', 'cpu_steal_006', 'High CPU steal time detected: 18% on VM instance i-ghi789', now() - INTERVAL 26 MINUTE);
