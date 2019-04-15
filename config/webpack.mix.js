let mix = require('laravel-mix');

// 禁止系统通知，这个不断通知挺烦人的。
mix.disableNotifications();
// mix.setPublicPath("resources/assets/");

mix.js("resources/assets/js/app.js", "resources/public/js/");

mix.sass("resources/assets/sass/app.scss", "resources/public/css/");