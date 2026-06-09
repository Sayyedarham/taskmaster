export interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'manager' | 'member';
}

export interface Task {
  id: string;
  title: string;
  description?: string;
  status: 'todo' | 'in_progress' | 'done' | 'blocked';
  priority: 'low' | 'medium' | 'high' | 'urgent';
  assignee_id?: string;
  creator_id: string;
  team_id: string;
  due_date?: string;
  created_at: string;
  updated_at: string;
}

export interface Team {
  id: string;
  name: string;
  sla_days: number;
  task_limit: number;
  created_at: string;
}

export interface WSMessage {
  type: 'task:created' | 'task:updated' | 'task:deleted';
  payload: Task | { id: string };
}
