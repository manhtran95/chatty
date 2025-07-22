import { StatusCodes } from 'http-status-codes'

interface APIResponse<T> {
    data: T
    message: string
    status: string
}

interface SignupData {
    name: string
    email: string
    password: string
}

interface LoginData {
    email: string
    password: string
}

// from backend, Capitalize the first letter of each key
interface SignupResponseData {
    NonFieldErrors: string[]
    FieldErrors: { [key: string]: string }
}
interface UserInfo {
    id: string
    name: string
    email: string
}

interface LoginResponse {
    nonFieldErrors?: string[]
    fieldErrors?: { [key: string]: string }
    userInfo?: UserInfo
    accessToken?: string
}

type SignupResult =
    | { success: true; redirect: string }
    | {
          success: false
          formData: {
              nonFieldErrors: string[]
              fieldErrors: { [key: string]: string }
          }
      }

export interface LoginResult {
    success: boolean
    data: {
        nonFieldErrors?: string[]
        userInfo?: UserInfo
        accessToken?: string
    }
}

class AuthService {
    private baseURL: string

    constructor() {
        this.baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    }

    async signupPost(userData: SignupData): Promise<SignupResult> {
        try {
            const formData = new URLSearchParams()
            Object.entries(userData).forEach(([key, value]) => {
                formData.append(key, value)
            })
            const response = await fetch(`${this.baseURL}/user/signup`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: formData.toString(),
            })

            if (response.status == 200) {
                return { success: true, redirect: '/login' }
            }

            const _response: APIResponse<SignupResponseData> =
                await response.json()
            return {
                success: false,
                formData: {
                    nonFieldErrors: _response.data.NonFieldErrors || [],
                    fieldErrors: _response.data.FieldErrors || {},
                },
            }
        } catch (error) {
            console.error('Signup error:', error)
            throw error
        }
    }

    async loginPost(loginData: LoginData): Promise<LoginResult> {
        try {
            const formData = new URLSearchParams()
            Object.entries(loginData).forEach(([key, value]) => {
                formData.append(key, value)
            })
            const response = await fetch(`${this.baseURL}/user/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: formData.toString(),
            })

            const _response: APIResponse<LoginResponse> = await response.json()
            if (response.status == StatusCodes.OK)
                return {
                    success: true,
                    data: {
                        userInfo: _response.data.userInfo,
                        accessToken: _response.data.accessToken,
                    },
                }
            else if (response.status == StatusCodes.UNPROCESSABLE_ENTITY)
                return {
                    success: false,
                    data: {
                        nonFieldErrors: _response.data.nonFieldErrors,
                    },
                }
            else
                return {
                    success: false,
                    data: {},
                }
        } catch (error) {
            console.error('Login error:', error)
            throw error
        }
    }
}

export default new AuthService()
