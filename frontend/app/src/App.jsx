import './styles/App.css';
import React, { useState } from 'react';
import Header from './components/Header.jsx';
import TabGenerator from './components/TabGenerator.jsx';
import AudioSeparation from './components/AudioSeparation.jsx';
import SavedTabs from './components/SavedTabs.jsx';
import { LanguageProvider } from './i18n/useLanguage.jsx';

export default function App() {
    const [activeTab, setActiveTab] = useState('generator');

    return (
        <LanguageProvider>
            <div className="app">
                <Header activeTab={activeTab} setActiveTab={setActiveTab} />

                <main className="main-content">
                    {activeTab === 'generator' && <TabGenerator />}
                    {activeTab === 'separation' && <AudioSeparation />}
                    {activeTab === 'saved' && <SavedTabs />}
                </main>
            </div>
        </LanguageProvider>
    );
}