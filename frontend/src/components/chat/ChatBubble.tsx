import React, { useState } from 'react';
import ChatAssistant from './ChatAssistant';
import { cn } from '../../utils/cn';

const ChatBubble: React.FC = () => {
    const [isOpen, setIsOpen] = useState(false);

    return (
        <>
            {/* Floating Button */}
            <button
                onClick={() => setIsOpen(!isOpen)}
                className={cn(
                    'fixed bottom-6 right-6 p-4 rounded-full shadow-lg transition-all duration-300 z-20',
                    isOpen
                        ? 'bg-red-600 hover:bg-red-700 animate-in scale-100'
                        : 'bg-blue-600 hover:bg-blue-700 animate-in fade-in slide-in-from-bottom-5'
                )}
                aria-label="Open chat assistant"
            >
                {isOpen ? (
                    <svg className="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z" />
                    </svg>
                ) : (
                    <svg className="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm3.5-9c.83 0 1.5-.67 1.5-1.5S16.33 8 15.5 8 14 8.67 14 9.5s.67 1.5 1.5 1.5zm-7 0c.83 0 1.5-.67 1.5-1.5S9.33 8 8.5 8 7 8.67 7 9.5 7.67 11 8.5 11zm3.5 6.5c2.33 0 4.31-1.46 5.11-3.5H6.89c.8 2.04 2.78 3.5 5.11 3.5z" />
                    </svg>
                )}
            </button>

            {/* Chat Assistant Panel */}
            <ChatAssistant isOpen={isOpen} onClose={() => setIsOpen(false)} />
        </>
    );
};

export default ChatBubble;
