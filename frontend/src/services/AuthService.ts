interface APIResponse<T> {
    data?: T
    message?: string
    status?: string
}

interface SignupData {
    name: string
    email: string
    password: string
}

interface FormData {
    id: string
    name: string
    email: string
    NonFieldErrors: string[]
    FieldErrors: { [key: string]: string }
}

type SignupResult =
    | { success: true; redirect: string }
    | { success: false; formData: FormData }

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

            console.log('signupPost Response:', response)

            if (response.status == 200) {
                return { success: true, redirect: '/login' }
            }

            const _response: APIResponse<FormData> = await response.json()
            const data: FormData = _response.data || {
                id: '',
                name: '',
                email: '',
                NonFieldErrors: [],
                FieldErrors: {},
            }

            return { success: false, formData: data }
        } catch (error) {
            console.error('Signup error:', error)
            throw error
        }
    }
}

export default new AuthService()
