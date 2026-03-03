import '../styles/Header.css';
import React from 'react';
import { useLanguage } from '../i18n/useLanguage.jsx';

export default function LanguageSwitcher() {
    const { lang, setLang } = useLanguage();

    return (
        <div className="language-switcher">
            <button
                onClick={() => setLang('ru')}
                className={lang === 'ru' ? 'active' : ''}
            >
                RU
            </button>

            <button
                onClick={() => setLang('en')}
                className={lang === 'en' ? 'active' : ''}
            >
                EN
            </button>
        </div>
    );
}