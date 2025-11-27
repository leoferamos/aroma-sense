import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

const AuthBackLink: React.FC = () => {
    const { isAuthenticated } = useAuth();
    const target = isAuthenticated ? '/products' : '/register';
    return (
        <Link
            to={target}
            className="block w-full md:w-1/2 mx-auto text-white text-lg font-medium py-3 rounded-full mt-2 bg-blue-600 hover:bg-blue-700 text-center transition-colors"
        >
            Voltar
        </Link>
    );
};

export default AuthBackLink;