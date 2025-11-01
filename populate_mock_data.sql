-- Insert mock log data for testing
-- Simulating logs from two different panels with various error patterns

-- Insert logs for the first panel (API errors)
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:00:00', 'tmpl_001', 'Connection timeout to database server 10.0.1.5'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:01:00', 'tmpl_001', 'Connection timeout to database server 10.0.1.8'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:02:00', 'tmpl_001', 'Connection timeout to database server 10.0.1.12'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:03:00', 'tmpl_002', 'Failed to authenticate user john.doe@example.com'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:04:00', 'tmpl_002', 'Failed to authenticate user jane.smith@example.com'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:05:00', 'tmpl_003', 'Rate limit exceeded for API key abc123def456'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:06:00', 'tmpl_003', 'Rate limit exceeded for API key xyz789ghi012'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:07:00', 'tmpl_003', 'Rate limit exceeded for API key mno345pqr678'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:08:00', 'tmpl_004', 'NULL pointer exception in handler /api/v1/users'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:09:00', 'tmpl_004', 'NULL pointer exception in handler /api/v1/orders'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:10:00', 'tmpl_005', 'Disk space critical on node worker-03: 95% used'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', '2025-11-01 00:11:00', 'tmpl_005', 'Disk space critical on node worker-07: 97% used');

-- Insert logs for the second panel (Database errors)
INSERT INTO logs (org, dashboard, panel_name, metric_name, timestamp, template_id, message) VALUES
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', '2025-11-01 00:00:00', 'tmpl_101', 'Deadlock detected in transaction 12345'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', '2025-11-01 00:02:00', 'tmpl_101', 'Deadlock detected in transaction 12389'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', '2025-11-01 00:04:00', 'tmpl_102', 'Slow query warning: SELECT took 5.2s on table users'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', '2025-11-01 00:05:00', 'tmpl_102', 'Slow query warning: SELECT took 8.7s on table orders'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', '2025-11-01 00:06:00', 'tmpl_102', 'Slow query warning: SELECT took 12.1s on table analytics'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', '2025-11-01 00:08:00', 'tmpl_103', 'Connection pool exhausted: max 100 connections reached'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', '2025-11-01 00:09:00', 'tmpl_103', 'Connection pool exhausted: max 100 connections reached'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', '2025-11-01 00:10:00', 'tmpl_104', 'Replication lag detected: replica-02 is 45s behind master');

-- Insert template examples (representative logs for each template)
INSERT INTO template_examples (org, dashboard, panel_name, metric_name, template_id, message, timestamp) VALUES
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_001', 'Connection timeout to database server 10.0.1.5', '2025-11-01 00:00:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_001', 'Connection timeout to database server 10.0.1.8', '2025-11-01 00:01:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_001', 'Connection timeout to database server 10.0.1.12', '2025-11-01 00:02:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_002', 'Failed to authenticate user john.doe@example.com', '2025-11-01 00:03:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_002', 'Failed to authenticate user jane.smith@example.com', '2025-11-01 00:04:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_003', 'Rate limit exceeded for API key abc123def456', '2025-11-01 00:05:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_003', 'Rate limit exceeded for API key xyz789ghi012', '2025-11-01 00:06:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_003', 'Rate limit exceeded for API key mno345pqr678', '2025-11-01 00:07:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_004', 'NULL pointer exception in handler /api/v1/users', '2025-11-01 00:08:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_004', 'NULL pointer exception in handler /api/v1/orders', '2025-11-01 00:09:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_005', 'Disk space critical on node worker-03: 95% used', '2025-11-01 00:10:00'),
('acme-corp', 'production-dashboard', 'API Logs', 'error_rate', 'tmpl_005', 'Disk space critical on node worker-07: 97% used', '2025-11-01 00:11:00'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', 'tmpl_101', 'Deadlock detected in transaction 12345', '2025-11-01 00:00:00'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', 'tmpl_101', 'Deadlock detected in transaction 12389', '2025-11-01 00:02:00'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', 'tmpl_102', 'Slow query warning: SELECT took 5.2s on table users', '2025-11-01 00:04:00'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', 'tmpl_102', 'Slow query warning: SELECT took 8.7s on table orders', '2025-11-01 00:05:00'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', 'tmpl_102', 'Slow query warning: SELECT took 12.1s on table analytics', '2025-11-01 00:06:00'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', 'tmpl_103', 'Connection pool exhausted: max 100 connections reached', '2025-11-01 00:08:00'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', 'tmpl_103', 'Connection pool exhausted: max 100 connections reached', '2025-11-01 00:09:00'),
('acme-corp', 'production-dashboard', 'Database Logs', 'query_errors', 'tmpl_104', 'Replication lag detected: replica-02 is 45s behind master', '2025-11-01 00:10:00');
