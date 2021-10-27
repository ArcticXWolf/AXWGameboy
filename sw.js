var cacheName = 'axwgameboy-pwa';

var filesToCache = [
    "/",
    "/index.html",
    "/gameframe.html",
    "/wasm_exec.js",
    "/style.css",
    "/axwgameboy-wasm.wasm",
];

self.addEventListener('install', function (e) {
    e.waitUntil(
        caches.open(cacheName).then(function (cache) {
            return cache.addAll(filesToCache);
        })
    );
});

self.addEventListener('fetch', function (e) {
    e.respondWith(
        caches.match(e.request).then(function (response) {
            return response || fetch(e.request);
        })
    );
});