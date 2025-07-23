// AuthContext.tsx
import { createContext, useState, useContext } from 'react'
import type { ReactNode } from 'react'
import AuthService from '../../services/AuthService'
import type { LoginResult } from '../../services/AuthService'

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

    // useEffect(() => {
    //     const tryRefresh = async () => {
    //         try {
    //             const res = await api.post("/refresh");
    //             setUser(res.data.user);
    //             setAccessToken(res.data.accessToken);
    //         } catch {
    //             setUser(null);
    //         }
    //     };
    //     tryRefresh();
    // }, []);

    const login = async (
        email: string,
        password: string
    ): Promise<LoginResult> => {
        const result = await AuthService.loginPost({ email, password })

        console.log(`AuthContext result:`)
        console.log(result)
        if (result.success) {
            if (result.data.userInfo) {
                setUser(result.data.userInfo)
            }
            if (result.data.accessToken) {
                setAccessToken(result.data.accessToken)
            }
        }
        return result
    }

    const logout = async () => {
        // TODO: Implement logout API call
        // await api.post("/logout");
        setUser(null)
        setAccessToken('')
    }

    return (
        <AuthContext.Provider value={{ user, accessToken, isAuthenticated, login, logout }}>
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
        accessToken: auth.accessToken
    }
}
