import { Outlet, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function Layout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  if (!user) {
    navigate('/login')
    return null
  }

  return (
    <div>
      <nav className="nav">
        <div className="nav-brand">TaskMaster</div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
          <span style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>
            {user.name} ({user.role})
          </span>
          <button className="btn" onClick={logout} style={{ background: 'var(--border)' }}>
            Logout
          </button>
        </div>
      </nav>
      <main className="container">
        <Outlet />
      </main>
    </div>
  )
}
