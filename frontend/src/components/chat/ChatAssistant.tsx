import React, { useState, useRef, useEffect } from 'react';
import { cn } from '../../utils/cn';

interface Message {
    id: string;
    content: string;
    sender: 'user' | 'assistant';
    timestamp: Date;
}

interface ChatAssistantProps {
    isOpen: boolean;
    onClose: () => void;
}

const ChatAssistant: React.FC<ChatAssistantProps> = ({ isOpen, onClose }) => {
    const [messages, setMessages] = useState<Message[]>([
        {
            id: '1',
            content: 'Hello! I\'m your Aroma Sense assistant. How can I help you today?',
            sender: 'assistant',
            timestamp: new Date(),
        },
    ]);
    const [inputValue, setInputValue] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(() => {
        scrollToBottom();
    }, [messages]);


    useEffect(() => {
        if (isOpen) {
            const t = setTimeout(() => scrollToBottom(), 50);
            return () => clearTimeout(t);
        }
    }, [isOpen]);

    const handleSendMessage = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!inputValue.trim()) return;

        // Add user message
        const userMessage: Message = {
            id: Date.now().toString(),
            content: inputValue,
            sender: 'user',
            timestamp: new Date(),
        };

        setMessages((prev) => [...prev, userMessage]);
        setInputValue('');

    // TODO: Call AI API here

        // Simulate assistant response
        setIsLoading(true);
        setTimeout(() => {
            const assistantMessage: Message = {
                id: (Date.now() + 1).toString(),
                content: 'I received your message. This is a placeholder response until the AI integration is implemented.',
                sender: 'assistant',
                timestamp: new Date(),
            };
            setMessages((prev) => [...prev, assistantMessage]);
            setIsLoading(false);
        }, 1000);
    };

    if (!isOpen) return null;

    return (
        <>
            {/* Overlay */}
            <div
                className="fixed inset-0 bg-black/50 z-30 animate-in fade-in"
                onClick={onClose}
            />

            {/* Chat Panel (bottom sheet on mobile, side panel on desktop) */}
            <div
                className={cn(
                    "fixed inset-x-0 bottom-0 md:inset-y-0 md:right-0 md:left-auto w-full md:w-96 h-[70dvh] md:h-[100dvh] bg-white shadow-2xl z-40 animate-in md:slide-in-from-right slide-in-from-bottom flex flex-col overflow-hidden min-h-0 rounded-t-2xl md:rounded-none border-t md:border-t-0"
                )}
            >
                {/* Header */}
                <div className="flex items-center justify-between bg-gradient-to-r from-blue-600 to-blue-700 text-white p-4 pt-[calc(theme(spacing.4)+env(safe-area-inset-top))] md:pt-4 rounded-t-2xl md:rounded-tl-lg">
                    <div className="flex items-center gap-3">
                        <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm3.5-9c.83 0 1.5-.67 1.5-1.5S16.33 8 15.5 8 14 8.67 14 9.5s.67 1.5 1.5 1.5zm-7 0c.83 0 1.5-.67 1.5-1.5S9.33 8 8.5 8 7 8.67 7 9.5 7.67 11 8.5 11zm3.5 6.5c2.33 0 4.31-1.46 5.11-3.5H6.89c.8 2.04 2.78 3.5 5.11 3.5z" />
                        </svg>
                        <h2 className="text-lg font-semibold">Aroma Sense Assistant</h2>
                    </div>
                    <button
                        onClick={onClose}
                        className="p-1 hover:bg-blue-500 rounded-md transition-colors"
                        aria-label="Close chat"
                    >
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z" />
                        </svg>
                    </button>
                </div>

                {/* Messages Container */}
                <div className="flex-1 min-h-0 overflow-y-auto p-4 space-y-4 scroll-pb-24 overscroll-contain">
                    {messages.map((message) => (
                        <div
                            key={message.id}
                            className={cn(
                                'flex',
                                message.sender === 'user' ? 'justify-end' : 'justify-start'
                            )}
                        >
                            <div
                                className={cn(
                                    'max-w-[80%] sm:max-w-xs px-4 py-2 rounded-lg',
                                    message.sender === 'user'
                                        ? 'bg-blue-600 text-white rounded-br-none'
                                        : 'bg-gray-200 text-gray-900 rounded-bl-none'
                                )}
                            >
                                <p className="text-sm">{message.content}</p>
                                <span className="text-xs opacity-70 mt-1 block">
                                    {message.timestamp.toLocaleTimeString([], {
                                        hour: '2-digit',
                                        minute: '2-digit',
                                    })}
                                </span>
                            </div>
                        </div>
                    ))}
                    {isLoading && (
                        <div className="flex justify-start">
                            <div className="bg-gray-200 text-gray-900 px-4 py-2 rounded-lg rounded-bl-none">
                                <div className="flex gap-1">
                                    <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce" />
                                    <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce delay-100" />
                                    <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce delay-200" />
                                </div>
                            </div>
                        </div>
                    )}
                    <div ref={messagesEndRef} />
                </div>

                {/* Input Area */}
                <form
                    onSubmit={handleSendMessage}
                    className="border-t border-gray-200 p-4 pb-[calc(env(safe-area-inset-bottom)+0.5rem)] bg-gray-50"
                >
                    <div className="flex gap-2">
                        <input
                            type="text"
                            value={inputValue}
                            onChange={(e) => setInputValue(e.target.value)}
                            placeholder="Type your message..."
                            className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-[15px]"
                            disabled={isLoading}
                        />
                        <button
                            type="submit"
                            disabled={isLoading || !inputValue.trim()}
                            className={cn(
                                'p-2 rounded-md transition-colors',
                                isLoading || !inputValue.trim()
                                    ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                                    : 'bg-blue-600 text-white hover:bg-blue-700'
                            )}
                            aria-label="Send message"
                        >
                            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                                <path d="M16.6915026,12.4744748 L3.50612381,13.2599618 C3.19218622,13.2599618 3.03521743,13.4170592 3.03521743,13.5741566 L1.15159189,20.0151496 C0.8376543,20.8006365 0.99,21.89 1.77946707,22.52 C2.41,22.99 3.50612381,23.1 4.13399899,22.8429026 L21.714504,14.0454487 C22.6563168,13.5741566 23.1272231,12.6315722 22.9702544,11.6889879 L4.13399899,1.16141721 C3.34915502,0.9 2.40734225,0.9 1.77946707,1.42274535 C0.994623095,2.10604706 0.837654326,3.0486314 1.15159189,3.99701575 L3.03521743,10.4380088 C3.03521743,10.5951061 3.34915502,10.7522035 3.50612381,10.7522035 L16.6915026,11.5376905 C16.6915026,11.5376905 17.1624089,11.5376905 17.1624089,12.0089827 C17.1624089,12.4744748 16.6915026,12.4744748 16.6915026,12.4744748 Z" />
                            </svg>
                        </button>
                    </div>
                </form>
            </div>
        </>
    );
};

export default ChatAssistant;
