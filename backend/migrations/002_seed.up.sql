-- Seed data for local development and testing

-- Insert demo teams
INSERT INTO teams (id, name, sla_days, task_limit) VALUES 
    ('550e8400-e29b-41d4-a716-446655440000', 'Engineering', 14, 100),
    ('550e8400-e29b-41d4-a716-446655440001', 'Design', 7, 50),
    ('550e8400-e29b-41d4-a716-446655440002', 'Marketing', 21, 30)
ON CONFLICT (id) DO NOTHING;

-- Insert demo users with REAL bcrypt hash for "password123"
-- Hash: $2b$10$Lz1hTFm06SImd1vF4Zztx.C0vVmv.elGfgOI6rg4e1TjcZif4qvbu
INSERT INTO users (id, email, password, name, role, team_id, created_at) VALUES 
    ('550e8400-e29b-41d4-a716-446655440010', 'admin@taskmaster.com', '$2b$10$Lz1hTFm06SImd1vF4Zztx.C0vVmv.elGfgOI6rg4e1TjcZif4qvbu', 'Admin User', 'admin', '550e8400-e29b-41d4-a716-446655440000', NOW()),
    ('550e8400-e29b-41d4-a716-446655440011', 'manager@taskmaster.com', '$2b$10$Lz1hTFm06SImd1vF4Zztx.C0vVmv.elGfgOI6rg4e1TjcZif4qvbu', 'Manager User', 'manager', '550e8400-e29b-41d4-a716-446655440000', NOW()),
    ('550e8400-e29b-41d4-a716-446655440012', 'member@taskmaster.com', '$2b$10$Lz1hTFm06SImd1vF4Zztx.C0vVmv.elGfgOI6rg4e1TjcZif4qvbu', 'Member User', 'member', '550e8400-e29b-41d4-a716-446655440000', NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert team memberships
INSERT INTO team_members (team_id, user_id, joined_at) VALUES 
    ('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440010', NOW()),
    ('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440011', NOW()),
    ('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440012', NOW())
ON CONFLICT (team_id, user_id) DO NOTHING;

-- Insert demo tasks
INSERT INTO tasks (id, title, description, status, priority, assignee_id, creator_id, team_id, due_date, created_at, updated_at) VALUES 
    ('550e8400-e29b-41d4-a716-446655440020', 'Build authentication system', 'Implement JWT-based auth with bcrypt passwords', 'in_progress', 'high', '550e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440000', NOW() + INTERVAL '7 days', NOW(), NOW()),
    ('550e8400-e29b-41d4-a716-446655440021', 'Design task dashboard UI', 'Create React components for task listing and creation', 'todo', 'medium', '550e8400-e29b-41d4-a716-446655440012', '550e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440000', NOW() + INTERVAL '14 days', NOW(), NOW()),
    ('550e8400-e29b-41d4-a716-446655440022', 'Setup AWS infrastructure', 'Terraform modules for VPC, ECS, RDS, ElastiCache', 'todo', 'urgent', '550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440000', NOW() + INTERVAL '3 days', NOW(), NOW()),
    ('550e8400-e29b-41d4-a716-446655440023', 'Write integration tests', 'Testcontainers for PostgreSQL and Redis integration', 'todo', 'high', '550e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440000', NOW() + INTERVAL '10 days', NOW(), NOW()),
    ('550e8400-e29b-41d4-a716-446655440024', 'Setup CI/CD pipeline', 'GitHub Actions for lint, test, build, push, deploy', 'done', 'medium', '550e8400-e29b-41d4-a716-446655440012', '550e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440000', NOW() - INTERVAL '2 days', NOW() - INTERVAL '5 days', NOW())
ON CONFLICT (id) DO NOTHING;