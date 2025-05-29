function generateUUID(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

export function getSessionId(): string {
  const storageKey = 'browser_session_id';
  let sessionId = sessionStorage.getItem(storageKey);

  if (!sessionId) {
    // 首次打开浏览器会话时创建新的UUID
    sessionId = generateUUID();
    sessionStorage.setItem(storageKey, sessionId);
    console.log('已创建新的会话标识符:', sessionId);
  }

  return sessionId;
}

export function getPersistentId(): string {
  const storageKey = 'browser_persistent_id';
  let persistentId = localStorage.getItem(storageKey);

  if (!persistentId) {
    persistentId = generateUUID();
    localStorage.setItem(storageKey, persistentId);
  }

  return persistentId;
}