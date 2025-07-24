import React, { useState, useEffect } from 'react'
import AuthService from '../services/AuthService'
import { useNavigate } from 'react-router-dom'
import { config } from '../config'
import { useAuth } from '../modules/auth/AuthContext'
import { buttonClasses, inputClasses, formClasses, tagClasses } from '../utils/tailwindClasses'

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
        <div className="flex justify-center items-center min-h-screen p-8">
            <div className="bg-white/5 backdrop-blur-md p-8 rounded-xl border border-white/10 w-full max-w-md shadow-2xl">
                <div className="text-center mb-6">
                    <h2 className="text-center mb-2 text-2xl font-semibold">
                        Sign Up
                    </h2>
                    <div className="flex justify-center gap-2 mt-2">
                        <span className={tagClasses.success}>Free</span>
                        <span className={tagClasses.info}>Secure</span>
                    </div>
                </div>
                {errors.general && (
                    <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-md text-red-400 text-sm">
                        {errors.general}
                    </div>
                )}
                <form onSubmit={handleSubmit}>
                    <div className={formClasses.group}>
                        <label htmlFor="name" className={formClasses.label}>
                            Name
                        </label>
                        <input
                            type="text"
                            id="name"
                            name="name"
                            value={formData.name}
                            onChange={handleChange}
                            className={`w-full ${errors.name ? inputClasses.error : inputClasses.base}`}
                            disabled={isSubmitting}
                        />
                        {errors.name && (
                            <span className={formClasses.error}>{errors.name}</span>
                        )}
                    </div>

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
                                <span className={tagClasses.warning}>Creating account...</span>
                            </span>
                        ) : (
                            'Sign Up'
                        )}
                    </button>
                </form>

                <p className="text-center mt-6 text-sm">
                    Already have an account?{' '}
                    <a href="/login" className={tagClasses.primary}>
                        Login
                    </a>
                </p>
            </div>
        </div>
    )
}

export default Signup
