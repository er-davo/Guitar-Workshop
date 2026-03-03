import '../styles/TabGenerator.css';
import React, { useState, useCallback } from 'react';
import { createTask, getTask, saveTab } from '../api/api.js';
import { useLanguage } from '../i18n/useLanguage.jsx';
import { usePolling } from '../hooks/usePolling.js';

export default function TabGenerator() {
    const { t } = useLanguage();

    const [file, setFile] = useState(null);
    const [dragActive, setDragActive] = useState(false);
    const [separation, setSeparation] = useState(false);

    const [loading, setLoading] = useState(false);
    const [taskId, setTaskId] = useState(null);
    const [tabText, setTabText] = useState('');
    const [tabName, setTabName] = useState('');
    const [saveDisabled, setSaveDisabled] = useState(true);
    const [status, setStatus] = useState(null);

    const statusMap = {
        created: t('statusCreated'),
        pending: t('statusPending'),
        waiting_for_separation: t('statusWaitingForSeparation'),
        done: t('statusDone'),
        fail: t('statusFail'),
    };

    // ---------------- File Handlers ----------------
    const handleFileChange = (e) => {
        if (e.target.files?.[0]) {
            setFile(e.target.files[0]);
        }
    };

    const handleDrag = useCallback((e) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.type === 'dragenter' || e.type === 'dragover') {
            setDragActive(true);
        } else if (e.type === 'dragleave') {
            setDragActive(false);
        }
    }, []);

    const handleDrop = useCallback((e) => {
        e.preventDefault();
        e.stopPropagation();
        setDragActive(false);
        if (e.dataTransfer.files && e.dataTransfer.files[0]) {
            setFile(e.dataTransfer.files[0]);
        }
    }, []);

    // ---------------- Submit ----------------
    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!file) return alert(t('errorNoFile'));

        setLoading(true);
        setTabText('');
        setTaskId(null);
        setStatus(t('statusCreated'));
        setSaveDisabled(true);

        try {
            const data = await createTask(file, separation);
            const id = data.task.id;
            if (!id) throw new Error(t('errorNoTaskId'));
            setTaskId(id);
        } catch (err) {
            setTabText(`${t('errorGeneric')}: ${err.message}`);
            setStatus(null);
        } finally {
            setLoading(false);
        }
    };

    // ---------------- Polling ----------------
    usePolling(
        async () => {
            if (!taskId) return null;
            return await getTask(taskId);
        },
        8000,
        Infinity,
        async (result) => {
            if (!result?.task?.status) return false;

            const s = result.task.status.toLowerCase();
            setStatus(statusMap[s] || s);

            if (s === 'done' && result.tab?.presigned_url) {
                try {
                    const r = await fetch(result.tab.presigned_url);
                    if (!r.ok) throw new Error(r.statusText);
                    const text = await r.text();
                    setTabText(text);
                    setSaveDisabled(false);
                } catch {
                    setTabText(t('errorLoadingTab'));
                }
                return true;
            }

            if (s === 'fail') {
                setTabText(t('errorGeneratingTab'));
                setSaveDisabled(true);
                return true;
            }

            return false;
        },
        [taskId]
    );

    // ---------------- Actions ----------------
    const handleSave = async () => {
        if (!tabName.trim()) return alert(t('tabNamePlaceholder'));
        if (!taskId) return alert(t('errorNoTaskId'));

        await saveTab(taskId, tabName.trim());
        setTabName('');
        setSaveDisabled(true);
        alert(t('tabSaved'));
    };

    const handleCopy = () => {
        navigator.clipboard.writeText(tabText);
    };

    return (
        <div className="tab-generator">
            <h2>{t('tabGenTitle')}</h2>

            <form onSubmit={handleSubmit}>
                <div
                    className={`drop-zone ${dragActive ? 'active' : ''}`}
                    onDragEnter={handleDrag}
                    onDragLeave={handleDrag}
                    onDragOver={handleDrag}
                    onDrop={handleDrop}
                >
                    <input
                        type="file"
                        accept="audio/*"
                        onChange={handleFileChange}
                    />
                    <p>
                        {file
                            ? `${t('fileLabel')}: ${file.name}`
                            : t('dropzonePlaceholder')}
                    </p>
                </div>

                <label className="cc-checkbox">
                    <input
                        type="checkbox"
                        checked={separation}
                        onChange={(e) => setSeparation(e.target.checked)}
                    />
                    <span className="cc-box"></span>
                    <span className="cc-label-text">
                        {t('enableSeparation')}
                    </span>
                </label>

                <button type="submit" disabled={loading}>
                    {loading ? t('loading') : t('tabGenTitle')}
                </button>
            </form>

            {status && (
                <p className="status-text">
                    {t('statusLabel')}: {status}
                </p>
            )}

            {tabText && (
                <div className="tab-result">
                    <pre>{tabText}</pre>

                    <input
                        type="text"
                        placeholder={t('tabNamePlaceholder')}
                        value={tabName}
                        onChange={(e) => setTabName(e.target.value)}
                    />

                    <div className="tab-actions">
                        <button onClick={handleSave} disabled={saveDisabled}>
                            {t('saveTab')}
                        </button>
                        <button onClick={handleCopy}>
                            {t('copyTab')}
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}