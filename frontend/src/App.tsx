import './App.css'
import Home from './pages/Home'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import ChatApp from './apps/ChatApp/ChatApp'
import Signup from './pages/Signup'
import Login from './pages/Login'
import { AuthProvider } from './modules/auth/AuthContext'
import { WebSocketProvider } from './modules/websocket/WebSocketContext'

function App() {
    return (
        <BrowserRouter>
            <AuthProvider>
                <WebSocketProvider>
                    <Routes>
                        <Route path="/home" element={<Home />} />
                        <Route path="/signup" element={<Signup />} />
                        <Route path="/login" element={<Login />} />
                        <Route path="/" element={<ChatApp />} />
                    </Routes>
                </WebSocketProvider>
            </AuthProvider>
        </BrowserRouter>
    )
}

export default App
