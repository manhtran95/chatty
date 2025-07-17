import './App.css'
import Home from './pages/Home'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import ChatApp from './apps/ChatApp/ChatApp'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/home" element={<Home />} />
        <Route path="/" element={<ChatApp />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
