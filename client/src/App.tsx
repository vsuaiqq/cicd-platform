import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Provider } from 'react-redux'
import { store } from './store'
import { ProtectedRoute } from './components/ProtectedRoute'
import AppLayout from './components/AppLayout'
import AuthLayout from './components/AuthLayout'
import { ToastProvider } from './components/ui'
import { ThemeProvider } from './lib/themeContext'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import Profile from './pages/Profile'
import CreateProject from './pages/CreateProject'
import ProjectSetup from './pages/ProjectSetup'
import ProjectDetail from './pages/ProjectDetail'
import PipelineRunDetail from './pages/PipelineRunDetail'

function AppRoutes() {
  return (
    <Routes>
      <Route element={<AuthLayout />}>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
      </Route>
      <Route
        element={
          <ProtectedRoute>
            <AppLayout />
          </ProtectedRoute>
        }
      >
        <Route path="/" element={<Dashboard />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/projects/new" element={<CreateProject />} />
        <Route path="/projects/:id/setup" element={<ProjectSetup />} />
        <Route path="/projects/:id/runs/:runId" element={<PipelineRunDetail />} />
        <Route path="/projects/:id" element={<ProjectDetail />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default function App() {
  return (
    <Provider store={store}>
      <BrowserRouter>
        <ThemeProvider>
          <ToastProvider>
            <AppRoutes />
          </ToastProvider>
        </ThemeProvider>
      </BrowserRouter>
    </Provider>
  )
}
