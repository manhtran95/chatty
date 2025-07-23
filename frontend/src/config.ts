// Configuration file for the application
export const config = {
    // Authentication settings
    auth: {
        passwordMinLength: parseInt(import.meta.env.VITE_PASSWORD_MIN_LENGTH || '8'),
        passwordMaxLength: parseInt(import.meta.env.VITE_PASSWORD_MAX_LENGTH || '128'),
    },

    // API settings
    api: {
        baseUrl: import.meta.env.VITE_API_URL,
        timeout: parseInt(import.meta.env.VITE_API_TIMEOUT || '5000'),
    },

    // App settings
    app: {
        name: import.meta.env.VITE_APP_NAME || 'Chatty',
        version: import.meta.env.VITE_APP_VERSION || '1.0.0',
    },
} as const

// Type for the config object
export type Config = typeof config 