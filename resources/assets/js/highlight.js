var hljs = require('highlight.js/lib/highlight');

// hljs.registerLanguage('apache', require('highlight.js/lib/languages/apache'));
// hljs.registerLanguage('cpp', require('highlight.js/lib/languages/cpp'));
hljs.registerLanguage('xml', require('highlight.js/lib/languages/xml'));
// hljs.registerLanguage('awk', require('highlight.js/lib/languages/awk'));
hljs.registerLanguage('bash', require('highlight.js/lib/languages/bash'));
// hljs.registerLanguage('basic', require('highlight.js/lib/languages/basic'));
hljs.registerLanguage('cmake', require('highlight.js/lib/languages/cmake'));
hljs.registerLanguage('css', require('highlight.js/lib/languages/css'));
hljs.registerLanguage('markdown', require('highlight.js/lib/languages/markdown'));
hljs.registerLanguage('go', require('highlight.js/lib/languages/go'));
hljs.registerLanguage('ini', require('highlight.js/lib/languages/ini'));
// hljs.registerLanguage('java', require('highlight.js/lib/languages/java'));
hljs.registerLanguage('javascript', require('highlight.js/lib/languages/javascript'));
hljs.registerLanguage('json', require('highlight.js/lib/languages/json'));
// hljs.registerLanguage('makefile', require('highlight.js/lib/languages/makefile'));
hljs.registerLanguage('nginx', require('highlight.js/lib/languages/nginx'));
// hljs.registerLanguage('pgsql', require('highlight.js/lib/languages/pgsql'));
hljs.registerLanguage('php', require('highlight.js/lib/languages/php'));
// hljs.registerLanguage('python', require('highlight.js/lib/languages/python'));
hljs.registerLanguage('scss', require('highlight.js/lib/languages/scss'));
hljs.registerLanguage('shell', require('highlight.js/lib/languages/shell'));
hljs.registerLanguage('sql', require('highlight.js/lib/languages/sql'));
// hljs.registerLanguage('swift', require('highlight.js/lib/languages/swift'));
hljs.registerLanguage('yaml', require('highlight.js/lib/languages/yaml'));
// hljs.registerLanguage('typescript', require('highlight.js/lib/languages/typescript'));
hljs.registerLanguage('vim', require('highlight.js/lib/languages/vim'));

module.exports = hljs;

