window.navigator.serviceWorker.register("/assets/js/service-worker.js").then(reg => {
  reg.addEventListener('updatefound', () => {
    const newWorker = reg.installing;
    newWorker.addEventListener('statechange', () => {
      if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
        // New service worker is waiting to activate
        newWorker.postMessage({ action: 'skipWaiting' });
      }
    });
  });
});

navigator.serviceWorker.addEventListener('controllerchange', () => {
  window.location.reload();
});
