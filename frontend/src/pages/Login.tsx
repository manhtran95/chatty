import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import './Signup.css'
import { useAuth } from '../modules/auth/AuthContext'

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
        <div className="signup-container">
            <div className="signup-form">
                <h2>Login</h2>
                {errors.general && (
                    <div className="error-message general-error">
                        {errors.general}
                    </div>
                )}
                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="email">Email</label>
                        <input
                            type="email"
                            id="email"
                            name="email"
                            value={formData.email}
                            onChange={handleChange}
                            className={errors.email ? 'error' : ''}
                            disabled={isSubmitting}
                        />
                        {errors.email && (
                            <span className="error-message">
                                {errors.email}
                            </span>
                        )}
                    </div>

                    <div className="form-group">
                        <label htmlFor="password">Password</label>
                        <input
                            type="password"
                            id="password"
                            name="password"
                            value={formData.password}
                            onChange={handleChange}
                            className={errors.password ? 'error' : ''}
                            disabled={isSubmitting}
                        />
                        {errors.password && (
                            <span className="error-message">
                                {errors.password}
                            </span>
                        )}
                    </div>

                    <button
                        type="submit"
                        className="btn-primary"
                        disabled={isSubmitting}
                        style={{ width: '100%', marginTop: '1rem' }}
                    >
                        {isSubmitting ? 'Logging in...' : 'Login'}
                    </button>
                </form>

                <p>
                    Don&apos;t have an account?{' '}
                    <a href="/signup" className="link-primary">
                        Sign up
                    </a>
                </p>
            </div>
        </div>
    )
}

export default Login
