import { useState } from 'react'
import { useWebSocket } from '../../modules/websocket/WebSocketContext'
import { buttonClasses, inputClasses, modalClasses, formClasses } from '../../utils/tailwindClasses'

interface CreateChatModalProps {
    isOpen: boolean
    onClose: () => void
}

export default function CreateChatModal({ isOpen, onClose }: CreateChatModalProps) {
    const [recipientEmail, setRecipientEmail] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [error, setError] = useState('')
    const { createChat, isConnected } = useWebSocket()

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        if (!recipientEmail.trim()) {
            setError('Please enter a recipient email')
            return
        }

        if (!isConnected) {
            setError('Not connected to server')
            return
        }

        setIsSubmitting(true)
        setError('')

        try {
            const success = createChat('New Chat', [recipientEmail.trim()])

            if (success) {
                setRecipientEmail('')
                onClose()
            } else {
                setError('Failed to create chat. Please try again.')
            }
        } catch (error) {
            setError('An error occurred while creating the chat.')
        } finally {
            setIsSubmitting(false)
        }
    }

    const handleClose = () => {
        setRecipientEmail('')
        setError('')
        onClose()
    }

    if (!isOpen) return null

    return (
        <div className={modalClasses.overlay} onClick={handleClose}>
            <div className={modalClasses.content} onClick={(e) => e.stopPropagation()}>
                <div className={modalClasses.header}>
                    <h3 className="text-lg font-semibold text-gray-900 m-0">Create New Chat</h3>
                </div>

                <form onSubmit={handleSubmit} className={modalClasses.form}>
                    <div className={formClasses.group}>
                        <label htmlFor="recipientEmail" className={formClasses.label}>
                            Recipient Email
                        </label>
                        <input
                            type="email"
                            id="recipientEmail"
                            value={recipientEmail}
                            onChange={(e) => setRecipientEmail(e.target.value)}
                            placeholder="Enter recipient's email"
                            disabled={isSubmitting}
                            className={`w-11/12 ${isSubmitting ? inputClasses.disabled : error ? inputClasses.error : inputClasses.base}`}
                        />
                        {error && <span className={formClasses.error}>{error}</span>}
                    </div>

                    <div className={modalClasses.actions}>
                        <button
                            type="button"
                            onClick={handleClose}
                            className={buttonClasses.secondary}
                            disabled={isSubmitting}
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className={buttonClasses.primary}
                            disabled={isSubmitting || !recipientEmail.trim() || !isConnected}
                        >
                            {isSubmitting ? 'Creating...' : 'Create Chat'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
} 