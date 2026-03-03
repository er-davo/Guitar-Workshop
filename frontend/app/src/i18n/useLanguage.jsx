import React, { createContext, useContext, useState } from 'react';
import { translations } from './translations.js';

// Создаём контекст языка
const LanguageContext = createContext();

// Провайдер, оборачивает приложение
export function LanguageProvider({ children }) {
    const [lang, setLang] = useState('ru');

    // функция перевода
    const t = (key) => {
        return translations[lang]?.[key] || key;
    };

    return (
        <LanguageContext.Provider value={{ lang, setLang, t }}>
            {children}
        </LanguageContext.Provider>
    );
}

// хук для использования языка в компонентах
export function useLanguage() {
    return useContext(LanguageContext);
}