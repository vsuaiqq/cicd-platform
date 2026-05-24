import { useEffect } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../store'
import { validateSession } from '../store/authSlice'
import { AppShellSkeleton } from './ui'

interface ProtectedRouteProps {
  children: React.ReactNode
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const dispatch = useAppDispatch()
  const location = useLocation()
  const { isAuthenticated, loading, sessionValidated } = useAppSelector((s) => s.auth)




  useEffect(() => {
    if (!sessionValidated) {
      dispatch(validateSession())
    }
  }, [dispatch, sessionValidated])


  if (!sessionValidated && loading) {
    return <AppShellSkeleton />
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  return <>{children}</>
}
