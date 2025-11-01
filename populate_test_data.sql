-- Populate test data for log-stream-centric schema

-- 1. Create organization
INSERT INTO organizations (id, name) VALUES ('1', 'Acme Corp');

-- 2. Create metric
INSERT INTO metrics (id, org_id, dashboard_name, panel_title, metric_name)
VALUES ('cpu_metric_1', '1', 'CPU Usage', 'CPU Usage', 'cpu_usage');

-- 3. Create log streams
INSERT INTO log_streams (id, org_id, service, region, log_stream_name) VALUES
('stream_api_east', '1', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1'),
('stream_api_west', '1', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2'),
('stream_worker_east', '1', 'worker', 'us-east-1', '/aws/ecs/worker/us-east-1');

-- 4. Map log streams to the CPU metric
INSERT INTO metric_log_mappings (id, org_id, metric_id, log_stream_id, is_active) VALUES
('mapping_1', '1', 'cpu_metric_1', 'stream_api_east', 1),
('mapping_2', '1', 'cpu_metric_1', 'stream_api_west', 1),
('mapping_3', '1', 'cpu_metric_1', 'stream_worker_east', 1);

-- 5. Add baseline logs (1-2 hours ago) with normal CPU patterns
INSERT INTO logs (org_id, log_stream_id, service, region, log_stream_name, timestamp, template_id, message) VALUES
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 90 MINUTE, 'cpu_normal_001', 'INFO: CPU usage at 45% on api-server-01'),
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 85 MINUTE, 'cpu_normal_001', 'INFO: CPU usage at 42% on api-server-02'),
('1', 'stream_api_west', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2', now() - INTERVAL 80 MINUTE, 'cpu_normal_001', 'INFO: CPU usage at 48% on api-server-03'),
('1', 'stream_worker_east', 'worker', 'us-east-1', '/aws/ecs/worker/us-east-1', now() - INTERVAL 75 MINUTE, 'cpu_normal_001', 'INFO: CPU usage at 50% on worker-01'),
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 70 MINUTE, 'cpu_normal_001', 'INFO: CPU usage at 44% on api-server-01'),
('1', 'stream_api_west', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2', now() - INTERVAL 68 MINUTE, 'cpu_throttle_002', 'WARN: CPU throttling detected on api-server-03'),
('1', 'stream_worker_east', 'worker', 'us-east-1', '/aws/ecs/worker/us-east-1', now() - INTERVAL 65 MINUTE, 'cpu_context_003', 'DEBUG: Context switches: 1500/sec on worker-01'),
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 62 MINUTE, 'cpu_throttle_002', 'WARN: CPU throttling detected on api-server-01'),
('1', 'stream_api_west', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2', now() - INTERVAL 60 MINUTE, 'cpu_context_003', 'DEBUG: Context switches: 1450/sec on api-server-03');

-- 6. Add current logs (last hour) with elevated CPU issues
INSERT INTO logs (org_id, log_stream_id, service, region, log_stream_name, timestamp, template_id, message) VALUES
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 55 MINUTE, 'cpu_high_001', 'WARNING: CPU usage at 92% on api-server-01'),
('1', 'stream_api_west', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2', now() - INTERVAL 50 MINUTE, 'cpu_high_001', 'WARNING: CPU usage at 89% on api-server-03'),
('1', 'stream_worker_east', 'worker', 'us-east-1', '/aws/ecs/worker/us-east-1', now() - INTERVAL 48 MINUTE, 'cpu_high_001', 'WARNING: CPU usage at 95% on worker-01'),
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 45 MINUTE, 'cpu_high_001', 'WARNING: CPU usage at 91% on api-server-02'),
('1', 'stream_api_west', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2', now() - INTERVAL 42 MINUTE, 'cpu_high_001', 'WARNING: CPU usage at 93% on api-server-03'),
('1', 'stream_worker_east', 'worker', 'us-east-1', '/aws/ecs/worker/us-east-1', now() - INTERVAL 40 MINUTE, 'cpu_process_005', 'ERROR: Process consuming 80% CPU: java'),
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 38 MINUTE, 'cpu_process_005', 'ERROR: Process consuming 75% CPU: node'),
('1', 'stream_api_west', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2', now() - INTERVAL 35 MINUTE, 'cpu_process_005', 'ERROR: Process consuming 82% CPU: node'),
('1', 'stream_worker_east', 'worker', 'us-east-1', '/aws/ecs/worker/us-east-1', now() - INTERVAL 32 MINUTE, 'cpu_process_005', 'ERROR: Process consuming 78% CPU: python'),
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 30 MINUTE, 'cpu_context_003', 'DEBUG: Context switches: 3200/sec on api-server-01'),
('1', 'stream_api_west', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2', now() - INTERVAL 28 MINUTE, 'cpu_context_003', 'DEBUG: Context switches: 3500/sec on api-server-03'),
('1', 'stream_worker_east', 'worker', 'us-east-1', '/aws/ecs/worker/us-east-1', now() - INTERVAL 25 MINUTE, 'cpu_context_003', 'DEBUG: Context switches: 3100/sec on worker-01'),
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 22 MINUTE, 'cpu_context_003', 'DEBUG: Context switches: 3300/sec on api-server-02'),
('1', 'stream_api_west', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2', now() - INTERVAL 20 MINUTE, 'cpu_throttle_002', 'WARN: CPU throttling detected on api-server-03'),
('1', 'stream_worker_east', 'worker', 'us-east-1', '/aws/ecs/worker/us-east-1', now() - INTERVAL 18 MINUTE, 'cpu_throttle_002', 'WARN: CPU throttling detected on worker-01'),
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 15 MINUTE, 'cpu_throttle_002', 'WARN: CPU throttling detected on api-server-01'),
('1', 'stream_api_west', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2', now() - INTERVAL 12 MINUTE, 'cpu_steal_006', 'ALERT: CPU steal time at 15% on api-server-03'),
('1', 'stream_worker_east', 'worker', 'us-east-1', '/aws/ecs/worker/us-east-1', now() - INTERVAL 10 MINUTE, 'cpu_steal_006', 'ALERT: CPU steal time at 12% on worker-01'),
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 8 MINUTE, 'cpu_steal_006', 'ALERT: CPU steal time at 18% on api-server-01'),
('1', 'stream_api_west', 'api-server', 'us-west-2', '/aws/ecs/api-server/us-west-2', now() - INTERVAL 5 MINUTE, 'cpu_load_004', 'WARN: Load average: 8.5, 7.2, 6.1 on api-server-03'),
('1', 'stream_worker_east', 'worker', 'us-east-1', '/aws/ecs/worker/us-east-1', now() - INTERVAL 3 MINUTE, 'cpu_load_004', 'WARN: Load average: 9.1, 8.0, 7.5 on worker-01'),
('1', 'stream_api_east', 'api-server', 'us-east-1', '/aws/ecs/api-server/us-east-1', now() - INTERVAL 1 MINUTE, 'cpu_load_004', 'WARN: Load average: 8.8, 7.8, 6.9 on api-server-02');

