if(!self.define){const e=e=>{"require"!==e&&(e+=".js");let s=Promise.resolve();return a[e]||(s=new Promise(async s=>{if("document"in self){const a=document.createElement("script");a.src=e,document.head.appendChild(a),a.onload=s}else importScripts(e),s()})),s.then(()=>{if(!a[e])throw new Error(`Module ${e} didn’t register its module`);return a[e]})},s=(s,a)=>{Promise.all(s.map(e)).then(e=>a(1===e.length?e[0]:e))},a={require:Promise.resolve(s)};self.define=(s,c,n)=>{a[s]||(a[s]=Promise.resolve().then(()=>{let a={};const i={uri:location.origin+s.slice(1)};return Promise.all(c.map(s=>{switch(s){case"exports":return a;case"module":return i;default:return e(s)}})).then(e=>{const s=n(...e);return a.default||(a.default=s),a})}))}}define("./sw.js",["./workbox-c2b5e142"],(function(e){"use strict";importScripts(),e.skipWaiting(),e.clientsClaim(),e.precacheAndRoute([{url:"/_next/static/O75g-fC-n3EPzE8g222pJ/_buildManifest.js",revision:"1f23b682cc5add35aab7fa820189742a"},{url:"/_next/static/O75g-fC-n3EPzE8g222pJ/_ssgManifest.js",revision:"abee47769bf307639ace4945f9cfd4ff"},{url:"/_next/static/chunks/30eecaf7486f66aff4d0871082ffc97e8d526c71.7e06f054a157e7bb9983.js",revision:"ed18dd667a51d1b57918712e79b43827"},{url:"/_next/static/chunks/9c7cdbc4c6db4133863d402b80ed82e8f5de782e.23c5a76868458615b585.js",revision:"44c6ff2d1a61e372e6407a793f7ef443"},{url:"/_next/static/chunks/cb1608f2.0b3445cb746d7f436efd.js",revision:"88f3589d5244c11af58b4b86a5bb4acc"},{url:"/_next/static/chunks/commons.d4d525acfdb127c07e33.js",revision:"effec38b53f7044e020bfd202c123f66"},{url:"/_next/static/chunks/framework.4df82c4704a0136f6a4b.js",revision:"7eafc2b810ea3395615465510119d273"},{url:"/_next/static/chunks/main-42d4f7c93b92c6049739.js",revision:"97f6e0595af101b427ed166ed4679a6d"},{url:"/_next/static/chunks/pages/_app-2d2b2892978d6d450823.js",revision:"2fd96694613a8da3e73bcb3354a7bce6"},{url:"/_next/static/chunks/pages/_error-5db79a1b9c0a6f29ad4e.js",revision:"5d53cfb2c5586bb7c7ab1a4025793a04"},{url:"/_next/static/chunks/pages/index-cdaf79ac31a15944514c.js",revision:"a45d8f792f02b20716d61c9254d2cf03"},{url:"/_next/static/chunks/pages/login-53bedb14aa789a5d3d67.js",revision:"19e5cbc2aa1473b7085c08b669ab3d34"},{url:"/_next/static/chunks/pages/success-d44a1fbdeef7c52b305f.js",revision:"f6d15700fe972a72ea76e40d18bbd0d8"},{url:"/_next/static/chunks/polyfills-4beebf4ac9054f0bf4e6.js",revision:"c8b961cfccce0518d96d73f45e46210d"},{url:"/_next/static/chunks/webpack-e067438c4cf4ef2ef178.js",revision:"8c19f623e8389f11131a054a7e17ff95"},{url:"/_next/static/css/4adfad60ab746d68942a.css",revision:"c6a3f91802aa60b8864770e1d2149017"},{url:"/favicon.ico",revision:"4ff59fef4ad8bd2547e3db47bac48f20"},{url:"/icons/icon-128x128.png",revision:"d626cfe7c65e6e5403bcbb9d13aa5053"},{url:"/icons/icon-144x144.png",revision:"e53a506b62999dc7a4f8b7222f8c5add"},{url:"/icons/icon-152x152.png",revision:"18b3958440703a9ecd3c246a0f3f7c72"},{url:"/icons/icon-192x192.png",revision:"27dc12f66697a47b6a8b3ee25ba96257"},{url:"/icons/icon-384x384.png",revision:"a40324a3fde2b0b26eeffd4f08bf8be8"},{url:"/icons/icon-512x512.png",revision:"93d6e8e15cfa78dfee55446f607d9a28"},{url:"/icons/icon-72x72.png",revision:"f2ffc41b3482888f3ae614e0dd2f6980"},{url:"/icons/icon-96x96.png",revision:"fba02a40f7ba6fc65be8a2f245480f6d"},{url:"/manifest.json",revision:"c96057f6fe080d95b52920d55437ade9"},{url:"/test.jpg",revision:"1796eb099970313cfaa469ca908437f8"},{url:"/test2.jpg",revision:"a4ae9fcc8af6123b63e915db6f98daa0"},{url:"/vatinator.svg",revision:"65cefee12ce4d8ae5441aae533ad7b55"},{url:"/vatinatorAsPath.svg",revision:"63c50c40bedf8ebe7151b39717fbe1bc"},{url:"/vercel.svg",revision:"4b4f1876502eb6721764637fe5c41702"}],{ignoreURLParametersMatching:[]}),e.cleanupOutdatedCaches(),e.registerRoute("/",new e.NetworkFirst({cacheName:"start-url",plugins:[new e.ExpirationPlugin({maxEntries:1,maxAgeSeconds:86400,purgeOnQuotaError:!0})]}),"GET"),e.registerRoute(/^https:\/\/fonts\.(?:googleapis|gstatic)\.com\/.*/i,new e.CacheFirst({cacheName:"google-fonts",plugins:[new e.ExpirationPlugin({maxEntries:4,maxAgeSeconds:31536e3,purgeOnQuotaError:!0})]}),"GET"),e.registerRoute(/\.(?:eot|otf|ttc|ttf|woff|woff2|font.css)$/i,new e.StaleWhileRevalidate({cacheName:"static-font-assets",plugins:[new e.ExpirationPlugin({maxEntries:4,maxAgeSeconds:604800,purgeOnQuotaError:!0})]}),"GET"),e.registerRoute(/\.(?:jpg|jpeg|gif|png|svg|ico|webp)$/i,new e.StaleWhileRevalidate({cacheName:"static-image-assets",plugins:[new e.ExpirationPlugin({maxEntries:64,maxAgeSeconds:86400,purgeOnQuotaError:!0})]}),"GET"),e.registerRoute(/\.(?:js)$/i,new e.StaleWhileRevalidate({cacheName:"static-js-assets",plugins:[new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400,purgeOnQuotaError:!0})]}),"GET"),e.registerRoute(/\.(?:css|less)$/i,new e.StaleWhileRevalidate({cacheName:"static-style-assets",plugins:[new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400,purgeOnQuotaError:!0})]}),"GET"),e.registerRoute(/\.(?:json|xml|csv)$/i,new e.NetworkFirst({cacheName:"static-data-assets",plugins:[new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400,purgeOnQuotaError:!0})]}),"GET"),e.registerRoute(/\/api\/.*$/i,new e.NetworkFirst({cacheName:"apis",networkTimeoutSeconds:10,plugins:[new e.ExpirationPlugin({maxEntries:16,maxAgeSeconds:86400,purgeOnQuotaError:!0})]}),"GET"),e.registerRoute(/.*/i,new e.NetworkFirst({cacheName:"others",networkTimeoutSeconds:10,plugins:[new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400,purgeOnQuotaError:!0})]}),"GET")}));
