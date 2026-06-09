import { useEffect, useState, useCallback } from 'react'
import { useAuth } from '../context/AuthContext'
import { useWebSocket } from '../hooks/useWebSocket'

interface Task {
  id: string
  title: string
  description?: string
  status: 'todo' | 'in_progress' | 'done' | 'blocked'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  assignee_id?: string
  creator_id: string
  team_id: string
  due_date?: string
  created_at: string
}

export default function Dashboard() {
  const { token, user } = useAuth()
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [newTask, setNewTask] = useState({ title: '', description: '', priority: 'medium' as const })
  const { lastMessage, connected } = useWebSocket('/ws')

  const fetchTasks = useCallback(async () => {
    const res = await fetch('/api/v1/tasks', {
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.ok) {
      const data = await res.json()
      setTasks(data.tasks || [])
    }
    setLoading(false)
  }, [token])

  useEffect(() => {
    fetchTasks()
  }, [fetchTasks])

  useEffect(() => {
    if (lastMessage) {
      try {
        const msg = JSON.parse(lastMessage)
        if (msg.type === 'task:created') {
          setTasks((prev) => [msg.payload, ...prev])
        } else if (msg.type === 'task:deleted') {
          setTasks((prev) => prev.filter((t) => t.id !== msg.payload.id))
        }
      } catch {
        // ignore invalid messages
      }
    }
  }, [lastMessage])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    const res = await fetch('/api/v1/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(newTask),
    })
    if (res.ok) {
      setNewTask({ title: '', description: '', priority: 'medium' })
      setShowForm(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this task?')) return
    await fetch(`/api/v1/tasks/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    })
  }

  const getStatusBadge = (status: string) => {
    const map: Record<string, string> = {
      todo: 'badge-todo',
      in_progress: 'badge-progress',
      done: 'badge-done',
      blocked: 'badge-blocked',
    }
    return map[status] || 'badge-todo'
  }

  const canCreate = user?.role === 'admin' || user?.role === 'manager'

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
        <div>
          <h1 style={{ fontSize: '1.5rem', fontWeight: 700 }}>Tasks</h1>
          <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>
            {connected ? '🟢 Real-time' : '🔴 Offline'} · {tasks.length} tasks
          </p>
        </div>
        {canCreate && (
          <button className="btn btn-primary" onClick={() => setShowForm(!showForm)}>
            {showForm ? 'Cancel' : '+ New Task'}
          </button>
        )}
      </div>

      {showForm && (
        <div className="card" style={{ marginBottom: '1.5rem' }}>
          <h3 style={{ marginBottom: '1rem' }}>Create Task</h3>
          <form onSubmit={handleCreate} style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
            <input
              className="input"
              placeholder="Task title"
              value={newTask.title}
              onChange={(e) => setNewTask({ ...newTask, title: e.target.value })}
              required
            />
            <input
              className="input"
              placeholder="Description (optional)"
              value={newTask.description}
              onChange={(e) => setNewTask({ ...newTask, description: e.target.value })}
            />
            <select
              className="input"
              value={newTask.priority}
              onChange={(e) => setNewTask({ ...newTask, priority: e.target.value as Task['priority'] })}
            >
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
              <option value="urgent">Urgent</option>
            </select>
            <button type="submit" className="btn btn-primary" style={{ alignSelf: 'flex-start' }}>
              Create Task
            </button>
          </form>
        </div>
      )}

      {loading ? (
        <p style={{ color: 'var(--text-muted)' }}>Loading tasks...</p>
      ) : tasks.length === 0 ? (
        <div className="card" style={{ textAlign: 'center', padding: '3rem' }}>
          <p style={{ color: 'var(--text-muted)' }}>No tasks yet. Create your first one!</p>
        </div>
      ) : (
        <div className="grid grid-3">
          {tasks.map((task) => (
            <div key={task.id} className="card" style={{ position: 'relative' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '0.5rem' }}>
                <span className={`badge ${getStatusBadge(task.status)}`}>
                  {task.status.replace('_', ' ')}
                </span>
                <span className={`priority-${task.priority}`} style={{ fontSize: '0.75rem', fontWeight: 600, textTransform: 'uppercase' }}>
                  {task.priority}
                </span>
              </div>
              <h3 style={{ fontSize: '1rem', fontWeight: 600, marginBottom: '0.5rem' }}>{task.title}</h3>
              {task.description && (
                <p style={{ fontSize: '0.875rem', color: 'var(--text-muted)', marginBottom: '0.75rem' }}>
                  {task.description}
                </p>
              )}
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', fontSize: '0.75rem', color: 'var(--text-muted)' }}>
                <span>{new Date(task.created_at).toLocaleDateString()}</span>
                <button
                  onClick={() => handleDelete(task.id)}
                  style={{
                    background: 'none',
                    border: 'none',
                    color: '#dc2626',
                    cursor: 'pointer',
                    fontSize: '0.75rem'
                  }}
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
