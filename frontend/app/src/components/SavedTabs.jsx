import '../styles/SavedTabs.css';
import '../styles/TabGenerator.css';

import React, { useState } from 'react';
import { searchTabs, getTab, deleteTab } from '../api/api.js';
import { useLanguage } from '../i18n/useLanguage.jsx';

export default function SavedTabs() {
    const { t } = useLanguage();

    const [query, setQuery] = useState('');
    const [results, setResults] = useState([]);
    const [selectedTab, setSelectedTab] = useState(null);
    const [tabText, setTabText] = useState('');

    const handleSearch = async () => {
        if (!query.trim()) return alert(t('searchPlaceholder'));

        try {
            const tabs = await searchTabs(query);
            setResults(tabs);
            setSelectedTab(null);
            setTabText('');
        } catch (err) {
            alert(`${t('errorGeneric')}: ${err.message}`);
            setResults([]);
        }
    };

    const handleSelect = async (tab) => {
        setSelectedTab(tab);
        setTabText('');

        try {
            const data = await getTab(tab.id);

            if (!data?.presigned_url) {
                throw new Error(t('errorNoFileLink'));
            }

            const resp = await fetch(data.presigned_url);
            if (!resp.ok) {
                throw new Error(t('errorFileDownload'));
            }

            const text = await resp.text();
            setTabText(text);

        } catch (err) {
            alert(err.message);
            setSelectedTab(null);
            setTabText('');
        }
    };

    const handleDelete = async (tab) => {
        if (!window.confirm(t('confirmDelete'))) return;

        try {
            await deleteTab(tab.id);

            setResults(prev => prev.filter(r => r.id !== tab.id));

            if (selectedTab?.id === tab.id) {
                setSelectedTab(null);
                setTabText('');
            }
        } catch (err) {
            alert(`${t('errorGeneric')}: ${err.message}`);
        }
    };

    const handleClose = () => {
        setSelectedTab(null);
        setTabText('');
    };

    return (
        <div className="saved-tabs">
            <h2>{t('savedTabsTitle')}</h2>

            <div className="search-bar">
                <input
                    type="text"
                    placeholder={t('searchPlaceholder')}
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                />
                <button onClick={handleSearch}>
                    {t('search')}
                </button>
            </div>

            <div className="tabs-list">
                {results.map(tab => (
                    <div
                        key={tab.id}
                        className={`tab-item ${selectedTab?.id === tab.id ? 'selected' : ''}`}
                    >
                        <span>{tab.name}</span>

                        <button onClick={() => handleSelect(tab)}>
                            {t('open')}
                        </button>

                        <button onClick={() => handleDelete(tab)}>
                            {t('delete')}
                        </button>
                    </div>
                ))}
            </div>

            {selectedTab && tabText && (
                <div className="tab-result">
                    <div className="tab-header">
                        <h3>{selectedTab.name}</h3>
                        <button onClick={handleClose}>
                            {t('close')}
                        </button>
                    </div>

                    <pre>{tabText}</pre>
                </div>
            )}
        </div>
    );
}