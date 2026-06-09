@echo off
setlocal

echo Seeding database...

:: Create temporary SQL file
echo -- Seed data for local development and testing > temp_seed.sql
echo. >> temp_seed.sql
echo -- Insert demo teams >> temp_seed.sql
echo INSERT INTO teams (id, name, sla_days, task_limit) VALUES >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440000', 'Engineering', 14, 100), >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440001', 'Design', 7, 50), >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440002', 'Marketing', 21, 30) >> temp_seed.sql
echo ON CONFLICT (id) DO NOTHING; >> temp_seed.sql
echo. >> temp_seed.sql
echo -- Insert demo users >> temp_seed.sql
echo INSERT INTO users (id, email, password, name, role, team_id, created_at) VALUES >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440010', 'admin@taskmaster.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqBmW1PqQh7/1tL3z7xL3z7xL3z7x', 'Admin User', 'admin', '550e8400-e29b-41d4-a716-446655440000', NOW()), >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440011', 'manager@taskmaster.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqBmW1PqQh7/1tL3z7xL3z7xL3z7x', 'Manager User', 'manager', '550e8400-e29b-41d4-a716-446655440000', NOW()), >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440012', 'member@taskmaster.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqBmW1PqQh7/1tL3z7xL3z7xL3z7x', 'Member User', 'member', '550e8400-e29b-41d4-a716-446655440000', NOW()) >> temp_seed.sql
echo ON CONFLICT (id) DO NOTHING; >> temp_seed.sql
echo. >> temp_seed.sql
echo -- Insert team memberships >> temp_seed.sql
echo INSERT INTO team_members (team_id, user_id, joined_at) VALUES >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440010', NOW()), >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440011', NOW()), >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440012', NOW()) >> temp_seed.sql
echo ON CONFLICT (team_id, user_id) DO NOTHING; >> temp_seed.sql
echo. >> temp_seed.sql
echo -- Insert demo tasks >> temp_seed.sql
echo INSERT INTO tasks (id, title, description, status, priority, assignee_id, creator_id, team_id, due_date, created_at, updated_at) VALUES >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440020', 'Build authentication system', 'Implement JWT-based auth with bcrypt passwords', 'in_progress', 'high', '550e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440000', NOW() + INTERVAL '7 days', NOW(), NOW()), >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440021', 'Design task dashboard UI', 'Create React components for task listing and creation', 'todo', 'medium', '550e8400-e29b-41d4-a716-446655440012', '550e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440000', NOW() + INTERVAL '14 days', NOW(), NOW()), >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440022', 'Setup AWS infrastructure', 'Terraform modules for VPC, ECS, RDS, ElastiCache', 'todo', 'urgent', '550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440000', NOW() + INTERVAL '3 days', NOW(), NOW()), >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440023', 'Write integration tests', 'Testcontainers for PostgreSQL and Redis integration', 'todo', 'high', '550e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440000', NOW() + INTERVAL '10 days', NOW(), NOW()), >> temp_seed.sql
echo     ('550e8400-e29b-41d4-a716-446655440024', 'Setup CI/CD pipeline', 'GitHub Actions for lint, test, build, push, deploy', 'done', 'medium', '550e8400-e29b-41d4-a716-446655440012', '550e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440000', NOW() - INTERVAL '2 days', NOW() - INTERVAL '5 days', NOW()) >> temp_seed.sql
echo ON CONFLICT (id) DO NOTHING; >> temp_seed.sql

:: Execute seed SQL via docker exec
docker exec -i taskmaster-postgres psql -U taskmaster -d taskmaster < temp_seed.sql

:: Clean up
del temp_seed.sql

echo Seed complete!
endlocal
