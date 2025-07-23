import { StatusCodes } from 'http-status-codes'
import { config } from '../config'

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

interface SignupResponse {
    form: {
        nonFieldErrors: string[]
        fieldErrors: { [key: string]: string }
    },
    redirect: boolean
}

interface SignupResult {
    success: boolean
    formData?: {
        nonFieldErrors: string[]
        fieldErrors: { [key: string]: string }
    }
    redirect?: string
}
interface LoginData {
    email: string
    password: string
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

export interface LoginResult {
    success: boolean
    data: {
        nonFieldErrors?: string[]
        userInfo?: UserInfo
        accessToken?: string
    }
}

interface RefreshResponse {
    accessToken: string
}

interface RefreshResult {
    success: boolean
    data: {
        accessToken: string
    }
}

interface LogoutResult {
    success: boolean
}

class AuthService {
    private baseURL: string

    constructor() {
        this.baseURL = config.api.baseUrl
    }

    async signup(userData: SignupData): Promise<SignupResult> {
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

            const _response: APIResponse<SignupResponse> =
                await response.json()
            return {
                success: false,
                formData: {
                    nonFieldErrors: _response.data.form.nonFieldErrors || [],
                    fieldErrors: _response.data.form.fieldErrors || {},
                },
            }
        } catch (error) {
            console.error('Signup error:', error)
            throw error
        }
    }

    async login(loginData: LoginData): Promise<LoginResult> {
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
                credentials: "include",
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

    async refresh(): Promise<RefreshResult> {
        const response = await fetch(`${this.baseURL}/user/refresh`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: "include",
        })
        const _response: APIResponse<RefreshResponse> = await response.json()
        if (response.status == StatusCodes.OK)
            return {
                success: true,
                data: {
                    accessToken: _response.data.accessToken,
                },
            }
        else
            return {
                success: false,
                data: {
                    accessToken: '',
                },
            }
    }

    async logout(): Promise<LogoutResult> {
        try {
            const response = await fetch(`${this.baseURL}/user/logout`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: "include",
            })

            if (response.status === StatusCodes.OK) {
                return { success: true }
            } else {
                return { success: false }
            }
        } catch (error) {
            console.error('Logout error:', error)
            return { success: false }
        }
    }
}

export default new AuthService()
