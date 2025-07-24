import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../modules/auth/AuthContext'
import { buttonClasses, inputClasses, formClasses, tagClasses } from '../utils/tailwindClasses'

interface LoginData {
    email: string
    password: string
}

interface LoginErrors {
    email?: string
    password?: string
    general?: string
}

const Login: React.FC = () => {
    const { login, isAuthenticated } = useAuth()!
    const navigate = useNavigate()
    const [formData, setFormData] = useState<LoginData>({
        email: '',
        password: '',
    })

    const [errors, setErrors] = useState<LoginErrors>({})
    const [isSubmitting, setIsSubmitting] = useState(false)

    // Redirect if already authenticated
    useEffect(() => {
        if (isAuthenticated) {
            navigate('/')
        }
    }, [isAuthenticated, navigate])

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

    const validateForm = (): boolean => {
        const newErrors: LoginErrors = {}

        if (!formData.email.trim()) {
            newErrors.email = 'Email is required'
        } else if (!emailRegex.test(formData.email)) {
            newErrors.email = 'Please enter a valid email address'
        }

        if (!formData.password) {
            newErrors.password = 'Password is required'
        }

        setErrors(newErrors)
        return Object.keys(newErrors).length === 0
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }))

        if (errors[name as keyof LoginErrors]) {
            setErrors((prev) => ({
                ...prev,
                [name]: undefined,
            }))
        }
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        if (!validateForm()) {
            return
        }

        setIsSubmitting(true)

        try {
            const result = await login(formData.email, formData.password)
            console.log(`Login result:`)
            console.log(result)
            if (result.success) {
                // Login successful - user and access token are set in AuthContext
                console.log('Login successful')
                // Redirect to ChatApp
                navigate('/')
            } else {
                // Handle login failure - result contains error information
                console.log('Login failed:', result.data)
                if (result.data.nonFieldErrors) {
                    setErrors((prev) => ({
                        ...prev,
                        general: result.data.nonFieldErrors!.join(', '),
                    }))
                }
            }
        } catch (error) {
            console.error('Login failed:', error)
            setErrors((prev) => ({
                ...prev,
                general: 'Login failed. Please try again.',
            }))
        } finally {
            setIsSubmitting(false)
        }
    }

    // Don't render if already authenticated (will redirect)
    if (isAuthenticated) {
        return null
    }

    return (
        <div className="flex justify-center items-center min-h-screen p-8">
            <div className="bg-white/5 backdrop-blur-md p-8 rounded-xl border border-white/10 w-full max-w-md shadow-2xl">
                <div className="text-center mb-6">
                    <h2 className="text-center mb-2 text-2xl font-semibold">
                        Login
                    </h2>
                    <span className={tagClasses.info}>Welcome back!</span>
                </div>
                {errors.general && (
                    <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-md text-red-400 text-sm">
                        {errors.general}
                    </div>
                )}
                <form onSubmit={handleSubmit}>
                    <div className={formClasses.group}>
                        <label htmlFor="email" className={formClasses.label}>
                            Email
                        </label>
                        <input
                            type="email"
                            id="email"
                            name="email"
                            value={formData.email}
                            onChange={handleChange}
                            className={`w-full ${errors.email ? inputClasses.error : inputClasses.base}`}
                            disabled={isSubmitting}
                        />
                        {errors.email && (
                            <span className={formClasses.error}>
                                {errors.email}
                            </span>
                        )}
                    </div>

                    <div className={formClasses.group}>
                        <label htmlFor="password" className={formClasses.label}>
                            Password
                        </label>
                        <input
                            type="password"
                            id="password"
                            name="password"
                            value={formData.password}
                            onChange={handleChange}
                            className={`w-full ${errors.password ? inputClasses.error : inputClasses.base}`}
                            disabled={isSubmitting}
                        />
                        {errors.password && (
                            <span className={formClasses.error}>
                                {errors.password}
                            </span>
                        )}
                    </div>

                    <button
                        type="submit"
                        className={`w-full mt-6 ${buttonClasses.primary}`}
                        disabled={isSubmitting}
                    >
                        {isSubmitting ? (
                            <span className="flex items-center justify-center gap-2">
                                <span className={tagClasses.warning}>Processing...</span>
                            </span>
                        ) : (
                            'Login'
                        )}
                    </button>
                </form>

                <p className="text-center mt-6 text-sm">
                    Don&apos;t have an account?{' '}
                    <a href="/signup" className={tagClasses.primary}>
                        Sign up
                    </a>
                </p>
            </div>
        </div>
    )
}

export default Login
