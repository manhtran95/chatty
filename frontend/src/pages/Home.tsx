import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../modules/auth/AuthContext'

function Home() {
    const navigate = useNavigate()
    const { isAuthenticated } = useAuth()!

    // Redirect if already authenticated
    useEffect(() => {
        if (isAuthenticated) {
            navigate('/')
        }
    }, [isAuthenticated, navigate])

    // Don't render if already authenticated (will redirect)
    if (isAuthenticated) {
        return null
    }

    return (
        <div className="home">
            <h1>Welcome to Chatty!</h1>
            <p>This is the home page of your chat application.</p>
            <p>Explore the features and enjoy chatting with others!</p>
        </div>
    )
}

export default Home
