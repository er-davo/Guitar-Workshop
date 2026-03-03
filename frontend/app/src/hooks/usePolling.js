import { useEffect, useRef } from 'react';

export function usePolling(fn, interval = 15000, maxAttempts = Infinity, callback, deps = []) {
    const attemptsRef = useRef(0);
    const timeoutRef = useRef(null);

    useEffect(() => {
        clearTimeout(timeoutRef.current);
        attemptsRef.current = 0;

        if (!fn) return;

        let stopped = false;

        const poll = async () => {
            if (stopped || attemptsRef.current >= maxAttempts) return;

            attemptsRef.current += 1;

            try {
                const result = await fn();
                if (result != null) {
                    const stop = await callback(result);
                    if (stop) {
                        stopped = true;
                        return;
                    }
                }
            } catch (err) {
                console.error('Polling error:', err);
            }

            if (!stopped) {
                timeoutRef.current = setTimeout(poll, interval);
            }
        };

        poll();

        return () => {
            stopped = true;
            clearTimeout(timeoutRef.current);
        };
    }, deps);
}