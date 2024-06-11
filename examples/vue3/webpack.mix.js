const path = require('path');
const mix = require('laravel-mix');

mix.js('resources/js/app.js', 'public/build/assets')
    .vue(3)
    .alias({'@': path.join(__dirname, 'resources/js/')})
    .setPublicPath('public')
    .version();