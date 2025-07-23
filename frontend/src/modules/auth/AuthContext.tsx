// AuthContext.tsx
import { createContext, useState, useContext, useEffect } from 'react'
import type { ReactNode } from 'react'
import AuthService from '../../services/AuthService'
import type { LoginResult } from '../../services/AuthService'

const USER_STORAGE_KEY = 'user'

interface UserInfo {
    id: string
    name: string
    email: string
}

type AuthContextType = {
    user: UserInfo | null
    accessToken: string
    isAuthenticated: boolean
    login: (email: string, password: string) => Promise<LoginResult>
    logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | null>(null)

interface AuthProviderProps {
    children: ReactNode
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
    const [user, setUser] = useState<UserInfo | null>(null)
    const [accessToken, setAccessToken] = useState('')

    // Computed property for authentication status
    const isAuthenticated = !!(user && accessToken)

    useEffect(() => {
        const tryRefresh = async () => {
            try {
                const res = await AuthService.refresh()
                if (res.success) {
                    setAccessToken(res.data.accessToken)
                } else {
                    setAccessToken('')
                }
            } catch {
                setAccessToken('')
            }
        }
        tryRefresh()
        const savedUser = localStorage.getItem(USER_STORAGE_KEY)
        if (savedUser) {
            setUser(JSON.parse(savedUser))
        }
    }, [])

    const login = async (
        email: string,
        password: string
    ): Promise<LoginResult> => {
        const result = await AuthService.login({ email, password })

        console.log(`AuthContext result:`)
        console.log(result)
        if (result.success) {
            if (result.data.userInfo) {
                localStorage.setItem(
                    USER_STORAGE_KEY,
                    JSON.stringify(result.data.userInfo)
                )
                setUser(result.data.userInfo)
            }
            if (result.data.accessToken) {
                setAccessToken(result.data.accessToken)
            }
        }
        return result
    }

    const logout = async () => {
        try {
            const result = await AuthService.logout()
            if (result.success) {
                // Clear local storage
                localStorage.removeItem(USER_STORAGE_KEY)
                // Reset state
                setUser(null)
                setAccessToken('')
            } else {
                console.error('Logout failed')
                // Still clear local state even if API call fails
                localStorage.removeItem(USER_STORAGE_KEY)
                setUser(null)
                setAccessToken('')
            }
        } catch (error) {
            console.error('Logout error:', error)
            // Clear local state even if there's an error
            localStorage.removeItem(USER_STORAGE_KEY)
            setUser(null)
            setAccessToken('')
        }
    }

    return (
        <AuthContext.Provider
            value={{ user, accessToken, isAuthenticated, login, logout }}
        >
            {children}
        </AuthContext.Provider>
    )
}

export const useAuth = () => useContext(AuthContext)

// Custom hook for authentication status
export const useAuthStatus = () => {
    const auth = useAuth()
    if (!auth) {
        throw new Error('useAuthStatus must be used within an AuthProvider')
    }

    return {
        isAuthenticated: auth.isAuthenticated,
        user: auth.user,
        accessToken: auth.accessToken,
    }
}
