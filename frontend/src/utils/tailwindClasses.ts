// Common Tailwind class combinations for reusability

export const buttonClasses = {
    primary: "px-5 py-2.5 bg-blue-600 text-white border-none rounded-md text-sm font-medium cursor-pointer transition-colors duration-150 hover:bg-blue-700 disabled:bg-gray-500 disabled:cursor-not-allowed",
    secondary: "px-5 py-2.5 bg-gray-600 text-white border-none rounded-md text-sm font-medium cursor-pointer transition-colors duration-150 hover:bg-gray-700 disabled:bg-gray-400 disabled:cursor-not-allowed",
    danger: "px-5 py-2.5 bg-red-600 text-white border-none rounded-md text-sm font-medium cursor-pointer transition-colors duration-150 hover:bg-red-700 disabled:bg-red-400 disabled:cursor-not-allowed",
    small: "px-3 py-2 bg-blue-600 text-white border-none rounded-md text-sm font-medium cursor-pointer transition-colors duration-150 hover:bg-blue-700 disabled:bg-gray-500 disabled:cursor-not-allowed",
}

export const inputClasses = {
    base: "px-3 py-2 border rounded-md text-sm transition-colors duration-150 focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-200",
    error: "px-3 py-2 border border-red-500 rounded-md text-sm transition-colors duration-150 focus:outline-none focus:border-red-500 focus:ring-2 focus:ring-red-200",
    disabled: "px-3 py-2 border border-gray-300 rounded-md text-sm bg-gray-100 cursor-not-allowed",
}

export const modalClasses = {
    overlay: "fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50",
    content: "bg-white rounded-lg shadow-xl w-11/12 max-w-md max-h-[90vh] overflow-y-auto",
    header: "flex justify-between items-center p-5 pb-0 border-b border-gray-200",
    form: "p-5",
    actions: "flex gap-3 justify-end mt-6",
}

export const layoutClasses = {
    flexCenter: "flex items-center justify-center",
    flexBetween: "flex justify-between items-center",
    flexCol: "flex flex-col",
    container: "w-full h-full",
}

export const statusClasses = {
    connected: "px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800",
    disconnected: "px-3 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800",
}

export const chatClasses = {
    item: "p-3 border-b border-gray-200 cursor-pointer text-left transition-colors duration-150",
    itemSelected: "p-3 border-b border-gray-200 cursor-pointer text-left transition-colors duration-150 bg-blue-50",
    itemHover: "p-3 border-b border-gray-200 cursor-pointer text-left transition-colors duration-150 bg-white hover:bg-gray-50",
    header: "p-4 border-b border-gray-200 bg-gray-50",
    list: "flex-1 overflow-y-auto bg-white",
}

export const formClasses = {
    group: "mb-5",
    label: "block mb-2 font-medium text-gray-700 text-sm",
    error: "text-red-500 text-xs mt-1 block",
}

export const tagClasses = {
    base: "inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium",
    primary: "inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800",
    secondary: "inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800",
    success: "inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800",
    warning: "inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800",
    danger: "inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800",
    info: "inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-cyan-100 text-cyan-800",
    purple: "inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-800",
    removable: "inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800 hover:bg-blue-200 transition-colors duration-150",
    closeButton: "ml-1.5 -mr-1 h-4 w-4 rounded-full inline-flex items-center justify-center text-blue-400 hover:bg-blue-200 hover:text-blue-500 focus:outline-none focus:bg-blue-200 focus:text-blue-500",
}

export const messageClasses = {
    container: "p-2 border-b border-gray-200 text-left",
    sender: "font-bold mb-1 text-gray-900",
    content: "text-gray-700",
    timestamp: "text-xs text-gray-500 mt-1",
    ownMessage: "p-2 border-b border-gray-200 text-right",
} 