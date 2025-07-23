import React, { useState, useEffect } from 'react'
import AuthService from '../services/AuthService'
import './Signup.css'
import { useNavigate } from 'react-router-dom'
import { config } from '../config'
import { useAuth } from '../modules/auth/AuthContext'

interface FormData {
    name: string
    email: string
    password: string
}

interface FormErrors {
    name?: string
    email?: string
    password?: string
    general?: string
}

const Signup: React.FC = () => {
    const navigate = useNavigate()
    const { isAuthenticated } = useAuth()!
    const [formData, setFormData] = useState<FormData>({
        name: '',
        email: '',
        password: '',
    })

    const [errors, setErrors] = useState<FormErrors>({})
    const [isSubmitting, setIsSubmitting] = useState(false)

    // Redirect if already authenticated
    useEffect(() => {
        if (isAuthenticated) {
            navigate('/')
        }
    }, [isAuthenticated, navigate])

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

    const validateForm = (): boolean => {
        const newErrors: FormErrors = {}

        if (!formData.name.trim()) {
            newErrors.name = 'Name is required'
        }

        if (!formData.email.trim()) {
            newErrors.email = 'Email is required'
        } else if (!emailRegex.test(formData.email)) {
            newErrors.email = 'Please enter a valid email address'
        }

        if (!formData.password) {
            newErrors.password = 'Password is required'
        } else if (formData.password.length < config.auth.passwordMinLength) {
            newErrors.password = `Password must be at least ${config.auth.passwordMinLength} characters long`
        } else if (formData.password.length > config.auth.passwordMaxLength) {
            newErrors.password = `Password must be no more than ${config.auth.passwordMaxLength} characters long`
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

        if (errors[name as keyof FormErrors]) {
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
            const result = await AuthService.signup(formData)
            console.log('Signup result:', result)

            if (result.success) {
                console.log('Success, Redirecting to:', result.redirect)
                navigate(result.redirect || '/login')
            } else {
                if (result.formData?.nonFieldErrors) {
                    setErrors((prev) => ({
                        ...prev,
                        general: result.formData?.nonFieldErrors.join(', '),
                    }))
                }

                if (result.formData?.fieldErrors) {
                    setErrors((prev) => ({
                        ...prev,
                        ...result.formData?.fieldErrors,
                    }))
                }
            }
        } catch (error) {
            console.error('Signup failed:', error)
            setErrors((prev) => ({
                ...prev,
                general: 'Signup failed. Please try again.',
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
                <h2>Sign Up</h2>
                {errors.general && (
                    <div className="error-message general-error">
                        {errors.general}
                    </div>
                )}
                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="name">Name</label>
                        <input
                            type="text"
                            id="name"
                            name="name"
                            value={formData.name}
                            onChange={handleChange}
                            className={errors.name ? 'error' : ''}
                            disabled={isSubmitting}
                        />
                        {errors.name && (
                            <span className="error-message">{errors.name}</span>
                        )}
                    </div>

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
                        {isSubmitting ? 'Signing up...' : 'Sign Up'}
                    </button>
                </form>

                <p>
                    Already have an account?{' '}
                    <a href="/login" className="link-primary">
                        Login
                    </a>
                </p>
            </div>
        </div>
    )
}

export default Signup
