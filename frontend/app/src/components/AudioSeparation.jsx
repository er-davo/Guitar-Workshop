import '../styles/AudioSeparation.css';
import React, { useState, useCallback } from 'react';
import { startSeparation, getSeparationTask } from '../api/api.js';
import { useLanguage } from '../i18n/useLanguage.jsx';
import { usePolling } from '../hooks/usePolling.js';

export default function AudioSeparation() {
    const { t } = useLanguage();

    const [file, setFile] = useState(null);
    const [dragActive, setDragActive] = useState(false);

    const [loading, setLoading] = useState(false);
    const [taskId, setTaskId] = useState(null);
    const [status, setStatus] = useState(null);
    const [stems, setStems] = useState(null);
    const [error, setError] = useState(null);

    const statusMap = {
        created: t('statusCreated'),
        pending: t('statusPending'),
        done: t('statusDone'),
        fail: t('statusFail'),
    };

    const handleFileChange = (e) => {
        if (e.target.files?.[0]) {
            setFile(e.target.files[0]);
            setStems(null);
            setError(null);
            setStatus(null);
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
            setStems(null);
            setError(null);
            setStatus(null);
        }
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!file) return alert(t('errorNoFile'));

        setLoading(true);
        setStems(null);
        setError(null);
        setTaskId(null);
        setStatus(t('statusCreated'));

        try {
            const data = await startSeparation(file);
            if (!data?.id) throw new Error(t('errorNoTaskId'));
            setTaskId(data.id);
        } catch (err) {
            setError(err.message);
            setStatus(null);
        } finally {
            setLoading(false);
        }
    };

    usePolling(
        async () => {
            if (!taskId) return null;
            return await getSeparationTask(taskId);
        },
        7000,
        60,
        (result) => {
            if (!result?.status) return false;

            const s = result.status.toLowerCase();
            setStatus(statusMap[s] || s);

            if (s === 'done' && result.separated_audio_signed_urls) {
                setStems(result.separated_audio_signed_urls);
                return true;
            }

            if (s === 'fail') {
                setError(result.error || t('errorSeparation'));
                return true;
            }

            return false;
        },
        [taskId]
    );

    return (
        <div className="audio-separation">
            <h2>{t('audioSepTitle')}</h2>

            <form onSubmit={handleSubmit} className="audio-form">
                <div
                    className={`drop-zone ${dragActive ? 'active' : ''}`}
                    onDragEnter={handleDrag}
                    onDragLeave={handleDrag}
                    onDragOver={handleDrag}
                    onDrop={handleDrop}
                >
                    <input type="file" accept="audio/*" onChange={handleFileChange} />
                    <p>
                        {file
                            ? `${t('fileLabel')}: ${file.name}`
                            : t('dropzonePlaceholder')}
                    </p>
                </div>

                <button type="submit" disabled={loading}>
                    {loading ? t('loading') : t('audioSepTitle')}
                </button>
            </form>

            {status && (
                <p className="status-text">
                    {t('statusLabel')}: {status}
                </p>
            )}

            {error && <p className="error-text">{error}</p>}

            {stems && (
                <div className="stems-container">
                    {Object.entries(stems).map(([name, url]) => (
                        <div key={name} className="stem-item">
                            <p className="stem-name">{name}</p>
                            <audio controls src={url} />
                            <a href={url} download>
                                {t('download')}
                            </a>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}