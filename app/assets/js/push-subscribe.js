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
    alert("Push messaging is not supported in your browser.");
    return;
  }
  if (Notification.permission === "default") {
    await Notification.requestPermission();
  }
  try {
    const reg = await navigator.serviceWorker.register("/assets/js/service-worker.js");
    if (! reg.active) {
      alert("Service worker is not active yet.");
      return;
    }
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
    console.log("Push subscription:");
  } catch (e) {
    console.error(e);
  }
}

async function unregisterPush() {
  const id = "1412"; // Set a unique ID for each user
  await fetch(`/unsubscribe?id=${id}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }
  });
  console.log("Push unsubscription:");
}

document.addEventListener("DOMContentLoaded", function () {
  const checkbox = document.getElementById('push-subscribe-checkbox');
  checkbox.addEventListener('change', function (event) {
    checkbox.disabled = true; // Prevent multiple subscriptions
    if (event.target.checked) {
      console.log("Checkbox checked:", event.target.checked);
      registerPush().then(() => checkbox.disabled = false);
    } else {
      console.log("Checkbox unchecked:", event.target.checked);
      unregisterPush().then(() => checkbox.disabled = false);
    }
  });
});
