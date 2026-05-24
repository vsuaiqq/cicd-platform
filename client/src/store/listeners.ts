import { createListenerMiddleware, isAnyOf } from '@reduxjs/toolkit'
import { apiSlice } from './api/apiSlice'
import { logout } from './authSlice'

export const listenerMiddleware = createListenerMiddleware()

listenerMiddleware.startListening({
  matcher: isAnyOf(logout.fulfilled),
  effect: async (_action, api) => {
    api.dispatch(apiSlice.util.resetApiState())
  },
})
