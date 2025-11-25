import { useState, useCallback } from 'react';

export interface ViaCepResponse {
    cep?: string;
    logradouro?: string;
    complemento?: string;
    bairro?: string;
    localidade?: string;
    uf?: string;
    ddd?: string;
    erro?: boolean;
}

export default function useCepLookup() {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const lookupCep = useCallback(async (rawCep?: string): Promise<ViaCepResponse | null> => {
        const cep = ((rawCep ?? '') || '').replace(/\D/g, '');
        if (cep.length !== 8) return null;

        setError(null);
        setLoading(true);
        try {
            const res = await fetch(`https://viacep.com.br/ws/${cep}/json/`);
            if (!res.ok) {
                setError('Erro na consulta do CEP');
                return null;
            }
            const data = (await res.json()) as ViaCepResponse;
            if (data.erro) {
                setError('CEP não encontrado');
                return null;
            }
            return data;
        } catch (err) {
            console.error('Error fetching CEP:', err);
            setError('Error in CEP lookup');
            return null;
        } finally {
            setLoading(false);
        }
    }, []);

    const clearError = useCallback(() => setError(null), []);

    return { lookupCep, loading, error, clearError } as const;
}
