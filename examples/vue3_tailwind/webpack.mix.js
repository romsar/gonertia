const path = require('path');
const mix = require('laravel-mix');
require('mix-tailwindcss');

mix
    .js('resources/js/app.js', 'public/build/assets')
    .postCss('resources/css/app.css', 'build/assets', [])
    .vue(3)
    .alias({'@': path.join(__dirname, 'resources/js/')})
    .setPublicPath('public')
    .version();