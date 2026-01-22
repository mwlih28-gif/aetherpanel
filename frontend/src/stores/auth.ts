import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface User {
  id: string
  email: string
  username: string
  firstName: string
  lastName: string
  role: string
  avatar?: string
  twoFactorEnabled: boolean
  credits: number
}

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isLoading: boolean
  setUser: (user: User | null) => void
  setTokens: (accessToken: string, refreshToken: string) => void
  login: (user: User, accessToken: string, refreshToken: string) => void
  logout: () => void
  updateUser: (updates: Partial<User>) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,

      setUser: (user: User | null) => set({ user, isAuthenticated: !!user }),
      
      setTokens: (accessToken: string, refreshToken: string) => 
        set({ accessToken, refreshToken }),
      
      login: (user: User, accessToken: string, refreshToken: string) =>
        set({
          user,
          accessToken,
          refreshToken,
          isAuthenticated: true,
        }),
      
      logout: () =>
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
        }),
      
      updateUser: (updates: Partial<User>) =>
        set((state: AuthState) => ({
          user: state.user ? { ...state.user, ...updates } : null,
        })),
    }),
    {
      name: 'aether-auth',
      partialize: (state: AuthState) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
