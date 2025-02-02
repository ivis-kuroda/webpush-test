function urlBase64ToUint8Array(base64String) {
    const padding = '='.repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding)
        .replace(/-/g, '+')
        .replace(/_/g, '/');
    const rawData = window.atob(base64);
    const outputArray = new Uint8Array(rawData.length);
    for (let i = 0; i < rawData.length; ++i) {
        outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
}

async function registerPush() {
    if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
        return;
    }
    try {
        const reg = await navigator.serviceWorker.register("/assets/js/service-worker.js");
        await navigator.serviceWorker.ready; // Wait for the Service Worker to be active
        const applicationServerKey = urlBase64ToUint8Array("BDOpUfHEw7LFRJWhDxF5TW7SR-kiaOY-_6iFrVweY8rfmi9ySzjxSGWbbm-wwriXwAYWVX5808Pb2U2ApYXYKLc");
        const sub = await reg.pushManager.subscribe({
            userVisibleOnly: true,
            applicationServerKey: applicationServerKey
        });
        const id = "1412"; // Set a unique ID for each user
        await fetch(`/subscribe?id=${id}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(sub.toJSON())
        });
    } catch (e) {
        console.error(e);
    }
}
document.addEventListener("DOMContentLoaded", registerPush);
