const CURRENT_CACHE = 'axwgameboy-pwa';

const filesToCache = [
    "./",
    "./index.html",
    "./gameframe.html",
    "./wasm_exec.js",
    "./style.css",
    "./axwgameboy-wasm.wasm",
];

self.addEventListener('activate', evt =>
    evt.waitUntil(
        caches.keys().then(cacheNames => {
            return Promise.all(
                cacheNames.map(cacheName => {
                    if (cacheName !== CURRENT_CACHE) {
                        return caches.delete(cacheName);
                    }
                })
            );
        })
    )
);

self.addEventListener('install', function (e) {
    e.waitUntil(
        caches.open(CURRENT_CACHE).then(function (cache) {
            return cache.addAll(filesToCache);
        })
    );
});

const fromNetwork = (request, timeout) =>
    new Promise((fulfill, reject) => {
        const timeoutId = setTimeout(reject, timeout);
        fetch(request).then(response => {
            clearTimeout(timeoutId);
            fulfill(response);
            update(request);
        }, reject);
    });

const fromCache = request =>
    caches
        .open(CURRENT_CACHE)
        .then(cache =>
            cache
                .match(request)
                .then(matching => matching || cache.match('/offline/'))
        );

const update = request =>
    caches
        .open(CURRENT_CACHE)
        .then(cache =>
            fetch(request).then(response => cache.put(request, response))
        );

self.addEventListener('fetch', evt => {
    evt.respondWith(
        fromNetwork(evt.request, 10000).catch(() => fromCache(evt.request))
    );
    evt.waitUntil(update(evt.request));
});