-- 7. Add template examples
INSERT INTO template_examples (org_id, log_stream_id, service, region, template_id, message, timestamp) VALUES
('1', 'stream_api_east', 'api-server', 'us-east-1', 'cpu_normal_001', 'INFO: CPU usage at 45% on api-server-01', now() - INTERVAL 90 MINUTE),
('1', 'stream_api_east', 'api-server', 'us-east-1', 'cpu_high_001', 'WARNING: CPU usage at 92% on api-server-01', now() - INTERVAL 55 MINUTE),
('1', 'stream_worker_east', 'worker', 'us-east-1', 'cpu_process_005', 'ERROR: Process consuming 80% CPU: java', now() - INTERVAL 40 MINUTE),
('1', 'stream_api_east', 'api-server', 'us-east-1', 'cpu_context_003', 'DEBUG: Context switches: 3200/sec on api-server-01', now() - INTERVAL 30 MINUTE),
('1', 'stream_api_west', 'api-server', 'us-west-2', 'cpu_throttle_002', 'WARN: CPU throttling detected on api-server-03', now() - INTERVAL 20 MINUTE),
('1', 'stream_api_east', 'api-server', 'us-east-1', 'cpu_steal_006', 'ALERT: CPU steal time at 18% on api-server-01', now() - INTERVAL 8 MINUTE),
('1', 'stream_worker_east', 'worker', 'us-east-1', 'cpu_load_004', 'WARN: Load average: 9.1, 8.0, 7.5 on worker-01', now() - INTERVAL 3 MINUTE);
