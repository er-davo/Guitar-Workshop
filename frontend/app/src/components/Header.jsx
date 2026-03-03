import '../styles/Header.css';
import React from 'react';
import LanguageSwitcher from './LanguageSwitcher.jsx';
import { useLanguage } from '../i18n/useLanguage.jsx';

export default function Header({ activeTab, setActiveTab }) {
    const { t } = useLanguage();

    return (
        <header className="app-header">
            <h1>{t('appTitle')}</h1>

            <nav className="tab-nav">
                <button
                    onClick={() => setActiveTab('generator')}
                    className={activeTab === 'generator' ? 'active' : ''}
                >
                    {t('tabGenTitle')}
                </button>

                <button
                    onClick={() => setActiveTab('separation')}
                    className={activeTab === 'separation' ? 'active' : ''}
                >
                    {t('audioSepTitle')}
                </button>

                <button
                    onClick={() => setActiveTab('saved')}
                    className={activeTab === 'saved' ? 'active' : ''}
                >
                    {t('savedTabsTitle')}
                </button>
            </nav>

            <LanguageSwitcher />
        </header>
    );
}