/**
 * Entry point for the React application.
 * Renders the App component inside the root element.
 */
import { createRoot } from 'react-dom/client';
import App from './App';
import './index.css';
import './i18n';

createRoot(document.getElementById('root')!).render(
    <App />
);
