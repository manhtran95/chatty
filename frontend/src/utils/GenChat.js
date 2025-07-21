function generateUUID() {
    // Basic UUID v4 generator
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
        const r = (Math.random() * 16) | 0
        const v = c === 'x' ? r : (r & 0x3) | 0x8
        return v.toString(16)
    })
}

function getRandomUsername() {
    const names = [
        'alice',
        'bob',
        'carol',
        'dave',
        'eve',
        'frank',
        'grace',
        'heidi',
    ]
    return (
        names[Math.floor(Math.random() * names.length)] +
        Math.floor(Math.random() * 100)
    )
}

function getRandomMessage() {
    const messages = [
        'Hey!',
        "How's it going?",
        "What's up?",
        "Let's meet tomorrow.",
        'Got it.',
        'Thanks!',
        'Good morning!',
        'Check this out.',
        'Any updates?',
        'See you soon.',
    ]
    return messages[Math.floor(Math.random() * messages.length)]
}

function generateGroupChat() {
    const userCount = Math.floor(Math.random() * 3) + 2 // 2 to 4 users
    const userList = []
    const usernameSet = new Set()

    while (userList.length < userCount) {
        const username = getRandomUsername()
        if (!usernameSet.has(username)) {
            userList.push({ id: generateUUID(), username })
            usernameSet.add(username)
        }
    }

    const messages = Array.from({ length: 10 }, () => {
        const sender = userList[Math.floor(Math.random() * userList.length)]
        return {
            senderName: sender.username,
            content: getRandomMessage(),
        }
    })

    return { userList, messages }
}

const groups = Array.from({ length: 5 }, generateGroupChat)

console.log(JSON.stringify(groups, null, 2))
