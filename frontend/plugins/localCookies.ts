export default defineNuxtPlugin(() => {
    return {
        provide: {
            cookies: {
                getItem: (key: string): string | null => {
                    const cookie = useCookie(key);
                    return cookie.value || null;
                },
                setItem: (key: string, value: string, options?: { maxAge?: number; path?: string }) => {
                    const cookie = useCookie(key, {
                        maxAge: options?.maxAge || 60 * 60 * 24 * 7, // 7 days default
                        path: options?.path || '/',
                        sameSite: 'strict'
                    });
                    cookie.value = value;
                },
                removeItem: (key: string) => {
                    const cookie = useCookie(key);
                    cookie.value = null;
                }
            }
        }
    }
});
